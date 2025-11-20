package http

import (
	"errors"
	"net/http"

	"github.com/bikes2road/authentication/internal/domain"
	"github.com/bikes2road/authentication/internal/ports"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authService ports.AuthService
}

// NewAuthHandler crea una nueva instancia del handler de autenticación
func NewAuthHandler(authService ports.AuthService) ports.AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

// Login godoc
// @Summary      Login de usuario
// @Description  Autentica un usuario con email y password, retorna tokens JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body domain.LoginRequest true "Credenciales de login"
// @Success      200 {object} domain.LoginResponse "Login exitoso"
// @Failure      400 {object} ErrorResponse "Request inválido"
// @Failure      401 {object} ErrorResponse "Credenciales inválidas"
// @Failure      500 {object} ErrorResponse "Error interno del servidor"
// @Router       /auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Validate godoc
// @Summary      Validar token JWT
// @Description  Valida un token JWT y retorna sus claims si es válido
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body domain.ValidateRequest true "Token a validar"
// @Success      200 {object} domain.ValidateResponse "Token validado"
// @Failure      400 {object} ErrorResponse "Request inválido"
// @Failure      500 {object} ErrorResponse "Error interno del servidor"
// @Router       /auth/validate [post]
func (h *authHandler) Validate(c *gin.Context) {
	var req domain.ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	response, err := h.authService.ValidateToken(c.Request.Context(), req.Token)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Refresh godoc
// @Summary      Refrescar token JWT
// @Description  Genera un nuevo par de tokens usando un refresh token válido
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body domain.RefreshRequest true "Refresh token"
// @Success      200 {object} domain.RefreshResponse "Tokens refrescados"
// @Failure      400 {object} ErrorResponse "Request inválido"
// @Failure      401 {object} ErrorResponse "Token inválido o expirado"
// @Failure      500 {object} ErrorResponse "Error interno del servidor"
// @Router       /auth/refresh [post]
func (h *authHandler) Refresh(c *gin.Context) {
	var req domain.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// handleError maneja los errores y retorna la respuesta HTTP apropiada
func (h *authHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid credentials",
		})
	case errors.Is(err, domain.ErrUserNotFound):
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not found",
		})
	case errors.Is(err, domain.ErrUserInactive):
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User is inactive",
		})
	case errors.Is(err, domain.ErrInvalidToken):
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid token",
		})
	case errors.Is(err, domain.ErrTokenExpired):
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Token has expired",
		})
	case errors.Is(err, domain.ErrTokenMalformed):
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Token is malformed",
		})
	case errors.Is(err, domain.ErrUserServiceUnavailable):
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error:   "Service Unavailable",
			Message: "User service is unavailable",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "An unexpected error occurred",
		})
	}
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
