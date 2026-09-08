package models

import "github.com/google/uuid"

type Product struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	MarketID uuid.UUID `json:"market_id" validate:"required"`
	Title    string    `json:"title" validate:"required"`
	Price    int       `json:"price" validate:"required,min=1,max=99999999"`
}

type Filter struct {
	MarketsID []uuid.UUID `json:"market_id"`
	Title     string      `json:"title"`
	MinPrice  int         `json:"min_price" validate:"min=1,max=99999999"`
	MaxPrice  int         `json:"max_price" validate:"min=0,max=99999999"`
	Offset    int         `json:"offset" validate:"min=0"`
	Limit     int         `json:"limit" validate:"required,min=1,max=100"`
}
