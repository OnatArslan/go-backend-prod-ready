package products

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	svc Service
}

func NewHandler(svc Service) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.ListProductsHandler)

	return r
}

func (h *handler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	// call the services for ListProducts

	// return json in an HTTP response
}
