package config

import "go-service-template/internal/utils/env"

type Config struct {
	Level          string
	JSONFormat     bool
	MaxMsgLength   int
	AddGoroutineID bool
	OtelEnabled    bool
	OtelEndpoint   string
	ServiceName    string
}

func Load() *Config {
	return &Config{
		Level:          env.String("LOG_LEVEL", "info"),
		JSONFormat:     env.Bool("LOG_JSON", false),
		MaxMsgLength:   env.Int("LOG_MAX_MSG_LENGTH", 0),
		AddGoroutineID: env.Bool("LOG_GOROUTINE_ID", false),
		OtelEnabled:    env.Bool("OTEL_ENABLED", false),
		OtelEndpoint:   env.String("OTEL_ENDPOINT", ""),
		ServiceName:    env.String("SERVICE_NAME", "go-service-template"),
	}
}
