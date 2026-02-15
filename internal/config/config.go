package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"

	healthconfig "go-service-template/internal/health/config"
)

var ErrMissingPort = errors.New("missing required config: SERVER_PORT or PORT")

// Server holds server-level configuration.
type Server struct {
	Port string
}

// Config is the root configuration struct. It holds server settings and all feature configs.
type Config struct {
	Server Server
	Health *healthconfig.Config
}

// Load loads configuration from the environment. It loads the .env file first (via godotenv),
// then reads all config from the environment. If ENV_PATH is set in .env, that file is loaded instead of .env.
// Returns config and an error if required fields are missing.
func Load() (*Config, error) {
	_ = godotenv.Load(".env")
	if p := os.Getenv("ENV_PATH"); p != "" {
		_ = godotenv.Load(p)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		return nil, ErrMissingPort
	}

	cfg := &Config{
		Server: Server{Port: port},
		Health: healthconfig.Load(),
	}
	return cfg, nil
}
