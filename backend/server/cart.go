package server

import (
	"context"
	"errors"
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

// @Tags cart
// @Summary Add product to cart
// @Accept json
// @Param addRequest body AddToCartReq true "Add something to cart"
// @Success 200
// @Failure 401 {object} error "Unauthorized"
// @Failure 403 {object} error "Forbidden"
// @Failure 404 {object} error "Not Found"
// @Failure 500 {object} error "Internal Server Error"
// @Router /api/cart/add [post]
func (a *Server) AddProductToCart(c fiber.Ctx) error {
	var req AddToCartReq
	req.UserID = models.TestUserId
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

// @Tags cart
// @Summary Delete product from cart
// @Accept json
// @Param deleteRequest body DeleteFromCartReq true "Delete product from cart"
// @Success 200
// @Failure 400 {object} error "Bad request"
// @Failure 401 {object} error "Unauthorized"
// @Failure 403 {object} error "Forbidden"
// @Failure 404 {object} error "Not Found"
// @Failure 500 {object} error "Internal Server Error"
// @Router /api/cart/delete [delete]
func (a *Server) DeleteProductFromCart(c fiber.Ctx) error {
	var req DeleteFromCartReq
	req.UserID = models.TestUserId
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

// @Tags cart
// @Summary Get cart
// @Accept json
// @Param getCartRequest body GetCartReq true "Get cart by user"
// @Success 200 {object} models.Cart
// @Failure 400 {object} error "Bad request"
// @Failure 401 {object} error "Unauthorized"
// @Failure 404 {object} error "Not Found"
// @Failure 500 {object} error "Internal Server Error"
// @Router /api/cart [get]
func (a *Server) GetCart(c fiber.Ctx) error {
	var req GetCartReq
	req.UserID = models.TestUserId
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

// @Tags cart
// @Summary Order cart
// @Accept json
// @Param orderRequest body OrderReq true "Submit an order for the cart"
// @Success 200
// @Failure 400 {object} error "Bad request"
// @Failure 401 {object} error "Unauthorized"
// @Failure 403 {object} error "Forbidden"
// @Failure 404 {object} error "Not Found"
// @Failure 500 {object} error "Internal Server Error"
// @Router /api/cart/order [post]
func (a *Server) Order(c fiber.Ctx) error {
	var req OrderReq
	req.UserID = models.TestUserId
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

func mapError(err error) error {
	switch {
	case errors.Is(err, models.ErrCartNotFound):
		return fiber.NewError(404, "Cart not found")
	case errors.Is(err, models.ErrCartNotYours):
		return fiber.NewError(401, "Cart not yours or already finished")
	case errors.Is(err, models.ErrCartInProcess):
		return fiber.NewError(403, "Cart in process, can't change it")
	default:
		return fiber.NewError(fiber.StatusInternalServerError, "unexpected internal error")
	}
}
