package repo

import (
	"context"
	"errors"

	"backend/models"
	rdb "backend/rdb/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CartRepo struct {
	q  *rdb.Queries
	db *pgx.Conn
}

func NewCartRepo(q *rdb.Queries) CartRepo {
	return CartRepo{q: q}
}

func (c *CartRepo) GetCart(ctx context.Context, userID uuid.UUID) (models.Cart, error) {
	cart, err := c.q.GetCart(ctx, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return models.Cart{}, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		cart, err := c.q.CreateCart(ctx, userID)
		if err != nil {
			return models.Cart{}, err
		}
		return models.Cart{ID: cart.ID, UserID: cart.UserID}, nil
	}

	cartProducts, err := c.q.GetCartProducts(ctx, cart.ID)
	if err != nil {
		return models.Cart{}, err
	}

	products := make([]models.Product, len(cartProducts))
	for i, cp := range cartProducts {
		products[i] = models.Product{ID: cp.ID, Title: cp.Title, Price: cp.Price, MarketID: cp.MarketID}
	}

	return models.Cart{ID: cart.ID, UserID: cart.UserID, Products: products}, nil
}

func (c *CartRepo) AddProduct(ctx context.Context, cartID, userID, productID uuid.UUID) error {
	cart, err := c.q.GetCart(ctx, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrCartNotFound
	}

	if cart.ID != cartID {
		return models.ErrCartNotYours
	}

	if cart.InProcess.Bool {
		return models.ErrCartInProcess
	}

	return c.q.AddToCart(ctx, rdb.AddToCartParams{
		CartID:    cartID,
		ProductID: productID,
	})
}

func (c *CartRepo) DeleteProduct(ctx context.Context, cartID, userID, productID uuid.UUID) error {
	cart, err := c.q.GetCart(ctx, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrCartNotFound
	}

	if cart.ID != cartID {
		return models.ErrCartNotYours
	}

	if cart.InProcess.Bool {
		return models.ErrCartInProcess
	}

	return c.q.DeleteFromCart(ctx, rdb.DeleteFromCartParams{CartID: cartID, ProductID: productID})
}

func (c *CartRepo) Order(ctx context.Context, cartID, userID uuid.UUID) (models.Order, error) {
	cart, err := c.q.GetCart(ctx, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return models.Order{}, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Order{}, models.ErrCartNotFound
	}

	if cart.ID != cartID {
		return models.Order{}, models.ErrCartNotYours
	}

	if cart.InProcess.Bool {
		return models.Order{}, models.ErrCartInProcess
	}

	tx, err := c.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return models.Order{}, err
	}
	defer tx.Rollback(ctx)
	err = c.q.OrderStartProcess(ctx, cartID)
	if err != nil {
		return models.Order{}, err
	}
	reqProducts, err := c.q.GetProductsForOrder(ctx, cartID)
	if err != nil {
		return models.Order{}, err
	}

	products := make(map[string][]uuid.UUID)
	for _, rp := range reqProducts {
		products[rp.ApiLink] = append(products[rp.ApiLink], rp.ProductID)
	}

	return models.Order{ID: cartID, Products: products}, nil
}

func (c *CartRepo) FinishOrder(ctx context.Context, orderID uuid.UUID) error {
	return c.q.FinishOrder(ctx, orderID)
}
