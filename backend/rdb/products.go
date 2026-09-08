package repo

import (
	"context"
	"log/slog"

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
	slog.Info("ProductsRepo: GetMarkets: got called")
	markets, err := p.q.GetMarkets(ctx)
	if err != nil {
		slog.Error("ProductsRepo: GetMarkets: returned error", "err", err)
		return nil, err
	}

	result := make([]models.Market, len(markets))
	for i, market := range markets {
		result[i] = models.Market{ID: market.ID, Title: market.Title}
	}

	return result, nil
}

func (p *ProductsRepo) FilterProducts(ctx context.Context, filter models.Filter) ([]models.Product, error) {
	slog.Info("ProductsRepo: FilterProducts: got called", "filter", filter)
	products, err := p.q.FilterProducts(ctx, filter)
	if err != nil {
		slog.Error("ProductsRepo: FilterProducts: returned error", "err", err)
		return nil, err
	}
	return products, nil
}

func (p *ProductsRepo) GetProductCard(ctx context.Context, id uuid.UUID) (models.Product, error) {
	slog.Info("ProductsRepo: GetProductCard: got called", "id", id)
	res, err := p.q.GetProductCard(ctx, id)
	if err != nil {
		slog.Error("ProductsRepo: GetProductCard: returned error", "err", err)
		return models.Product{}, err
	}
	return models.Product{ID: res.ID, MarketID: res.MarketID, Title: res.Title, Price: res.Price}, nil
}
