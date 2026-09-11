// Package config centralizes how the application reads its runtime
// configuration. Nothing else in the codebase calls os.Getenv
// directly — this is the one place that knows about environment
// variables, so changing how config is sourced (env vars today,
// maybe a config file or secrets manager later) only touches this file.
package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string // empty string means "use the in-memory repository"
}

// Load reads configuration from environment variables, applying
// sensible defaults for local development so `go run` still works
// without any .env file being present.
func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

// Addr returns the address net/http should listen on, e.g. ":8080".
func (c Config) Addr() string {
	return fmt.Sprintf(":%s", c.Port)
}
