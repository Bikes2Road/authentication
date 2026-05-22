package ports

import (
	"context"

	"github.com/bikes2road/authentication/internal/domain"
)

// UserRepository defines the interface for user data persistence
// This is a port (output port) that will be implemented by adapters layer
type UserRepository interface {
	// Create persists a new user to the database
	Create(ctx context.Context, user *domain.User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (*domain.User, error)

	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// GetByEmailOrNickName retrieves a user by their email or nick name
	GetByEmailOrNickName(ctx context.Context, emailOrNickName string) (*domain.User, error)

	// GetByNickName retrieves a user by their nick name
	GetByNickName(ctx context.Context, nickName string) (*domain.User, error)

	// GetAll retrieves all users with pagination
	GetAll(ctx context.Context, limit, offset int) ([]*domain.User, error)

	// Update updates an existing user's information
	Update(ctx context.Context, user *domain.User) error

	// Delete removes a user from the database
	Delete(ctx context.Context, id string) error

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
