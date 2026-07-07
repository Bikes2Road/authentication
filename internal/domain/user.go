package domain

import "time"

// User representa la información básica del usuario necesaria para autenticación
type UserAuth struct {
	ID        string `json:"id"`
	NickName  string `json:"nick_name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  bool   `json:"is_active"`
}

// CompanyInfo agrupa los datos de empresa asociados a un usuario.
// Se popula desde users.company_id, users.company_role y companies.suscription_type
// mediante un LEFT JOIN en el repositorio. Es nil cuando el usuario no tiene
// empresa asociada.
type CompanyInfo struct {
	ID              string `json:"company_id"`
	Role            string `json:"company_role"`
	SuscriptionType string `json:"suscription_type"`
}

type User struct {
	ID              string       `json:"id"`
	NickName        string       `json:"nick_name"`
	FirstName       string       `json:"first_name"`
	LastName        string       `json:"last_name"`
	Email           string       `json:"email"`
	PhoneNumber     string       `json:"phone_number"`
	Password        string       `json:"password"`
	HasPassword     bool         `json:"has_password"`
	IsActive        bool         `json:"is_active"`
	Role            string       `json:"role"`
	SuscriptionType string       `json:"suscription_type,omitempty"`
	Company         *CompanyInfo `json:"company,omitempty"`
	DateCreated     time.Time    `json:"date_created"`
	DateUpdated     time.Time    `json:"date_updated"`
}

// IsValid verifica si el usuario es válido para autenticación
func (u *User) IsValid() bool {
	return u.ID != "" && u.Email != "" && u.IsActive
}
