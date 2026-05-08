package ports

import (
	"github.com/gin-gonic/gin"
)

// AuthHandler define la interfaz para los handlers de autenticación
type AuthHandler interface {
	Login(c *gin.Context)
	OauthLogin(c *gin.Context)
	Validate(c *gin.Context)
	Refresh(c *gin.Context)
}

// HealthHandler define la interfaz para el handler de health check
type HealthHandler interface {
	Health(c *gin.Context)
}
