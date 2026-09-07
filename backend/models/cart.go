package models

import "github.com/google/uuid"

type Cart struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`
	Products []Product `json:"products"`
}
