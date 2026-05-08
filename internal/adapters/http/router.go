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

	// API v1 routes
	v1 := router.Group("/v1")
	{
		v1.POST("/login", authHandler.Login)
		v1.POST("/login/oauth", authHandler.OauthLogin)
		v1.POST("/validate", authHandler.Validate)
		v1.POST("/refresh", authHandler.Refresh)
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}
