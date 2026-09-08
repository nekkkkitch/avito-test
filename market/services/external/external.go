package external

import (
	"context"
	"market/models"

	"github.com/google/uuid"
)

type MarketRepo interface {
	SaveProduct(ctx context.Context, product models.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo MarketRepo
}

func New(repo MarketRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveProduct(ctx context.Context, product models.Product) error {
	return s.repo.SaveProduct(ctx, product)
}

func (s *Service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, id)
}
