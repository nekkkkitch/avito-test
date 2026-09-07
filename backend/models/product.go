package models

import "github.com/google/uuid"

type Product struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	MarketID uuid.UUID `json:"market_id" validate:"required"`
	Name     string    `json:"name" validate:"required"`
	Price    int       `json:"price" validate:"required,min=1,max=99999999"`
}
