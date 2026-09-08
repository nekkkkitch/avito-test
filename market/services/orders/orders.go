package orders

import (
	"context"
	"log/slog"

	"market/models"
)

type MarketRepo interface {
	SaveOrder(ctx context.Context, order models.Order) error
	ListOrders(ctx context.Context) ([]models.Order, error)
}

type Service struct {
	repo MarketRepo
}

func New(repo MarketRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveOrder(ctx context.Context, order models.Order) error {
	slog.Info("OrdersService: SaveOrder: got called", "order", order)
	return s.repo.SaveOrder(ctx, order)
}

func (s *Service) ListOrders(ctx context.Context) ([]models.Order, error) {
	slog.Info("OrdersService: ListOrder: got called")
	return s.repo.ListOrders(ctx)
}
