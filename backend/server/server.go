package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type Server struct {
	app    *fiber.App
	market ProductsSvc
	cart   CartSvc
	ext    ExternalSvc
	val    *validator.Validate
}

func New(app *fiber.App, m ProductsSvc, c CartSvc, e ExternalSvc) *Server {
	s := &Server{
		app:    app,
		market: m,
		cart:   c,
		ext:    e,
		val:    validator.New(validator.WithRequiredStructEnabled()),
	}
	return s
}

func (a *Server) SetupRoutes() {
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
