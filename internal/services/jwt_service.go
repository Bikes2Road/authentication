package services

import (
	"fmt"
	"time"

	"github.com/bikes2road/authentication/internal/domain"
	"github.com/bikes2road/authentication/internal/ports"
	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey              []byte
	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
}

// NewJWTService crea una nueva instancia del servicio JWT
func NewJWTService(secretKey string, accessTokenExpiration, refreshTokenExpiration time.Duration) ports.JWTService {
	return &jwtService{
		secretKey:              []byte(secretKey),
		accessTokenExpiration:  accessTokenExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}
}

// GenerateTokenPair genera un par de tokens (access y refresh) para un usuario
func (s *jwtService) GenerateTokenPair(user *domain.User) (*domain.TokenPair, error) {
	// Generar access token
	accessToken, err := s.generateToken(user, domain.AccessToken, s.accessTokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generar refresh token
	refreshToken, err := s.generateToken(user, domain.RefreshToken, s.refreshTokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTokenExpiration.Seconds()),
	}, nil
}

// generateToken genera un token JWT
func (s *jwtService) generateToken(user *domain.User, tokenType domain.TokenType, expiration time.Duration) (string, error) {
	now := time.Now()
	expirationTime := now.Add(expiration)

	claims := &domain.JWTClaims{
		UserID:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		// Campos explícitos para swagger
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		Issuer:    "bikes2road-auth",
		Subject:   user.ID,
		ID:        fmt.Sprintf("%s-%s-%d", user.ID, tokenType, now.Unix()),
		// También llenar RegisteredClaims para compatibilidad con jwt library
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "bikes2road-auth",
			Subject:   user.ID,
			ID:        fmt.Sprintf("%s-%s-%d", user.ID, tokenType, now.Unix()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken valida un token y retorna sus claims
func (s *jwtService) ValidateToken(tokenString string, tokenType domain.TokenType) (*domain.JWTClaims, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Verificar que el token no haya expirado usando el campo explícito
	if claims.ExpiresAt > 0 && claims.ExpiresAt < time.Now().Unix() {
		return nil, domain.ErrTokenExpired
	}

	return claims, nil
}

// ParseToken parsea un token y retorna sus claims sin validar expiración
func (s *jwtService) ParseToken(tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firma sea el esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		if err == jwt.ErrTokenMalformed {
			return nil, domain.ErrTokenMalformed
		} else if err == jwt.ErrTokenExpired {
			return nil, domain.ErrTokenExpired
		}
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*domain.JWTClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}
