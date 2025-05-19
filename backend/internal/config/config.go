package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config representa la configuración de la aplicación
type Config struct {
	// Servidor
	ServerPort string
	ServerHost string

	// Base de datos
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	// JWT
	JWTSecret        string
	JWTRefreshSecret string
	JWTExpiration    time.Duration

	// Aplicación
	LogLevel string
	Env      string
}

// LoadConfig carga la configuración desde variables de entorno o archivo .env
func LoadConfig() *Config {
	// Intentar cargar .env si existe
	godotenv.Load()

	// Configuración por defecto
	config := &Config{
		// Servidor
		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),

		// Base de datos
		DBHost:     getEnv("DB_HOST", "sistema-tours-db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "sistema_tours"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// JWT
		JWTSecret:        getEnv("JWT_SECRET", "sistema-tours-secret-key"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "sistema-tours-refresh-secret-key"),
		JWTExpiration:    time.Hour * 24, // 1 día por defecto

		// Aplicación
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Env:      getEnv("APP_ENV", "development"),
	}

	// Parsear duración de JWT si está definida
	if jwtExp := getEnv("JWT_EXPIRATION_HOURS", ""); jwtExp != "" {
		if hours, err := strconv.Atoi(jwtExp); err == nil {
			config.JWTExpiration = time.Hour * time.Duration(hours)
		}
	}

	return config
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
