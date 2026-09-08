package external

import (
	"context"
	"log/slog"

	"backend/models"

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
	slog.Info("ExternalService: SetProduct: got called", "products", products)
	return s.external.SetProduct(ctx, products)
}

func (s *Service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	slog.Info("ExternalService: DeleteProduct: got called", "id", id)
	return s.external.DeleteProduct(ctx, id)
}
