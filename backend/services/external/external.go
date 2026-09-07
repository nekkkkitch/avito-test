package external

import (
	"backend/models"

	"github.com/google/uuid"
)

type RDBExternal interface {
	SetProduct([]models.Product) error
	DeleteProduct(id uuid.UUID) error
}

type Service struct {
	external RDBExternal
}

func New(e RDBExternal) *Service {
	return &Service{
		external: e,
	}
}

func (s *Service) SetProduct(products []models.Product) error {
	return s.external.SetProduct(products)
}

func (s *Service) DeleteProduct(id uuid.UUID) error {
	return s.external.DeleteProduct(id)
}
