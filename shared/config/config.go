package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	PostgresURL  string
	RedisURL     string
	JWTSecret    string
	Environment  string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		PostgresURL: getEnv("POSTGRES_URL", "postgres://postgres:password@localhost:5432/ua_game?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}