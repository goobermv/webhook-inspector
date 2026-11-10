package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr            string
	DatabaseURL         string
	AdminAPIKey         string
	DefaultMaxRequests  int
	RetentionHours      int
	RateLimitPerMinute  int
	MigrationsDir       string
	AllowPrivateTargets bool
}

func Load() Config {
	return Config{
		HTTPAddr:            getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		AdminAPIKey:         getEnv("ADMIN_API_KEY", ""),
		DefaultMaxRequests:  getEnvInt("DEFAULT_BIN_MAX_REQUESTS", 1000),
		RetentionHours:      getEnvInt("RETENTION_HOURS", 24),
		RateLimitPerMinute:  getEnvInt("RATE_LIMIT_PER_MINUTE", 120),
		MigrationsDir:       getEnv("MIGRATIONS_DIR", "file://internal/database/migrations"),
		AllowPrivateTargets: getEnvBool("ALLOW_PRIVATE_TARGETS", false),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
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

func getEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
