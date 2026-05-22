package container

import (
	"fmt"

	"github.com/bikes2road/authentication/cmd/api/config"
	httpAdapter "github.com/bikes2road/authentication/internal/adapters/http"
	"github.com/bikes2road/authentication/internal/adapters/postgres"
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
	pool, err := postgres.NewClient(postgres.ClientConfig{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize postgres client: %w", err)
	}

	if err := postgres.RunMigrations(pool); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	userRepository := postgres.NewUserRepository(pool)
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
