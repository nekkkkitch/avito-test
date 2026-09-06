package api

import (
	"backend/models"

	"github.com/gofiber/fiber/v3"
)

type SCart interface {
	AddProductToCart(userID, productID string) error
	DeleteProductFromCart(userID, productID string) error
	GetCart(userID string) (models.Cart, error)
	Order(userID string) error
}

func (a *API) AddProductToCart(c *fiber.Ctx) error {
	return nil
}

func (a *API) DeleteProductFromCart(c *fiber.Ctx) error {
	return nil
}

func (a *API) GetCart(c *fiber.Ctx) error {
	return nil
}

func (a *API) Order(c *fiber.Ctx) error {
	return nil
}
