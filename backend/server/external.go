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

func (a *Server) SetProduct(c fiber.Ctx) error {
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

func (a *Server) DeleteProduct(c fiber.Ctx) error {
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
