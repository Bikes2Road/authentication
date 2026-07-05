package services

import (
	"context"
	"fmt"

	"github.com/bikes2road/authentication/internal/domain"
	"github.com/bikes2road/authentication/internal/ports"
)

type authService struct {
	jwtService  ports.JWTService
	userService ports.UserService
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(jwtService ports.JWTService, userService ports.UserService) ports.AuthService {
	return &authService{
		jwtService:  jwtService,
		userService: userService,
	}
}

// resolveUserContext limpia o conserva user.Company según la rol del usuario.
// La información de empresa (ID, Role, SuscriptionType) ya viene embebida
// en domain.User desde el repositorio mediante un LEFT JOIN. Esta función
// garantiza que el claim Company solo se emita cuando es válido.
//
// Reglas:
//   - role != "company"           → user.Company = nil
//   - role == "company" pero sin ID → user.Company = nil
//   - role == "company" con ID    → user.Company se conserva tal cual
func (s *authService) resolveUserContext(user *domain.User) {
	if user.Role != domain.RoleCompany {
		user.Company = nil
		return
	}
	if user.Company == nil || user.Company.ID == "" {
		user.Company = nil
	}
}

// Login autentica un usuario y genera tokens JWT
func (s *authService) Login(ctx context.Context, req ports.VerifyUserRequest) (*domain.LoginResponse, error) {
	// Obtener usuario del servicio de usuarios
	user, err := s.userService.VerifyUser(ctx, req)
	if err != nil {
		return nil, err
	}

	s.resolveUserContext(user)

	// Generar tokens
	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Construir respuesta
	response := &domain.LoginResponse{
		User: &domain.UserInfo{
			ID:          user.ID,
			Email:       user.Email,
			NickName:    user.NickName,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Role:        user.Role,
			HasPassword: user.HasPassword,
			Company:     user.Company,
		},
		Tokens: tokens,
	}

	return response, nil
}

func (s *authService) OauthLogin(ctx context.Context, req ports.UserInfoOAuth) (*domain.LoginResponse, error) {
	user, err := s.userService.GetUserByID(ctx, req.ID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to load oauth user: %w", err)
	}

	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	s.resolveUserContext(user)

	// Generar tokens
	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Construir respuesta
	response := &domain.LoginResponse{
		User: &domain.UserInfo{
			ID:          user.ID,
			Email:       user.Email,
			NickName:    user.NickName,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Role:        user.Role,
			HasPassword: user.HasPassword,
			Company:     user.Company,
		},
		Tokens: tokens,
	}

	return response, nil
}

// ValidateToken valida un token JWT
func (s *authService) ValidateToken(ctx context.Context, token string) (*domain.ValidateResponse, error) {
	claims, err := s.jwtService.ValidateToken(token, domain.AccessToken)
	if err != nil {
		return &domain.ValidateResponse{
			Valid:  false,
			Claims: nil,
		}, nil
	}

	return &domain.ValidateResponse{
		Valid:  true,
		Claims: claims,
	}, nil
}

// RefreshToken refresca un token JWT usando el refresh token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshResponse, error) {
	// Validar el refresh token
	claims, err := s.jwtService.ValidateToken(refreshToken, domain.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Obtener usuario actualizado del servicio de usuarios
	user, err := s.userService.GetUserByEmailOrNickName(ctx, claims.NickName)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidToken
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verificar que el usuario siga activo
	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	s.resolveUserContext(user)

	// Generar nuevos tokens
	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &domain.RefreshResponse{
		Tokens: tokens,
	}, nil
}
