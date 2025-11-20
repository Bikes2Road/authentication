package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Server ServerConfig
	JWT    JWTConfig
	Users  UsersServiceConfig
}

// ServerConfig contiene la configuración del servidor HTTP
type ServerConfig struct {
	Port string
	Host string
}

// JWTConfig contiene la configuración de JWT
type JWTConfig struct {
	SecretKey              string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

// UsersServiceConfig contiene la configuración del servicio de usuarios
type UsersServiceConfig struct {
	BaseURL string
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8084"),
			Host: getEnv("HOST", "0.0.0.0"),
		},
		JWT: JWTConfig{
			SecretKey:              getEnv("JWT_SECRET_KEY", ""),
			AccessTokenExpiration:  getDurationEnv("JWT_ACCESS_TOKEN_EXPIRATION", 24*time.Hour),
			RefreshTokenExpiration: getDurationEnv("JWT_REFRESH_TOKEN_EXPIRATION", 24*time.Hour),
		},
		Users: UsersServiceConfig{
			BaseURL: getEnv("USERS_SERVICE_URL", "http://localhost:8083"),
		},
	}

	// Validar configuración requerida
	if config.JWT.SecretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is required")
	}

	return config, nil
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getDurationEnv obtiene una duración desde una variable de entorno (en horas)
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	hours, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return time.Duration(hours) * time.Hour
}
