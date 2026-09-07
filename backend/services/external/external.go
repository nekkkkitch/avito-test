package external

import (
	"backend/models"
	"context"

	"github.com/google/uuid"
)

type RDBExternal interface {
	SetProduct(ctx context.Context, products []models.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	external RDBExternal
}

func New(e RDBExternal) *Service {
	return &Service{
		external: e,
	}
}

func (s *Service) SetProduct(ctx context.Context, products []models.Product) error {
	return s.external.SetProduct(ctx, products)
}

func (s *Service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.external.DeleteProduct(ctx, id)
}
