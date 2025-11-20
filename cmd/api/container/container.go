package container

import (
	"github.com/bikes2road/authentication/cmd/api/config"
	httpAdapter "github.com/bikes2road/authentication/internal/adapters/http"
	"github.com/bikes2road/authentication/internal/adapters/httpclient"
	"github.com/bikes2road/authentication/internal/ports"
	"github.com/bikes2road/authentication/internal/services"
	"github.com/gin-gonic/gin"
)

// Container contiene todas las dependencias de la aplicación
type Container struct {
	Config        *config.Config
	AuthHandler   ports.AuthHandler
	HealthHandler ports.HealthHandler
	Router        *gin.Engine
}

// New crea un nuevo container con todas las dependencias inyectadas
func New(cfg *config.Config) *Container {
	// Crear adaptador de cliente HTTP para el servicio de usuarios
	userServiceClient := httpclient.NewUserServiceClient(cfg.Users.BaseURL)

	// Crear servicios
	jwtService := services.NewJWTService(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenExpiration,
		cfg.JWT.RefreshTokenExpiration,
	)
	authService := services.NewAuthService(jwtService, userServiceClient)

	// Crear handlers
	authHandler := httpAdapter.NewAuthHandler(authService)
	healthHandler := httpAdapter.NewHealthHandler()

	// Configurar router
	router := httpAdapter.SetupRouter(authHandler, healthHandler)

	return &Container{
		Config:        cfg,
		AuthHandler:   authHandler,
		HealthHandler: healthHandler,
		Router:        router,
	}
}
