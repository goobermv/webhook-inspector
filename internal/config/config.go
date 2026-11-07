package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr           string
	DatabaseURL        string
	AdminAPIKey        string
	DefaultMaxRequests int
	RetentionHours     int
	RateLimitPerMinute int
	MigrationsDir      string
}

func Load() Config {
	return Config{
		HTTPAddr:           getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		AdminAPIKey:        getEnv("ADMIN_API_KEY", ""),
		DefaultMaxRequests: getEnvInt("DEFAULT_BIN_MAX_REQUESTS", 1000),
		RetentionHours:     getEnvInt("RETENTION_HOURS", 24),
		RateLimitPerMinute: getEnvInt("RATE_LIMIT_PER_MINUTE", 120),
		MigrationsDir:      getEnv("MIGRATIONS_DIR", "file://internal/database/migrations"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
