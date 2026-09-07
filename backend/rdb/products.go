package repo

import (
	"backend/models"
	rdb "backend/rdb/postgres"
	"context"
)

type ProductsRepo struct {
	q *rdb.Queries
}

func NewProductsRepo(q *rdb.Queries) *ProductsRepo {
	return &ProductsRepo{q: q}
}

func (p *ProductsRepo) GetMarkets(ctx context.Context) ([]models.Market, error) {
	markets, err := p.q.GetMarkets(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.Market, len(markets))
	for i, market := range markets {
		result[i] = models.Market{ID: market.ID, Title: market.Title}
	}

	return result, nil
}

func (p *ProductsRepo) FilterProducts(ctx context.Context, filter models.Filter) ([]models.Product, error) {
	return p.q.FilterProducts(ctx, filter)
}
