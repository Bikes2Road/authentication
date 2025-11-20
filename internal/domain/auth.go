package domain

// LoginRequest representa la solicitud de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse representa la respuesta exitosa de login
type LoginResponse struct {
	User   *UserInfo  `json:"user"`
	Tokens *TokenPair `json:"tokens"`
}

// UserInfo representa la información del usuario en la respuesta
type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// ValidateRequest representa la solicitud de validación de token
type ValidateRequest struct {
	Token string `json:"token" binding:"required"`
}

// ValidateResponse representa la respuesta de validación de token
type ValidateResponse struct {
	Valid  bool       `json:"valid"`
	Claims *JWTClaims `json:"claims,omitempty"`
}

// RefreshRequest representa la solicitud de refresh de token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshResponse representa la respuesta de refresh de token
type RefreshResponse struct {
	Tokens *TokenPair `json:"tokens"`
}
