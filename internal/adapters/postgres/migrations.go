package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bikes2road/authentication/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool) error {
	log.Println("Running PostgreSQL migrations...")

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		nick_name VARCHAR(255) NOT NULL,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		is_active BOOLEAN NOT NULL DEFAULT true,
		role VARCHAR(50) NOT NULL DEFAULT 'user',
		phone_number VARCHAR(50),
		has_password BOOLEAN DEFAULT false,
		date_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		date_updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_nick_name ON users(nick_name);
	CREATE INDEX IF NOT EXISTS idx_users_date_created ON users(date_created);
	`

	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("PostgreSQL migrations completed successfully")
	return nil
}

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
	HasPassword bool      `json:"has_password"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

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
		HasPassword: user.HasPassword,
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}
}