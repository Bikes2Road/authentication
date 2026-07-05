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
	repo        ports.UserRepository
	companyRepo ports.CompanyRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo ports.UserRepository, companyRepo ports.CompanyRepository) ports.UserService {
	return &userService{
		repo:        repo,
		companyRepo: companyRepo,
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

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
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

// GetCompanyIDForUser returns the company_id associated with the given user id.
// Returns (nil, nil) when the user has no associated company (per current auth policy).
func (s *userService) GetCompanyIDForUser(ctx context.Context, userID string) (*string, error) {
	companyID, err := s.companyRepo.GetCompanyIDByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company id for user: %w", err)
	}
	return companyID, nil
}

// GetCompanySuscriptionType returns the suscription_type of the company identified by companyID.
// Returns (nil, nil) when the company has no suscription_type.
func (s *userService) GetCompanySuscriptionType(ctx context.Context, companyID string) (*string, error) {
	suscriptionType, err := s.companyRepo.GetSuscriptionTypeByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company suscription type: %w", err)
	}
	return suscriptionType, nil
}
