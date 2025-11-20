package main

import (
	"fmt"
	"log"

	"github.com/bikes2road/authentication/cmd/api/config"
	"github.com/bikes2road/authentication/cmd/api/container"

	_ "github.com/bikes2road/authentication/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Bikes2Road Authentication API
// @version         1.0
// @description     Microservicio de autenticación para validar y crear JWT tokens
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@bikes2road.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Crear container con dependencias
	c := container.New(cfg)

	// Configurar Swagger
	c.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Iniciar servidor
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting authentication service on %s", addr)
	log.Printf("Swagger documentation available at http://%s/swagger/index.html", addr)

	if err := c.Router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
