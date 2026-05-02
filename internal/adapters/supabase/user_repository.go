package supabase

import (
	"context"
	"fmt"

	"github.com/bikes2road/authentication/internal/domain"
	"github.com/bikes2road/authentication/internal/ports"
	"github.com/supabase-community/postgrest-go"
	supabase "github.com/supabase-community/supabase-go"
)

// userRepository implements the UserRepository port using Supabase
type userRepository struct {
	client *supabase.Client
}

// NewUserRepository creates a new Supabase user repository
func NewUserRepository(client *supabase.Client) ports.UserRepository {
	return &userRepository{
		client: client,
	}
}

// Create inserts a new user into the database
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	supaUser := toSupabaseUser(user)

	var result []User
	_, err := r.client.From("users").Insert(supaUser, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var users []User
	_, err := r.client.From("users").
		Select("*", "", false).
		Eq("id", id).
		Limit(1, "").
		ExecuteTo(&users)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	if len(users) == 0 {
		return nil, domain.ErrUserNotFound
	}

	return toDomainUser(&users[0]), nil
}

// GetByEmail retrieves a user by their email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var users []User
	_, err := r.client.From("users").
		Select("*", "", false).
		Eq("email", email).
		ExecuteTo(&users)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if len(users) == 0 {
		return nil, domain.ErrUserNotFound
	}

	return toDomainUser(&users[0]), nil
}

// GetByNickName retrieves a user by their nick name
func (r *userRepository) GetByNickName(ctx context.Context, nickName string) (*domain.User, error) {
	var users []User
	_, err := r.client.From("users").
		Select("*", "", false).
		Eq("nick_name", nickName).
		ExecuteTo(&users)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by nick name: %w", err)
	}

	if len(users) == 0 {
		return nil, domain.ErrUserNotFound
	}

	return toDomainUser(&users[0]), nil
}

// GetByEmailOrNickName retrieves a user by their email or nick name
func (r *userRepository) GetByEmailOrNickName(ctx context.Context, emailOrNickName string) (*domain.User, error) {
	var users []User
	filter := fmt.Sprintf("nick_name.eq.%s,email.eq.%s", emailOrNickName, emailOrNickName)
	_, err := r.client.From("users").
		Select("id, nick_name, first_name, last_name, email, password, is_active, phone_number", "", false).
		Or(filter, "").
		ExecuteTo(&users)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email or nick name: %w", err)
	}

	if len(users) == 0 {
		return nil, domain.ErrUserNotFound
	}

	return toDomainUser(&users[0]), nil
}

// GetAll retrieves all users with pagination
func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []User
	_, err := r.client.From("users").
		Select("*", "", false).
		Order("date_created", &postgrest.OrderOpts{Ascending: false}).
		Range(offset, offset+limit-1, "").
		ExecuteTo(&users)

	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	domainUsers := make([]*domain.User, 0, len(users))
	for i := range users {
		domainUsers = append(domainUsers, toDomainUser(&users[i]))
	}

	return domainUsers, nil
}

// Update updates an existing user's information
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	supaUser := toSupabaseUser(user)

	var result []User
	_, err := r.client.From("users").
		Update(supaUser, "", "").
		Eq("id", user.ID).
		ExecuteTo(&result)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if len(result) == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete removes a user from the database
func (r *userRepository) Delete(ctx context.Context, id string) error {
	var result []User
	_, err := r.client.From("users").
		Delete("", "").
		Eq("id", id).
		ExecuteTo(&result)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if len(result) == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var users []User
	_, err := r.client.From("users").
		Select("id", "", false).
		Eq("email", email).
		ExecuteTo(&users)

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return len(users) > 0, nil
}
