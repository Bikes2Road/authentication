package http

import (
	"github.com/bikes2road/authentication/internal/ports"
	"github.com/gin-gonic/gin"
)

// SetupRouter configura las rutas de la aplicación
func SetupRouter(authHandler ports.AuthHandler, healthHandler ports.HealthHandler) *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", healthHandler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/validate", authHandler.Validate)
			auth.POST("/refresh", authHandler.Refresh)
		}
	}

	return router
}
