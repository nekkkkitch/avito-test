package repo

import (
	"context"
	"log/slog"

	"backend/models"
	rdb "backend/rdb/postgres"
	"github.com/google/uuid"
)

type ExternalRepo struct {
	q *rdb.Queries
}

func NewExternalRepo(q *rdb.Queries) ExternalRepo {
	return ExternalRepo{q}
}

func (e *ExternalRepo) SetProduct(ctx context.Context, products []models.Product) error {
	slog.Info("ExternalRepo: SetProduct: got called", "products", products)
	for _, p := range products {
		err := e.q.SetProduct(ctx, rdb.SetProductParams{ID: p.ID, MarketID: p.MarketID, Title: p.Title, Price: p.Price})
		if err != nil {
			slog.Error("ExternalRepo: SetProduct: returned error", "err", err)
			return err
		}
	}
	return nil
}

func (e *ExternalRepo) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	slog.Info("ExternalRepo: DeleteProduct: got called", "id", id)
	err := e.q.DeleteProduct(ctx, id)
	if err != nil {
		slog.Error("ExternalRepo: DeleteProduct: returned error", "err", err)
		return err
	}
	return nil
}
