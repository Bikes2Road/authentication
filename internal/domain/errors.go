package domain

import "errors"

var (
	// ErrInvalidCredentials se retorna cuando las credenciales son inválidas
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserNotFound se retorna cuando el usuario no existe
	ErrUserNotFound = errors.New("user not found")

	// ErrUserInactive se retorna cuando el usuario está inactivo
	ErrUserInactive = errors.New("user is inactive")

	// ErrInvalidToken se retorna cuando el token es inválido
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenExpired se retorna cuando el token ha expirado
	ErrTokenExpired = errors.New("token has expired")

	// ErrTokenMalformed se retorna cuando el token está mal formado
	ErrTokenMalformed = errors.New("token is malformed")

	// ErrUnauthorized se retorna cuando no hay autorización
	ErrUnauthorized = errors.New("unauthorized")

	// ErrUserServiceUnavailable se retorna cuando el servicio de usuarios no está disponible
	ErrUserServiceUnavailable = errors.New("user service unavailable")

	// ErrInternalServer se retorna para errores internos del servidor
	ErrInternalServer = errors.New("internal server error")
)
