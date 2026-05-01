package domain

// User representa la información básica del usuario necesaria para autenticación
type User struct {
	ID        string `json:"id"`
	NickName  string `json:"nick_name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  bool   `json:"is_active"`
}

// IsValid verifica si el usuario es válido para autenticación
func (u *User) IsValid() bool {
	return u.ID != "" && u.Email != "" && u.IsActive
}
