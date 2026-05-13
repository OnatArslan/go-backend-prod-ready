package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App  AppConfig
	HTTP HTTPConfig
	DB   DBConfig
	Log  LogConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type HTTPConfig struct {
	Addr            string
	RequestTimeout  time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type DBConfig struct {
	DSN            string
	MaxConns       int32
	ConnectTimeout time.Duration
}

type LogConfig struct {
	Level  slog.Level
	Format string
}

func Load() (Config, error) {
	cfg := Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "go-backend-prod-ready"),
			Env:  strings.ToLower(getEnv("APP_ENV", "local")),
		},
		HTTP: HTTPConfig{
			Addr:            getEnv("HTTP_ADDR", ":8080"),
			RequestTimeout:  getDuration("HTTP_REQUEST_TIMEOUT", 30*time.Second),
			ReadTimeout:     getDuration("HTTP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDuration("HTTP_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:     getDuration("HTTP_IDLE_TIMEOUT", time.Minute),
			ShutdownTimeout: getDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		DB: DBConfig{
			DSN:            getEnv("DATABASE_URL", ""),
			MaxConns:       int32(getInt("DB_MAX_CONNS", 10)),
			ConnectTimeout: getDuration("DB_CONNECT_TIMEOUT", 5*time.Second),
		},
		Log: LogConfig{
			Level:  getLogLevel("LOG_LEVEL", slog.LevelInfo),
			Format: getLogFormat("LOG_FORMAT", "json"),
		},
	}

	if cfg.DB.DSN == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.DB.MaxConns < 1 {
		return Config{}, fmt.Errorf("DB_MAX_CONNS must be greater than 0")
	}

	return cfg, nil
}

func (cfg Config) IsProduction() bool {
	return cfg.App.Env == "production"
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return duration
}

func getInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return n
}

func getLogLevel(key string, fallback slog.Level) slog.Level {
	switch strings.ToLower(getEnv(key, "")) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return fallback
	}
}

func getLogFormat(key, fallback string) string {
	switch strings.ToLower(getEnv(key, fallback)) {
	case "json":
		return "json"
	case "text":
		return "text"
	default:
		return fallback
	}
}
