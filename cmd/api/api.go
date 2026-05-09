package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/OnatArslan/go-backend-prod-ready/internal/config"
	"github.com/OnatArslan/go-backend-prod-ready/internal/db"
	"github.com/OnatArslan/go-backend-prod-ready/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	cfg config.Config
	// logger
	q *db.Queries
}

// mount
func (app *application) mount() http.Handler {

	productService := products.NewService(nil)
	productHandler := products.NewHandler(productService)

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// health endpoint
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	r.Route("/ecom/v1", func(r chi.Router) {
		// register product routes
		r.Mount("/product", productHandler.Routes())

	})

	return r
}

// run
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.cfg.HTTP.Addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	slog.Info("server has started", "addr", app.cfg.HTTP.Addr)

	return srv.ListenAndServe()

}

// 47:32
