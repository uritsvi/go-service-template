package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var ErrMissingPort = errors.New("missing required config: SERVER_PORT")

// Server holds server-level configuration.
type Server struct {
	Port string
}

// Config is the root configuration struct. It holds server settings and feature configs.
type Config struct {
	Server Server
}

// Load loads configuration from the environment. It loads the .env file first (via godotenv),
// then reads all config from the environment.
// Returns config and an error if required fields are missing.
func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return nil, ErrMissingPort
	}

	return &Config{Server: Server{Port: port}}, nil
}
