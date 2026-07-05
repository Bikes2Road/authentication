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

	// Limpieza para entornos de desarrollo: el nuevo schema no es compatible
	// con el anterior. CASCADE elimina dependencias (FKs, constraints, etc.).
	if _, err := pool.Exec(context.Background(),
		`DROP TABLE IF EXISTS subscription_plans, companies, users CASCADE;`); err != nil {
		return fmt.Errorf("failed to drop existing tables: %w", err)
	}

	query := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	DO $$ BEGIN
		CREATE TYPE user_role AS ENUM ('admin', 'user', 'mainteiner', 'company');
	EXCEPTION WHEN duplicate_object THEN NULL; END $$;

	DO $$ BEGIN
		CREATE TYPE company_user_role AS ENUM ('admin', 'colaborator');
	EXCEPTION WHEN duplicate_object THEN NULL; END $$;

	DO $$ BEGIN
		CREATE TYPE suscription_type AS ENUM ('basic', 'standard', 'premium', 'none');
	EXCEPTION WHEN duplicate_object THEN NULL; END $$;

	DO $$ BEGIN
		CREATE TYPE suscription_type_company AS ENUM ('basic', 'standard', 'premium', 'none');
	EXCEPTION WHEN duplicate_object THEN NULL; END $$;

	DO $$ BEGIN
		CREATE TYPE register_platform AS ENUM ('bikes2road', 'google', 'facebook', 'apple');
	EXCEPTION WHEN duplicate_object THEN NULL; END $$;

	-- Tablas independientes primero (sin FKs hacia users/companies)
	CREATE TABLE IF NOT EXISTS "companies" (
		"company_id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
		"company_name" varchar NOT NULL,
		"nit" varchar NOT NULL CONSTRAINT "companies_nit_key" UNIQUE,
		"address" varchar[],
		"phone_number" varchar,
		"website_url" varchar,
		"logo_url" varchar,
		"email_company" varchar,
		"facebook_url" varchar,
		"instagram_url" varchar,
		"is_verified" boolean DEFAULT false NOT NULL,
		"verified_at" timestamp with time zone,
		"verified_by" uuid,
		"date_created" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
		"date_updated" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
		"nick_company" varchar NOT NULL CONSTRAINT "companies_nick_company_key" UNIQUE,
		"tiktok_url" varchar,
		"suscription" boolean DEFAULT false NOT NULL,
		"date_finish_suscription" timestamp with time zone,
		"suscription_type" suscription_type_company DEFAULT 'none' NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_companies_nit ON "companies"(nit);
	CREATE INDEX IF NOT EXISTS idx_companies_nick_company ON "companies"(nick_company);

	CREATE TABLE IF NOT EXISTS "subscription_plans" (
		"plan_id" varchar PRIMARY KEY,
		"name" varchar NOT NULL,
		"max_bikes_allowed" integer DEFAULT 0 NOT NULL,
		"date_created" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
		"date_updated" timestamp with time zone DEFAULT CURRENT_TIMESTAMP
	);

	-- users al final: tiene FK hacia companies
	CREATE TABLE IF NOT EXISTS "users" (
		"id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
		"nick_name" varchar NOT NULL CONSTRAINT "users_nick_name_key" UNIQUE,
		"first_name" varchar,
		"last_name" varchar,
		"email" varchar NOT NULL CONSTRAINT "users_email_key" UNIQUE,
		"password" varchar NOT NULL,
		"phone_number" varchar,
		"role" user_role DEFAULT 'user' NOT NULL,
		"register_platform" register_platform DEFAULT 'bikes2road' NOT NULL,
		"is_active" boolean DEFAULT true NOT NULL,
		"has_password" boolean DEFAULT true NOT NULL,
		"email_verified" boolean DEFAULT true NOT NULL,
		"send_emails" boolean DEFAULT false NOT NULL,
		"is_banned" boolean DEFAULT false,
		"suscription" boolean DEFAULT false NOT NULL,
		"suscription_type" suscription_type DEFAULT 'none' NOT NULL,
		"date_finish_suscription" timestamp with time zone,
		"payment_customer_id" varchar CONSTRAINT "users_payment_customer_id_key" UNIQUE,
		"refresh_token" text,
		"date_created" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
		"date_updated" timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
		"company_id" uuid,
		"company_role" company_user_role,
		CONSTRAINT fk_users_company FOREIGN KEY ("company_id") REFERENCES "companies"("company_id") ON DELETE SET NULL
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON "users"(email);
	CREATE INDEX IF NOT EXISTS idx_users_nick_name ON "users"(nick_name);
	CREATE INDEX IF NOT EXISTS idx_users_date_created ON "users"(date_created);
	CREATE INDEX IF NOT EXISTS idx_users_company_id ON "users"(company_id);

	-- FK circular: companies.verified_by → users.id, agregada al final
	DO $$ BEGIN
		ALTER TABLE "companies"
			ADD CONSTRAINT fk_companies_verified_by FOREIGN KEY ("verified_by") REFERENCES "users"("id") ON DELETE SET NULL;
	EXCEPTION WHEN duplicate_object THEN NULL; END $$;
	`

	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("PostgreSQL migrations completed successfully")
	return nil
}

type User struct {
	ID                     string
	NickName               string
	FirstName              *string
	LastName               *string
	Email                  string
	Password               string
	IsActive               bool
	Role                   string
	PhoneNumber            *string
	HasPassword            bool
	SuscriptionType        *string
	CompanyID              *string
	CompanyRole            *string
	CompanySuscriptionType *string
	DateCreated            time.Time
	DateUpdated            time.Time
}

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
