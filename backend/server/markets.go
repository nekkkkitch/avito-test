package server

import (
	"context"

	"backend/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ProductsSvc interface {
	GetMarkets(ctx context.Context) ([]models.Market, error)
	FilterProducts(ctx context.Context, filter models.Filter) ([]models.Product, error)
	GetProductCard(ctx context.Context, id uuid.UUID) (models.Product, error)
}

func (a *Server) GetMarkets(c fiber.Ctx) error {
	markets, err := a.market.GetMarkets(c.Context())
	if err != nil {
		return mapError(err)
	}

	return c.JSON(markets)
}

func (a *Server) FilterProducts(c fiber.Ctx) error {
	var filter models.Filter
	if err := c.Bind().Body(&filter); err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	if err := a.val.Struct(filter); err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	products, err := a.market.FilterProducts(c.Context(), filter)
	if err != nil {
		return mapError(err)
	}

	return c.JSON(products)
}

func (a *Server) GetProductCard(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad id")
	}

	product, err := a.market.GetProductCard(c.Context(), id)
	if err != nil {
		return mapError(err)
	}

	return c.JSON(product)
}
