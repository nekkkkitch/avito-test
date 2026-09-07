package server

import (
	"context"
	"log/slog"

	"backend/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CartSvc interface {
	AddProductToCart(ctx context.Context, cartID, userID, productID uuid.UUID) error
	DeleteProductFromCart(ctx context.Context, cartID, userID, productID uuid.UUID) error
	GetCart(ctx context.Context, userID uuid.UUID) (models.Cart, error)
	Order(ctx context.Context, cartID, userID uuid.UUID) error
	// todo add get history
}

type AddToCartReq struct {
	CartID    uuid.UUID `json:"cart_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}

func (a *Server) AddProductToCart(c fiber.Ctx) error {
	var req AddToCartReq
	if err := c.Bind().Body(&req); err != nil {
		slog.Error("API: AddProductToCart: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	if err := a.val.Struct(req); err != nil {
		slog.Error("API: AddProductToCart: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	err := a.cart.AddProductToCart(c.Context(), req.CartID, req.UserID, req.ProductID)
	if err != nil {
		return mapError(err)
	}

	return c.SendStatus(200)
}

type DeleteFromCartReq struct {
	CartID    uuid.UUID `json:"cart_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}

func (a *Server) DeleteProductFromCart(c fiber.Ctx) error {
	var req DeleteFromCartReq
	if err := c.Bind().Body(&req); err != nil {
		slog.Error("API: DeleteProductFromCart: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	if err := a.val.Struct(req); err != nil {
		slog.Error("API: DeleteProductFromCart: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	err := a.cart.DeleteProductFromCart(c.Context(), req.CartID, req.UserID, req.ProductID)
	if err != nil {
		return mapError(err)
	}
	return c.SendStatus(200)
}

type GetCartReq struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (a *Server) GetCart(c fiber.Ctx) error {
	var req GetCartReq
	if err := c.Bind().Body(&req); err != nil {
		slog.Error("API: GetCart: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	if err := a.val.Struct(req); err != nil {
		slog.Error("API: GetCart: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	cart, err := a.cart.GetCart(c.Context(), req.UserID)
	if err != nil {
		return mapError(err)
	}

	if err := c.JSON(cart); err != nil {
		slog.Error("API: GetCart: can't send response", "err", err)
		return fiber.NewError(fiber.ErrInternalServerError.Code, "can't send response")
	}

	return nil
}

type OrderReq struct {
	CartID uuid.UUID `json:"cart_id" validate:"required"`
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (a *Server) Order(c fiber.Ctx) error {
	var req OrderReq
	if err := c.Bind().Body(&req); err != nil {
		slog.Error("API: Order: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	if err := a.val.Struct(req); err != nil {
		slog.Error("API: Order: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	err := a.cart.Order(c.Context(), req.CartID, req.UserID)
	if err != nil {
		return mapError(err)
	}

	return c.SendStatus(200)
}

func mapError(err error) error { /*
		switch {
		case errors.Is(err, errors.ErrUnsupported):
			return fiber.NewError()
		}*/
	return nil
}
