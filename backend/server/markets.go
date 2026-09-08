package server

import (
	"context"
	"log/slog"

	"backend/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ProductsSvc interface {
	GetMarkets(ctx context.Context) ([]models.Market, error)
	FilterProducts(ctx context.Context, filter models.Filter) ([]models.Product, error)
	GetProductCard(ctx context.Context, id uuid.UUID) (models.Product, error)
}

// @Tags market
// @Summary Get all markets
// @Success 200 {array} models.Market
// @Failure 500 {object} error "Internal Server Error"
// @Router /api/markets [get]
func (a *Server) GetMarkets(c fiber.Ctx) error {
	slog.Info("Server: GetMarkets: got called")
	markets, err := a.market.GetMarkets(c.Context())
	if err != nil {
		return mapError(err)
	}

	return c.JSON(markets)
}

// @Tags market
// @Summary Filter products
// @Accept json
// @Param filter body models.Filter true "Filter criteria"
// @Success 200 {array} models.Product
// @Failure 400 {object} error "Bad request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /api/filter [post]
func (a *Server) FilterProducts(c fiber.Ctx) error {
	slog.Info("Server: FilterProducts: got called", "body", c.Body())
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

// @Tags market
// @Summary Get product card
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} error "Bad request"
// @Failure 404 {object} error "Not Found"
// @Failure 500 {object} error "Internal Server Error"
// @Router /api/product/{id} [get]
func (a *Server) GetProductCard(c fiber.Ctx) error {
	slog.Info("Server: GetProductCard: got called", "id", c.Params("id"))
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
