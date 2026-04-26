package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS configura las políticas de Cross-Origin Resource Sharing.
func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	// Permitir todos los orígenes en desarrollo (ajustar en producción usando variables de entorno u orígenes específicos)
	config.AllowAllOrigins = true
	// Permitir los métodos comunes
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	// Permitir las cabeceras comunes
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "Accept"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	return cors.New(config)
}
