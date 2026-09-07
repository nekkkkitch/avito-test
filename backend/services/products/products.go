package products

import (
	"backend/models"
	"context"

	"github.com/google/uuid"
)

type RDBProducts interface {
	GetMarkets(ctx context.Context) ([]models.Market, error)
	FilterProducts(ctx context.Context, filter models.Filter) ([]models.Product, error)
	GetProductCard(ctx context.Context, id uuid.UUID) (models.Product, error)
}

type Service struct {
	products RDBProducts
}

type DBProducts interface {
	GetMarkets() ([]string, error)
	FilterProducts(filter models.Filter) ([]models.Product, error)
	GetProductCard(id uuid.UUID) (models.Product, error)
}

func New(p RDBProducts) *Service {
	return &Service{
		products: p,
	}
}

func (s *Service) GetMarkets(ctx context.Context) ([]models.Market, error) {
	return s.products.GetMarkets(ctx)
}

func (s *Service) FilterProducts(ctx context.Context, filter models.Filter) ([]models.Product, error) {
	return s.products.FilterProducts(ctx, filter)
}

func (s *Service) GetProductCard(ctx context.Context, id uuid.UUID) (models.Product, error) {
	return s.products.GetProductCard(ctx, id)
}
