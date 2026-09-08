package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Products  []Product  `json:"products"`
	InProcess bool       `json:"in_process"`
	OrderedAt *time.Time `json:"ordered_at"`
}

type Order struct {
	ID       uuid.UUID              `json:"id"`
	UserID   uuid.UUID              `json:"user_id"`
	Products map[string][]uuid.UUID `json:"products"`
}

type SingleOrder struct {
	ID       uuid.UUID   `json:"id"`
	UserID   uuid.UUID   `json:"user_id"`
	Products []uuid.UUID `json:"products"`
}
