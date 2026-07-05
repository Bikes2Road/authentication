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
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nick_name VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(50),
    role user_role NOT NULL DEFAULT 'user',
    register_platform register_platform NOT NULL DEFAULT 'bikes2road',
    is_active BOOLEAN NOT NULL DEFAULT true,
    has_password BOOLEAN NOT NULL DEFAULT true,
    email_verified BOOLEAN NOT NULL DEFAULT true,
    send_emails BOOLEAN NOT NULL DEFAULT false,
    is_banned BOOLEAN DEFAULT false,
    suscription BOOLEAN NOT NULL DEFAULT false,
    suscription_type suscription_type NOT NULL DEFAULT 'none',
    date_finish_suscription TIMESTAMPTZ,
    payment_customer_id VARCHAR(255) UNIQUE,
    refresh_token TEXT,
    date_created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    date_updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    company_id UUID,
    company_role company_user_role,
    CONSTRAINT fk_users_company FOREIGN KEY (company_id) REFERENCES companies(company_id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS companies (
    company_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_name VARCHAR(255) NOT NULL,
    nit VARCHAR(255) NOT NULL UNIQUE,
    address VARCHAR(255)[],
    phone_number VARCHAR(255),
    website_url VARCHAR(255),
    logo_url VARCHAR(255),
    email_company VARCHAR(255),
    facebook_url VARCHAR(255),
    instagram_url VARCHAR(255),
    is_verified BOOLEAN NOT NULL DEFAULT false,
    verified_at TIMESTAMPTZ,
    verified_by UUID,
    date_created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    date_updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    nick_company VARCHAR(255) NOT NULL UNIQUE,
    tiktok_url VARCHAR(255),
    suscription BOOLEAN NOT NULL DEFAULT false,
    date_finish_suscription TIMESTAMPTZ,
    suscription_type suscription_type_company NOT NULL DEFAULT 'none'
);

CREATE TABLE IF NOT EXISTS subscription_plans (
    plan_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    max_bikes_allowed INTEGER NOT NULL DEFAULT 0,
    date_created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    date_updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
		`)
		return fmt.Errorf("users table not found, please create it in Supabase Dashboard: %w", err)
	}

	log.Println("Users table exists and is accessible")
	return nil
}

// User represents the database schema for Supabase.
// Los campos CompanyID/CompanyRole/CompanySuscriptionType se populan
// mediante la relación FK con companies.
type User struct {
	ID                     string    `json:"id"`
	NickName               string    `json:"nick_name"`
	FirstName              *string   `json:"first_name"`
	LastName               *string   `json:"last_name"`
	Email                  string    `json:"email"`
	Password               string    `json:"password"`
	IsActive               bool      `json:"is_active"`
	Role                   string    `json:"role"`
	PhoneNumber            *string   `json:"phone_number"`
	HasPassword            bool      `json:"has_password"`
	SuscriptionType        *string   `json:"suscription_type"`
	CompanyID              *string   `json:"company_id"`
	CompanyRole            *string   `json:"company_role"`
	CompanySuscriptionType *string   `json:"company_suscription_type"`
	DateCreated            time.Time `json:"date_created"`
	DateUpdated            time.Time `json:"date_updated"`
}

// toSupabaseUser converts domain.User to Supabase User
func toSupabaseUser(user *domain.User) *User {
	var firstName, lastName, phoneNumber *string
	if user.FirstName != "" {
		firstName = &user.FirstName
	}
	if user.LastName != "" {
		lastName = &user.LastName
	}
	if user.PhoneNumber != "" {
		phoneNumber = &user.PhoneNumber
	}
	var suscriptionType *string
	if user.SuscriptionType != "" {
		suscriptionType = &user.SuscriptionType
	}
	var companyID, companyRole *string
	if user.Company != nil {
		companyID = &user.Company.ID
		if user.Company.Role != "" {
			companyRole = &user.Company.Role
		}
	}
	return &User{
		ID:              user.ID,
		NickName:        user.NickName,
		FirstName:       firstName,
		LastName:        lastName,
		Email:           user.Email,
		Password:        user.Password,
		IsActive:        user.IsActive,
		Role:            user.Role,
		PhoneNumber:     phoneNumber,
		HasPassword:     user.HasPassword,
		SuscriptionType: suscriptionType,
		CompanyID:       companyID,
		CompanyRole:     companyRole,
		DateCreated:     user.DateCreated,
		DateUpdated:     user.DateUpdated,
	}
}

// toDomainUser converts Supabase User to domain.User
func toDomainUser(user *User) *domain.User {
	firstName := ""
	if user.FirstName != nil {
		firstName = *user.FirstName
	}
	lastName := ""
	if user.LastName != nil {
		lastName = *user.LastName
	}
	phoneNumber := ""
	if user.PhoneNumber != nil {
		phoneNumber = *user.PhoneNumber
	}
	suscriptionType := ""
	if user.SuscriptionType != nil {
		suscriptionType = *user.SuscriptionType
	}

	d := &domain.User{
		ID:              user.ID,
		NickName:        user.NickName,
		FirstName:       firstName,
		LastName:        lastName,
		Email:           user.Email,
		Password:        user.Password,
		IsActive:        user.IsActive,
		Role:            user.Role,
		PhoneNumber:     phoneNumber,
		HasPassword:     user.HasPassword,
		SuscriptionType: suscriptionType,
		DateCreated:     user.DateCreated,
		DateUpdated:     user.DateUpdated,
	}

	if user.CompanyID != nil && *user.CompanyID != "" {
		info := &domain.CompanyInfo{ID: *user.CompanyID}
		if user.CompanyRole != nil {
			info.Role = *user.CompanyRole
		}
		if user.CompanySuscriptionType != nil {
			info.SuscriptionType = *user.CompanySuscriptionType
		}
		d.Company = info
	}

	return d
}
