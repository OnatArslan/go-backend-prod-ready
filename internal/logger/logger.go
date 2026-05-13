package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/OnatArslan/go-backend-prod-ready/internal/config"
)

func New(cfg config.LogConfig) *slog.Logger {
	return NewWithWriter(os.Stdout, cfg)
}

func NewWithWriter(w io.Writer, cfg config.LogConfig) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: cfg.Level,
	}

	if cfg.Format == "text" {
		return slog.New(slog.NewTextHandler(w, opts))
	}

	return slog.New(slog.NewJSONHandler(w, opts))
}
