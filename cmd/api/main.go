package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	api := application{
		config: cfg,
	}

	// structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	slog.SetDefault(logger)

	h := api.mount()
	if err := api.run(h); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			slog.Error("server has failed to start, err")
			os.Exit(1)
		}
		slog.Error("internal server error")
		os.Exit(1)
	}

}
