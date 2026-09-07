package api

import (
	"backend/models"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type SCart interface {
	AddProductToCart(userID, productID uuid.UUID) error
	DeleteProductFromCart(userID, productID uuid.UUID) error
	GetCart(userID uuid.UUID) (models.Cart, error)
	Order(userID uuid.UUID) error
}

type AddToCartReq struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}

func (a *API) AddProductToCart(c fiber.Ctx) error {
	var req AddToCartReq
	if err := c.Bind().Body(&req); err != nil {
		slog.Error("API: AddProductToCart: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	if err := a.val.Struct(req); err != nil {
		slog.Error("API: AddProductToCart: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	err := a.cart.AddProductToCart(req.UserID, req.ProductID)
	if err != nil {
		return mapError(err)
	}

	return c.SendStatus(200)
}

type DeleteFromCartReq struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}

func (a *API) DeleteProductFromCart(c fiber.Ctx) error {
	var req DeleteFromCartReq
	if err := c.Bind().Body(&req); err != nil {
		slog.Error("API: DeleteProductFromCart: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	if err := a.val.Struct(req); err != nil {
		slog.Error("API: DeleteProductFromCart: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	err := a.cart.DeleteProductFromCart(req.UserID, req.ProductID)
	if err != nil {
		return mapError(err)
	}
	return c.SendStatus(200)
}

type GetCartReq struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (a *API) GetCart(c fiber.Ctx) error {
	var req GetCartReq
	if err := c.Bind().Body(&req); err != nil {
		slog.Error("API: GetCart: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	if err := a.val.Struct(req); err != nil {
		slog.Error("API: GetCart: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	cart, err := a.cart.GetCart(req.UserID)
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
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (a *API) Order(c fiber.Ctx) error {
	var req OrderReq
	if err := c.Bind().Body(&req); err != nil {
		slog.Error("API: Order: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	if err := a.val.Struct(req); err != nil {
		slog.Error("API: Order: can't bind body", "err", err)
		return fiber.NewError(fiber.ErrBadRequest.Code, "bad body")
	}

	err := a.cart.Order(req.UserID)
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
