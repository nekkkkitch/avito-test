package products

import (
	"context"
	"market/models"

	"github.com/google/uuid"
)

type MarketRepo interface {
	ListProducts(ctx context.Context) ([]models.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (models.Product, error)
}

type Service struct {
	repo MarketRepo
}

func New(repo MarketRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListProducts(ctx context.Context) ([]models.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *Service) GetProductByID(ctx context.Context, id uuid.UUID) (models.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}
