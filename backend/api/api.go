package api

import "github.com/gofiber/fiber/v3"

type API struct {
	app *fiber.App
}

func New(app *fiber.App) *API {
	return &API{app}
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
