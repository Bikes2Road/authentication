package http

import (
	"github.com/bikes2road/authentication/internal/adapters/http/middleware"
	"github.com/bikes2road/authentication/internal/ports"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configura las rutas de la aplicación
func SetupRouter(authHandler ports.AuthHandler, healthHandler ports.HealthHandler) *gin.Engine {
	router := gin.Default()

	// Aplicar middlewares globales
	router.Use(middleware.CORS())
	router.Use(middleware.SecurityHeaders())

	// Health check endpoint
	router.GET("/health", healthHandler.Health)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
