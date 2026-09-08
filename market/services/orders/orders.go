package orders

import (
	"context"
	"market/models"
)

type MarketRepo interface {
	SaveOrder(ctx context.Context, order models.Order) error
}

type Service struct {
	repo MarketRepo
}

func New(repo MarketRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveOrder(ctx context.Context, order models.Order) error {
	return s.repo.SaveOrder(ctx, order)
}
