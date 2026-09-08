package cart

import (
	"context"
	"log/slog"

	"backend/models"

	"github.com/gofiber/fiber/v3/client"
	"github.com/google/uuid"
)

type RDBCart interface {
	GetCart(ctx context.Context, userID uuid.UUID) (models.Cart, error)
	AddProduct(ctx context.Context, cartID, userID, productID uuid.UUID) error
	DeleteProduct(ctx context.Context, cartID, userID, productID uuid.UUID) error
	Order(ctx context.Context, cartID, userID uuid.UUID) (models.Order, error)
	FinishOrder(ctx context.Context, orderID uuid.UUID) error
}

type Service struct {
	cart RDBCart
}

func New(c RDBCart) *Service {
	return &Service{
		cart: c,
	}
}

func (s *Service) GetCart(ctx context.Context, userID uuid.UUID) (models.Cart, error) {
	slog.Info("CartService: GetCart: got called", "userID", userID)
	return s.cart.GetCart(ctx, userID)
}

func (s *Service) AddProductToCart(ctx context.Context, cartID, userID, productID uuid.UUID) error {
	slog.Info("CartService: AddProductToCart: got called", "cartID", cartID, "userID", userID, "productID", productID)
	return s.cart.AddProduct(ctx, cartID, userID, productID)
}

func (s *Service) DeleteProductFromCart(ctx context.Context, cartID, userID, productID uuid.UUID) error {
	slog.Info("CartService: DeleteProductFromCart: got called", "cartID", cartID, "userID", userID, "productID", productID)
	return s.cart.DeleteProduct(ctx, cartID, userID, productID)
}

type Order struct {
	ID       uuid.UUID   `json:"id"`
	Products []uuid.UUID `json:"products"`
}

func (s *Service) Order(ctx context.Context, cartID, userID uuid.UUID) error {
	slog.Info("CartService: Order: got called", "cartID", cartID, "userID", userID)
	cart, err := s.cart.Order(ctx, cartID, userID)
	if err != nil {
		return err
	}

	for k := range cart.Products {
		miniOrder := models.SingleOrder{
			ID:       cart.ID,
			UserID:   userID,
			Products: cart.Products[k],
		}
		resp, err := client.Post("http://"+k+"/api/order", client.Config{
			Body: miniOrder,
		})
		if err != nil && resp.StatusCode() != 409 {
			return err
		}
	}

	err = s.cart.FinishOrder(ctx, cartID)
	if err != nil {
		return err
	}

	return nil
}
