package repo

import (
	"context"

	"backend/models"
	rdb "backend/rdb/postgres"
	"github.com/google/uuid"
)

type ProductsRepo struct {
	q *rdb.Queries
}

func NewProductsRepo(q *rdb.Queries) ProductsRepo {
	return ProductsRepo{q: q}
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

func (p *ProductsRepo) GetProductCard(ctx context.Context, id uuid.UUID) (models.Product, error) {
	res, err := p.q.GetProductCard(ctx, id)
	if err != nil {
		return models.Product{}, err
	}
	return models.Product{ID: res.ID, MarketID: res.MarketID, Title: res.Title, Price: res.Price}, nil
}
