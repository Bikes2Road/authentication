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

// resolveUserContext fills user.CompanyID and user.CompanySuscriptionType when
// the user has the company role. It also clears both fields for non-company
// users so the JWT claims don't leak stale values.
// Returns an error only on database failure; users with role=company and no
// associated company are allowed to authenticate with company_id=nil.
func (s *authService) resolveUserContext(ctx context.Context, user *domain.User) error {
	if user.Role != domain.RoleCompany {
		user.CompanyID = nil
		user.CompanySuscriptionType = nil
		return nil
	}
	companyID, err := s.userService.GetCompanyIDForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	user.CompanyID = companyID

	if companyID == nil {
		user.CompanySuscriptionType = nil
		return nil
	}

	suscriptionType, err := s.userService.GetCompanySuscriptionType(ctx, *companyID)
	if err != nil {
		return err
	}
	user.CompanySuscriptionType = suscriptionType
	return nil
}

// Login autentica un usuario y genera tokens JWT
func (s *authService) Login(ctx context.Context, req ports.VerifyUserRequest) (*domain.LoginResponse, error) {
	// Obtener usuario del servicio de usuarios
	user, err := s.userService.VerifyUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := s.resolveUserContext(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to resolve user context: %w", err)
	}

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
			CompanyID:   user.CompanyID,
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

	if err := s.resolveUserContext(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to resolve user context: %w", err)
	}

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
			CompanyID:   user.CompanyID,
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

	if err := s.resolveUserContext(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to resolve user context: %w", err)
	}

	// Generar nuevos tokens
	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &domain.RefreshResponse{
		Tokens: tokens,
	}, nil
}
