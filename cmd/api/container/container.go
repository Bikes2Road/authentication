package container

import (
	"fmt"

	"github.com/bikes2road/authentication/cmd/api/config"
	httpAdapter "github.com/bikes2road/authentication/internal/adapters/http"
	"github.com/bikes2road/authentication/internal/adapters/supabase"
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
func New(cfg *config.Config) (*Container, error) {
	// Initialize supabase client
	client, err := supabase.NewClient(supabase.ClientConfig{
		URL:    cfg.Supabase.URL,
		APIKey: cfg.Supabase.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize supabase client: %w", err)
	}

	if err := supabase.RunMigrations(client); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	userRepository := supabase.NewUserRepository(client)
	userService := services.NewUserService(userRepository)

	// Crear servicios
	jwtService := services.NewJWTService(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenExpiration,
		cfg.JWT.RefreshTokenExpiration,
	)
	authService := services.NewAuthService(jwtService, userService)

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
	}, nil
}
