package repo

import (
	"backend/models"
	rdb "backend/rdb/postgres"
	"context"

	"github.com/google/uuid"
)

type ExternalRepo struct {
	q *rdb.Queries
}

func NewExternalRepo(q *rdb.Queries) ExternalRepo {
	return ExternalRepo{q}
}

func (e *ExternalRepo) SetProduct(ctx context.Context, products []models.Product) error {
	for _, p := range products {
		err := e.q.SetProduct(ctx, rdb.SetProductParams{ID: p.ID, MarketID: p.MarketID, Title: p.Name, Price: p.Price})
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *ExternalRepo) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return e.q.DeleteProduct(ctx, id)
}
