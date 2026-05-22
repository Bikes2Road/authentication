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

type User struct {
	ID          string    `json:"id"`
	NickName    string    `json:"nick_name"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Password    string    `json:"password"`
	HasPassword bool      `json:"has_password"`
	IsActive    bool      `json:"is_active"`
	Role        string    `json:"role"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// IsValid verifica si el usuario es válido para autenticación
func (u *User) IsValid() bool {
	return u.ID != "" && u.Email != "" && u.IsActive
}
