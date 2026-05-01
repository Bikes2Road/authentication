package ports

import (
	"context"

	"github.com/bikes2road/authentication/internal/domain"
)

// AuthService define la interfaz para el servicio de autenticación
type AuthService interface {
	Login(ctx context.Context, emailOrNickName, password string) (*domain.LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*domain.ValidateResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshResponse, error)
}

// JWTService define la interfaz para el servicio de JWT
type JWTService interface {
	GenerateTokenPair(user *domain.User) (*domain.TokenPair, error)
	ValidateToken(tokenString string, tokenType domain.TokenType) (*domain.JWTClaims, error)
	ParseToken(tokenString string) (*domain.JWTClaims, error)
}

// UserServiceClient define la interfaz para el cliente del servicio de usuarios
type UserServiceClient interface {
	GetUserByEmailOrNickName(ctx context.Context, emailOrNickName, password string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}
