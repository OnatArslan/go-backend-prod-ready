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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	products := []int{1, 2, 3, 4, 5}
	// return json in an HTTP response
	httpx.Write(w, http.StatusOK, products)
}
