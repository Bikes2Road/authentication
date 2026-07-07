package postgres

import (
	"github.com/bikes2road/authentication/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type companyRepository struct {
	pool *pgxpool.Pool
}

// NewCompanyRepository se conserva por compatibilidad con el container,
// pero el CompanyRepository actual está vacío: la información de empresa
// viaja embebida en el User mediante el LEFT JOIN del userRepository.
func NewCompanyRepository(pool *pgxpool.Pool) ports.CompanyRepository {
	return &companyRepository{pool: pool}
}
