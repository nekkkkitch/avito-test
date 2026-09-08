package models

import (
	"github.com/google/uuid"
)

type Order struct {
	ID       uuid.UUID   `json:"id"`
	UserID   uuid.UUID   `json:"user_id"`
	Products []uuid.UUID `json:"products"`
}
