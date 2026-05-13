package products

import (
	"context"

	"github.com/OnatArslan/go-backend-prod-ready/internal/db"
)

type Repository interface {
	FindAll(ctx context.Context) ([]db.Product, error)
}

type svc struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &svc{
		repo: repo,
	}
}

func (svc *svc) ListProducts(context.Context) error {

	return nil
}
