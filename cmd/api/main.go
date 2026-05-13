package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/OnatArslan/go-backend-prod-ready/internal/config"
	"github.com/OnatArslan/go-backend-prod-ready/internal/db"
	applogger "github.com/OnatArslan/go-backend-prod-ready/internal/logger"
	"github.com/OnatArslan/go-backend-prod-ready/internal/postgres"
	"github.com/joho/godotenv"
)

func main() {
	loadDotenvForLocal()

	bootstrapLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(bootstrapLogger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "err", err)
		os.Exit(1)
	}

	logger := applogger.New(cfg.Log)
	slog.SetDefault(logger)

	pool, err := postgres.NewPool(context.Background(), cfg.DB)
	if err != nil {
		logger.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	q := db.New(pool)

	api := newApplication(cfg, logger, q)

	h := api.mount()
	if err := api.run(h); err != nil {
		logger.Error("server failed", "err", err)
		os.Exit(1)
	}
}

func loadDotenvForLocal() {
	if os.Getenv("APP_ENV") == "production" {
		return
	}

	_ = godotenv.Load()
}
