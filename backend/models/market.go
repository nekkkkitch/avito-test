package models

import "github.com/google/uuid"

type Market struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type Filter struct {
	MarketsID []uuid.UUID `json:"market_id"`
	Title     string      `json:"title"`
	MinPrice  int         `json:"min_price"`
	MaxPrice  int         `json:"max_price"`
}
