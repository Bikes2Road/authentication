package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders configura cabeceras de seguridad HTTP básicas recomendadas para la API.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Protege contra ataques XSS
		c.Header("X-XSS-Protection", "1; mode=block")
		// Protege contra MIME sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// Previene clickjacking
		c.Header("X-Frame-Options", "DENY")
		// Establece política de seguridad de contenido básica (opcional, ajustada según el caso de uso)
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		c.Next()
	}
}
