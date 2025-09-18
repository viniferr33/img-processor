package config

import (
	"fmt"
	"os"
)

var (
	// AppEnv is the application environment (e.g., development, production).
	AppEnv     = envOrDefault("ENV", "development")
	ServerPort = envOrDefault("SERVER_PORT", "8080")
	ServerHost = envOrDefault("SERVER_HOST", "0.0.0.0")

	// Database
	DBHost     = envOrPanic("DB_HOST")
	DBPort     = envOrDefault("DB_PORT", "5432")
	DBUser     = envOrPanic("DB_USER")
	DBPassword = envOrPanic("DB_PASSWORD")
	DBName     = envOrPanic("DB_NAME")
	DBSSLMode  = envOrDefault("DB_SSLMODE", "disable")

	// Migrations
	MigrationsPath = envOrDefault("MIGRATIONS_PATH", "migrations")

	// JWT
	JWTSecretKey     = envOrPanic("JWT_SECRET_KEY")
	JWTIssuer        = envOrDefault("JWT_ISSUER", "img-processor")
	JWTExpirationSec = int64(3600) // 1 hour
)

func envOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("missing environment variable: %s", key))
	}
	return value
}

func envOrDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}
