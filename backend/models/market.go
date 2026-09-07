package models

import "github.com/google/uuid"

type Market struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}
