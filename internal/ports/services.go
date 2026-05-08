package ports

import (
	"context"

	"github.com/bikes2road/authentication/internal/domain"
)

type VerifyUserRequest struct {
	EmailOrNickName string `json:"email_or_nick_name" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

type UserInfoOAuth struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	NickName  string `json:"nick_name"`
	Role      string `json:"role"`
}

// AuthService define la interfaz para el servicio de autenticación
type AuthService interface {
	Login(ctx context.Context, req VerifyUserRequest) (*domain.LoginResponse, error)
	OauthLogin(ctx context.Context, req UserInfoOAuth) (*domain.LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*domain.ValidateResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshResponse, error)
}

// JWTService define la interfaz para el servicio de JWT
type JWTService interface {
	GenerateTokenPair(user *domain.User) (*domain.TokenPair, error)
	ValidateToken(tokenString string, tokenType domain.TokenType) (*domain.JWTClaims, error)
	ParseToken(tokenString string) (*domain.JWTClaims, error)
}

// UserService define la interfaz para el cliente del servicio de usuarios
type UserService interface {
	GetUserByEmailOrNickName(ctx context.Context, emailOrNickName string) (*domain.User, error)
	VerifyUser(ctx context.Context, req VerifyUserRequest) (*domain.User, error)
}
