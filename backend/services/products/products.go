package products

import (
	"backend/models"

	"github.com/google/uuid"
)

type RDBProducts interface {
	GetMarkets() ([]models.Market, error)
	FilterProducts(filter models.Filter) ([]models.Product, error)
	GetProductCard(id uuid.UUID) (models.Product, error)
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

func (s *Service) GetMarkets() ([]models.Market, error) {
	return s.products.GetMarkets()
}

func (s *Service) FilterProducts(filter models.Filter) ([]models.Product, error) {
	return s.products.FilterProducts(filter)
}

func (s *Service) GetProductCard(id uuid.UUID) (models.Product, error) {
	return s.products.GetProductCard(id)
}
