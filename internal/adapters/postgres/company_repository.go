package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/bikes2road/authentication/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type companyRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyRepository(pool *pgxpool.Pool) ports.CompanyRepository {
	return &companyRepository{pool: pool}
}

func (r *companyRepository) GetCompanyIDByUserID(ctx context.Context, userID string) (*string, error) {
	var companyID string
	err := r.pool.QueryRow(ctx,
		`SELECT company_id FROM companies WHERE user_id = $1 LIMIT 1`,
		userID,
	).Scan(&companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get company id by user id: %w", err)
	}
	return &companyID, nil
}
