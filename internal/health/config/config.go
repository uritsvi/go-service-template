package config

import "os"

// Config holds health feature configuration. Part of the root config.
type Config struct {
	BasePath string
}

// Load builds health config from environment variables.
func Load() *Config {
	basePath := os.Getenv("HEALTH_BASE_PATH")
	if basePath == "" {
		basePath = "/health"
	}
	return &Config{BasePath: basePath}
}
