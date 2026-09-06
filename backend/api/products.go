package api

import (
	"backend/models"

	"github.com/gofiber/fiber/v3"
)

type SMarket interface {
	GetMarkets() ([]models.Market, error)
	FilterProducts(filter models.Filter) ([]models.Product, error)
}

func (a *API) GetMarkets(c *fiber.Ctx) error {

	return nil
}

func (a *API) FilterProducts(c *fiber.Ctx) error {
	return nil
}

func (a *API) GetProductCard(c *fiber.Ctx) error {
	return nil
}
