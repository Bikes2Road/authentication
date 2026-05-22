package services

import (
	"context"
	"fmt"

	"github.com/bikes2road/authentication/internal/domain"
	"github.com/bikes2road/authentication/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

// userService implements the UserService port
type userService struct {
	repo ports.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{
		repo: repo,
	}
}

// GetUserByEmailOrNickName retrieves a user by their email or nick name
func (s *userService) GetUserByEmailOrNickName(ctx context.Context, emailOrNickName string) (*domain.User, error) {
	user, err := s.repo.GetByEmailOrNickName(ctx, emailOrNickName)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email or nick name: %w", err)
	}
	return user, nil
}

// VerifyUser checks if the provided credentials are valid and returns user info
func (s *userService) VerifyUser(ctx context.Context, req ports.VerifyUserRequest) (*domain.User, error) {
	user, err := s.repo.GetByEmailOrNickName(ctx, req.EmailOrNickName)
	if err != nil {
		// Mask "not found" as invalid credentials to avoid user enumeration
		return nil, domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}
