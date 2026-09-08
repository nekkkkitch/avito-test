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

// @title Avito.Kitchen API
// @version 1.0
// @description Market service API
// @host localhost:8080
// @BasePath /

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

func (s *Server) SetupRoutes() {
	api := s.app.Group("/api")

	api.Get("/markets", s.GetMarkets)
	api.Post("/filter", s.FilterProducts)
	api.Get("/product/:id", s.GetProductCard)

	api.Post("/cart/add", s.AddProductToCart)
	api.Delete("/cart/delete", s.DeleteProductFromCart)
	api.Get("/cart", s.GetCart)
	api.Post("/cart/order", s.Order)

	api.Post("/product", s.SetProduct)
	api.Delete("/product/:id", s.DeleteProduct)

	s.setupSwagger()
}
