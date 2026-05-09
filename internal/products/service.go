package products

import "context"

type Repository interface {
	FindAll(ctx context.Context) error
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
