package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bikes2road/authentication/internal/domain"
	"github.com/bikes2road/authentication/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) ports.UserRepository {
	return &userRepository{pool: pool}
}

const userSelectColumns = `
	u.id, u.nick_name, u.first_name, u.last_name, u.email, u.password,
	u.is_active, u.role, u.phone_number, u.has_password, u.suscription_type,
	u.company_id, u.company_role, c.suscription_type AS company_suscription_type,
	u.date_created, u.date_updated
`

const userSelectFrom = `FROM users u LEFT JOIN companies c ON c.company_id = u.company_id`

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			id, nick_name, first_name, last_name, email, password,
			is_active, role, phone_number, has_password,
			suscription_type, company_id, company_role, date_created, date_updated
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
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
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.NickName, firstName, lastName,
		user.Email, user.Password, user.IsActive, user.Role,
		phoneNumber, user.HasPassword, suscriptionType,
		companyID, companyRole, time.Now(), time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) scanUser(row interface {
	Scan(dest ...any) error
}) (*User, error) {
	user := &User{}
	err := row.Scan(
		&user.ID, &user.NickName, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.IsActive, &user.Role,
		&user.PhoneNumber, &user.HasPassword, &user.SuscriptionType,
		&user.CompanyID, &user.CompanyRole, &user.CompanySuscriptionType,
		&user.DateCreated, &user.DateUpdated,
	)
	return user, err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT ` + userSelectColumns + userSelectFrom + ` WHERE u.id = $1 LIMIT 1`
	user, err := r.scanUser(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return toDomainUser(user), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT ` + userSelectColumns + userSelectFrom + ` WHERE u.email = $1 LIMIT 1`
	user, err := r.scanUser(r.pool.QueryRow(ctx, query, email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return toDomainUser(user), nil
}

func (r *userRepository) GetByNickName(ctx context.Context, nickName string) (*domain.User, error) {
	query := `SELECT ` + userSelectColumns + userSelectFrom + ` WHERE u.nick_name = $1 LIMIT 1`
	user, err := r.scanUser(r.pool.QueryRow(ctx, query, nickName))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by nick name: %w", err)
	}
	return toDomainUser(user), nil
}

func (r *userRepository) GetByEmailOrNickName(ctx context.Context, emailOrNickName string) (*domain.User, error) {
	query := `SELECT ` + userSelectColumns + userSelectFrom + ` WHERE u.nick_name = $1 OR u.email = $2 LIMIT 1`
	user, err := r.scanUser(r.pool.QueryRow(ctx, query, emailOrNickName, emailOrNickName))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		log.Printf("unexpected error looking for email/nickname: %v", err)
		return nil, err
	}
	return toDomainUser(user), nil
}

func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `SELECT ` + userSelectColumns + userSelectFrom + ` ORDER BY u.date_created DESC LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user, err := r.scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, toDomainUser(user))
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET
			nick_name = $1, first_name = $2, last_name = $3, email = $4, password = $5,
			is_active = $6, role = $7, phone_number = $8, has_password = $9,
			suscription_type = $10, company_id = $11, company_role = $12,
			date_updated = $13
		WHERE id = $14
	`
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
	result, err := r.pool.Exec(ctx, query,
		user.NickName, firstName, lastName, user.Email, user.Password,
		user.IsActive, user.Role, phoneNumber, user.HasPassword,
		suscriptionType, companyID, companyRole, time.Now(), user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return exists, nil
}
