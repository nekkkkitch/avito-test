package models

import "github.com/google/uuid"

type Product struct {
	ID       uuid.UUID `json:"id"`
	MarketID uuid.UUID `json:"market_id"`
	Name     string    `json:"name"`
	Price    int       `json:"price"`
}
