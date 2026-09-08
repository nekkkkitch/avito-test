package external

import (
	"context"
	"log/slog"

	"market/models"

	"github.com/gofiber/fiber/v3/client"
	"github.com/google/uuid"
)

type MarketRepo interface {
	SaveProduct(ctx context.Context, product models.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo    MarketRepo
	apiLink string
}

func New(repo MarketRepo, apiLink string) *Service {
	return &Service{repo: repo, apiLink: apiLink}
}

func (s *Service) SaveProduct(ctx context.Context, product models.Product) error {
	slog.Info("ExternalService: SaveProduct: got called", "product", product)
	_, err := client.Post(s.apiLink+"/api/product", client.Config{Body: []models.Product{product}})
	if err != nil {
		return err
	}

	return s.repo.SaveProduct(ctx, product)
}

func (s *Service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	slog.Info("ExternalService: DeleteProduct: got called", "id", id)
	_, err := client.Delete(s.apiLink + "/api/product/" + id.String())
	if err != nil {
		return err
	}

	return s.repo.DeleteProduct(ctx, id)
}
