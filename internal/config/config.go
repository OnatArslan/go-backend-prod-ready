package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTP HTTPConfig
	DB   DBConfig
}

type HTTPConfig struct {
	Addr string
}

type DBConfig struct {
	DSN string
}

func Load() (Config, error) {
	cfg := Config{
		HTTP: HTTPConfig{
			Addr: getEnv("HTTP_ADDR", ":8080"),
		},
		DB: DBConfig{
			DSN: os.Getenv("DATABASE_URL"),
		},
	}

	if cfg.DB.DSN == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
