package server

import (
	"context"
	"fmt"
	"log/slog"

	"backend/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ExternalSvc interface {
	SetProduct(ctx context.Context, products []models.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

// Компаниям должен выдаваться api ключ для идентификации

// @Tags product
// @Summary Set product
// @Accept json
// @Param product body []models.Product true "Set product data"
// @Success 200
// @Failure 400 {object} error "Bad request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /api/product [post]
func (a *Server) SetProduct(c fiber.Ctx) error {
	slog.Info("Server: SetProduct: got called", "body", c.Body())
	var products []models.Product
	if err := c.Bind().Body(&products); err != nil {
		slog.Error("API: SetProduct: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	for i := range products {
		if err := a.val.Struct(products[i]); err != nil {
			slog.Error("API: SetProduct: can't validate product", "err", err)
			return fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("bad body in product %v", products[i].ID))
		}
	}

	if err := a.ext.SetProduct(c.Context(), products); err != nil {
		return mapError(err)
	}

	return c.SendStatus(200)
}

type DeleteProductRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// @Tags product
// @Summary Delete product
// @Accept json
// @Param deleteRequest body DeleteProductRequest true "Delete product by ID"
// @Success 200
// @Failure 400 {object} error "Bad request"
// @Failure 404 {object} error "Not Found"
// @Failure 500 {object} error "Internal Server Error"
// @Router /api/product/{id} [delete]
func (a *Server) DeleteProduct(c fiber.Ctx) error {
	slog.Info("Server: DeleteProduct: got called", "body", c.Body())
	var req DeleteProductRequest
	if err := c.Bind().Body(&req); err != nil {
		slog.Error("API: DeleteProduct: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	if err := a.val.Struct(req); err != nil {
		slog.Error("API: DeleteProduct: can't validate request", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	if err := a.ext.DeleteProduct(c.Context(), req.ID); err != nil {
		return mapError(err)
	}

	return c.SendStatus(200)
}
