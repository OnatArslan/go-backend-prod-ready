package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/OnatArslan/go-backend-prod-ready/internal/config"
	"github.com/OnatArslan/go-backend-prod-ready/internal/db"
	"github.com/OnatArslan/go-backend-prod-ready/internal/httpx"
	"github.com/OnatArslan/go-backend-prod-ready/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	cfg    config.Config
	logger *slog.Logger
	q      *db.Queries
}

func newApplication(cfg config.Config, logger *slog.Logger, q *db.Queries) *application {
	return &application{
		cfg:    cfg,
		logger: logger,
		q:      q,
	}
}

func (app *application) mount() http.Handler {
	productService := products.NewService(app.q)
	productHandler := products.NewHandler(productService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(app.cfg.HTTP.RequestTimeout))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := httpx.Write(w, http.StatusOK, httpx.Envelope{
			"status":  "ok",
			"service": app.cfg.App.Name,
		}); err != nil {
			app.logger.Error("write health response failed", "err", err)
		}
	})

	r.Route("/ecom/v1", func(r chi.Router) {
		r.Mount("/product", productHandler.Routes())
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.cfg.HTTP.Addr,
		Handler:      h,
		ReadTimeout:  app.cfg.HTTP.ReadTimeout,
		WriteTimeout: app.cfg.HTTP.WriteTimeout,
		IdleTimeout:  app.cfg.HTTP.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		app.logger.Info("server started", "addr", app.cfg.HTTP.Addr, "env", app.cfg.App.Env)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		app.logger.Info("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), app.cfg.HTTP.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			_ = srv.Close()
			return fmt.Errorf("shutdown server: %w", err)
		}

		app.logger.Info("server stopped")
		return nil
	}
}
