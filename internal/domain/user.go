package domain

import "time"

// User representa la información básica del usuario necesaria para autenticación
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"` // Hash del password
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsValid verifica si el usuario es válido para autenticación
func (u *User) IsValid() bool {
	return u.ID != "" && u.Email != "" && u.IsActive
}

// FullName retorna el nombre completo del usuario
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}
