package products

import (
	"context"
	"log"
	"net/http"

	"github.com/OnatArslan/go-backend-prod-ready/internal/httpx"
	"github.com/go-chi/chi/v5"
)

type Service interface {
	ListProducts(ctx context.Context) error
}

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

	r.Get("/", h.listProducts)

	return r
}

func (h *handler) listProducts(w http.ResponseWriter, r *http.Request) {
	// call the services for ListProducts
	err := h.svc.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		_ = httpx.WriteError(w, http.StatusInternalServerError, "internal_server_error", "internal server error")
		return
	}
	products := []int{1, 2, 3, 4, 5}
	// return json in an HTTP response
	if err := httpx.Write(w, http.StatusOK, products); err != nil {
		log.Println(err)
	}
}
