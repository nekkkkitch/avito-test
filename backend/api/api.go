package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type API struct {
	app    *fiber.App
	market SMarket
	cart   SCart
	ext    SExternal
	val    *validator.Validate
}

func New(m SMarket, c SCart, e SExternal) *API {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})
	return &API{
		app:    app,
		market: m,
		cart:   c,
		ext:    e,
		val:    validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (a *API) SetupRoutes() {
	api := a.app.Group("/api")

	api.Get("/markets", a.GetMarkets)
	api.Post("/filter", a.FilterProducts)
	api.Get("/product/:id", a.GetProductCard)

	api.Post("/cart/add", a.AddProductToCart)
	api.Delete("/cart/delete", a.DeleteProductFromCart)
	api.Get("/cart", a.GetCart)
	api.Post("/cart/order", a.Order)

	api.Post("/product", a.SetProduct)
	api.Delete("/product/:id", a.DeleteProduct)
}
