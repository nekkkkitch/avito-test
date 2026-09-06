package api

import (
	"backend/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type SExternal interface {
	SetProduct(models.Product) error
	DeleteProduct(id uuid.UUID) error
}

func (a *API) SetProduct(c *fiber.Ctx) error {
	return nil
}

func (a *API) DeleteProduct(c *fiber.Ctx) error {
	return nil
}
