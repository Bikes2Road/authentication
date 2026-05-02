package supabase

import (
	"fmt"
	"log"
	"time"

	"github.com/bikes2road/authentication/internal/domain"
	supabase "github.com/supabase-community/supabase-go"
)

// RunMigrations executes database migrations using Supabase
// Note: Supabase migrations are typically managed through the Supabase Dashboard or CLI
// This function verifies the users table exists
func RunMigrations(client *supabase.Client) error {
	log.Println("Note: Supabase migrations are best managed through the Supabase Dashboard or CLI")
	log.Println("Checking if users table exists...")

	// Try to query the table to see if it exists
	// If it doesn't exist, you should create it through the Supabase Dashboard

	// Test query to check if table exists
	var testUsers []User
	_, err := client.From("users").Select("*", "", false).ExecuteTo(&testUsers)

	if err != nil {
		log.Println("Users table may not exist. Please create it using the Supabase Dashboard with the following schema:")
		log.Println(`
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY, -- ¡Cambiado a UUID!
    nick_name VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    role user_role NOT NULL DEFAULT 'user',
    phone_number VARCHAR(50),
    date_created TIMESTAMPTZ,
    date_updated TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_date_created ON users(date_created);
		`)
		return fmt.Errorf("users table not found, please create it in Supabase Dashboard: %w", err)
	}

	log.Println("Users table exists and is accessible")
	return nil
}

// User represents the database schema for Supabase
type User struct {
	ID          string    `json:"id"`
	NickName    string    `json:"nick_name"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	IsActive    bool      `json:"is_active"`
	Role        string    `json:"role"`
	PhoneNumber string    `json:"phone_number"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// toSupabaseUser converts domain.User to Supabase User
func toSupabaseUser(user *domain.User) *User {
	return &User{
		ID:          user.ID,
		NickName:    user.NickName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Password:    user.Password,
		IsActive:    user.IsActive,
		Role:        user.Role,
		PhoneNumber: user.PhoneNumber,
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}
}

// toDomainUser converts Supabase User to domain.User
func toDomainUser(user *User) *domain.User {
	return &domain.User{
		ID:          user.ID,
		NickName:    user.NickName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Password:    user.Password,
		IsActive:    user.IsActive,
		Role:        user.Role,
		PhoneNumber: user.PhoneNumber,
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}
}
