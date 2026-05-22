package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Users    UsersServiceConfig
	Postgres PostgresConfig
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

// PostgresConfig contiene la configuración de PostgreSQL
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8084"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		JWT: JWTConfig{
			SecretKey:              getEnv("JWT_SECRET_KEY", ""),
			AccessTokenExpiration:  getDurationEnv("JWT_ACCESS_TOKEN_EXPIRATION", 24*time.Hour),
			RefreshTokenExpiration: getDurationEnv("JWT_REFRESH_TOKEN_EXPIRATION", 24*time.Hour),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "auth_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
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
