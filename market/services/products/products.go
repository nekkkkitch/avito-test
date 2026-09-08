package products

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"market/models"
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
	slog.Info("ProductsService: ListProducts: got called")
	return s.repo.ListProducts(ctx)
}

func (s *Service) GetProductByID(ctx context.Context, id uuid.UUID) (models.Product, error) {
	slog.Info("ProductsService: GetProductByID: got called", "id", id)
	return s.repo.GetProductByID(ctx, id)
}
