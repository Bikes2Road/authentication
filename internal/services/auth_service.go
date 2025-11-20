package services

import (
	"context"
	"fmt"

	"github.com/bikes2road/authentication/internal/domain"
	"github.com/bikes2road/authentication/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	jwtService        ports.JWTService
	userServiceClient ports.UserServiceClient
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(jwtService ports.JWTService, userServiceClient ports.UserServiceClient) ports.AuthService {
	return &authService{
		jwtService:        jwtService,
		userServiceClient: userServiceClient,
	}
}

// Login autentica un usuario y genera tokens JWT
func (s *authService) Login(ctx context.Context, email, password string) (*domain.LoginResponse, error) {
	// Obtener usuario del servicio de usuarios
	user, err := s.userServiceClient.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verificar que el usuario esté activo
	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	// Verificar password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Generar tokens
	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Construir respuesta
	response := &domain.LoginResponse{
		User: &domain.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
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
	user, err := s.userServiceClient.GetUserByID(ctx, claims.UserID)
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

	// Generar nuevos tokens
	tokens, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &domain.RefreshResponse{
		Tokens: tokens,
	}, nil
}
