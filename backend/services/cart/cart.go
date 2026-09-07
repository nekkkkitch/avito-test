package cart

import (
	"backend/models"
	"fmt"

	"github.com/gofiber/fiber/v3/client"
	"github.com/google/uuid"
)

type RDBCart interface {
	AddProduct(userID, productID uuid.UUID) error
	DeleteProduct(userID, productID uuid.UUID) error
	Order(userID uuid.UUID) (models.Order, error)
	FinishOrder(orderID uuid.UUID) error
}

type Service struct {
	cart RDBCart
}

func New(c RDBCart) *Service {
	return &Service{
		cart: c,
	}
}

func (s *Service) AddProductToCart(userID, productID uuid.UUID) error {
	return s.cart.AddProduct(userID, productID)
}

func (s *Service) DeleteProductFromCart(userID, productID uuid.UUID) error {
	return s.cart.DeleteProduct(userID, productID)
}

type Order struct {
	ID       uuid.UUID   `json:"id"`
	Products []uuid.UUID `json:"products"`
}

func (s *Service) Order(userID uuid.UUID) error {
	cart, err := s.cart.Order(userID)
	if err != nil {
		return err
	}

	for k := range cart.Products {
		miniOrder := models.SingleOrder{
			ID:       cart.ID,
			UserID:   userID,
			Products: cart.Products[k],
		}
		_, err := client.Post(fmt.Sprint("http://"+k+"/order"), client.Config{
			Body: miniOrder})
		if err != nil {
			return err
		}
	}

	err = s.cart.FinishOrder(userID)
	if err != nil {
		return err
	}

	return nil
}
