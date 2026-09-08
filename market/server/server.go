package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"market/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// @title Market API
// @version 1.0
// @description Market service API
// @host localhost:8081
// @BasePath /api

type ProductService interface {
	ListProducts(ctx context.Context) ([]models.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (models.Product, error)
}

type ExternalService interface {
	SaveProduct(ctx context.Context, product models.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type OrderService interface {
	SaveOrder(ctx context.Context, order models.Order) error
	ListOrders(ctx context.Context) ([]models.Order, error)
}

type Server struct {
	app           *fiber.App
	productsSvc   ProductService
	externalSvc   ExternalService
	ordersSvc     OrderService
	mainServerURL string
}

func New(app *fiber.App, products ProductService, external ExternalService, orders OrderService, mainServerURL string) *Server {
	return &Server{
		app:           app,
		productsSvc:   products,
		externalSvc:   external,
		ordersSvc:     orders,
		mainServerURL: mainServerURL,
	}
}

func (s *Server) SetupRoutes() {
	slog.Info("Server: SetupRoutes: got called")
	api := s.app.Group("/api")

	api.Get("/ping", s.Ping)

	api.Get("/products", s.ListProducts)
	api.Post("/product", s.SaveProduct)
	api.Delete("/product", s.DeleteProduct)

	api.Post("/order", s.SaveOrder)
	api.Get("/orders", s.ListOrders)

	s.setupSwagger()
}

func (s *Server) Ping(c fiber.Ctx) error {
	slog.Info("Server: Ping: got called")
	return c.SendStatus(200)
}

// @Tags products
// @Summary List products
// @Success 200 {object} []models.Product
// @Failure 400 {object} error "Bad request"
// @Failure 502 {object} error "Main server notification failed"
// @Router /products [get]
func (s *Server) ListProducts(c fiber.Ctx) error {
	slog.Info("Server: ListProducts: got called")
	products, err := s.productsSvc.ListProducts(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(products)
}

type saveProductReq struct {
	Product models.Product `json:"product"`
}

// @Tags products
// @Summary Save or update product
// @Accept json
// @Param product body saveProductReq true "Product payload"
// @Success 200
// @Failure 400 {object} error "Bad request"
// @Failure 502 {object} error "Main server notification failed"
// @Router /product [post]
func (s *Server) SaveProduct(c fiber.Ctx) error {
	slog.Info("Server: SaveProduct: got called", "product", c.Body())
	var req saveProductReq
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "bad body")
	}

	p := req.Product
	if p.ID == uuid.Nil {
		p.ID, _ = uuid.NewRandom()
	}

	p.MarketID = models.MarketID

	if p.Title == "" || p.Price <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product")
	}

	if err := s.externalSvc.SaveProduct(c.Context(), p); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := s.notifyMainServerSetProduct([]models.Product{p}); err != nil {
		slog.Error("notify main server set product", "err", err)
		return fiber.NewError(fiber.StatusBadGateway, "main server notification failed")
	}

	return c.SendStatus(http.StatusOK)
}

type deleteProductReq struct {
	ID uuid.UUID `json:"id"`
}

// @Tags products
// @Summary Delete product
// @Accept json
// @Param deleteRequest body deleteProductReq true "Delete product payload"
// @Success 200
// @Failure 400 {object} error "Bad request"
// @Failure 404 {object} error "Not Found"
// @Failure 502 {object} error "Main server notification failed"
// @Router /product [delete]
func (s *Server) DeleteProduct(c fiber.Ctx) error {
	slog.Info("Server: DeleteProduct: got called", "body", c.Body())
	var req deleteProductReq
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "bad body")
	}
	if req.ID == uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	if err := s.externalSvc.DeleteProduct(c.Context(), req.ID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := s.notifyMainServerDeleteProduct(req.ID); err != nil {
		slog.Error("notify main server delete product", "err", err)
		return fiber.NewError(fiber.StatusBadGateway, "main server notification failed")
	}

	return c.SendStatus(http.StatusOK)
}

// @Tags orders
// @Summary Save order
// @Accept json
// @Param order body models.Order true "Order payload"
// @Success 200
// @Failure 400 {object} error "Bad request"
// @Failure 409 {object} error "Conflict"
// @Failure 500 {object} error "Internal Server Error"
// @Router /order [post]
func (s *Server) SaveOrder(c fiber.Ctx) error {
	slog.Info("Server: SaveOrder: got called", "order", c.Body())
	var order models.Order
	if err := c.Bind().Body(&order); err != nil {
		slog.Error("Server: SaveOrder: bind", "err", err)
		return fiber.NewError(fiber.StatusBadRequest, "bad body")
	}
	if order.ID == uuid.Nil || order.UserID == uuid.Nil || len(order.Products) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid order")
	}
	if err := s.ordersSvc.SaveOrder(c.Context(), order); err != nil {
		slog.Error("Server: SaveOrder: save", "err", err)
		if err.Error() == "duplicate order id" {
			slog.Error("Server: SaveOrder: save", "err", err)
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(http.StatusOK)
}

// @Tags orders
// @Summary List orders
// @Success 200 {object} []models.Order
// @Failure 400 {object} error "Bad request"
// @Router /products [get]
func (s *Server) ListOrders(c fiber.Ctx) error {
	slog.Info("Server: ListOrder: got called")

	orders, err := s.ordersSvc.ListOrders(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(orders)
}

func (s *Server) notifyMainServerSetProduct(products []models.Product) error {
	payload, err := json.Marshal(products)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.mainServerURL+"/api/product", "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("main server returned %s", resp.Status)
	}
	return nil
}

func (s *Server) notifyMainServerDeleteProduct(id uuid.UUID) error {
	payload, err := json.Marshal(map[string]string{"id": id.String()})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, s.mainServerURL+"/api/product/"+id.String(), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("main server returned %s", resp.Status)
	}
	return nil
}
