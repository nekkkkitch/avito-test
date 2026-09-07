package rdb

import (
	"backend/models"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (q *Queries) FilterProducts(ctx context.Context, filter models.Filter) ([]models.Product, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	belka := builder.Select("*").From("products")

	if filter.MarketsID != nil {
		belka = belka.Where(squirrel.Eq{"market_id": filter.MarketsID})
	}

	if filter.Title != "" {
		belka = belka.Where(squirrel.Like{"title": "%" + filter.Title + "%"})
	}

	if filter.MaxPrice == 0 {
		filter.MaxPrice = 9999999
	}

	belka = belka.Where("price >= ?", filter.MinPrice)
	belka = belka.Where("price <= ?", filter.MaxPrice)

	belka = belka.Limit(uint64(filter.Limit)).Offset(uint64(filter.Offset) * uint64(filter.Limit))

	sql, args, err := belka.ToSql()
	if err != nil {
		return nil, err
	}

	var rows pgx.Rows
	if len(args) == 0 {
		rows, err = q.db.Query(ctx, sql)
	} else {
		rows, err = q.db.Query(ctx, sql, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0, filter.Limit)
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.MarketID,
			&product.Name,
			&product.Price,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
