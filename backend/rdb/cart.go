package repo

import (
	"context"
	"errors"
	"log/slog"

	"backend/models"
	rdb "backend/rdb/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepo struct {
	q  *rdb.Queries
	db *pgxpool.Pool
}

func NewCartRepo(q *rdb.Queries, db *pgxpool.Pool) CartRepo {
	return CartRepo{q: q, db: db}
}

func (c *CartRepo) GetCart(ctx context.Context, userID uuid.UUID) (models.Cart, error) {
	slog.Info("CartRepo: GetCart: got called", "userID", userID)
	cart, err := c.q.GetCart(ctx, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("CartRepo: GetCart: returned error", "err", err)
		return models.Cart{}, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		cart, err := c.q.CreateCart(ctx, userID)
		if err != nil {
			slog.Error("CartRepo: GetCart: returned error", "err", err)
			return models.Cart{}, err
		}
		return models.Cart{ID: cart.ID, UserID: cart.UserID}, nil
	}

	cartProducts, err := c.q.GetCartProducts(ctx, cart.ID)
	if err != nil {
		slog.Error("CartRepo: GetCart: returned error", "err", err)
		return models.Cart{}, err
	}

	products := make([]models.Product, len(cartProducts))
	for i, cp := range cartProducts {
		products[i] = models.Product{ID: cp.ID, Title: cp.Title, Price: cp.Price, MarketID: cp.MarketID}
	}

	return models.Cart{ID: cart.ID, UserID: cart.UserID, Products: products}, nil
}

func (c *CartRepo) AddProduct(ctx context.Context, cartID, userID, productID uuid.UUID) error {
	slog.Info("CartRepo: AddProduct: got called", "cartID", cartID, "userID", userID, "productID", productID)
	cart, err := c.q.GetCart(ctx, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("CartRepo: AddProduct: returned error", "err", err)
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

	err = c.q.AddToCart(ctx, rdb.AddToCartParams{
		CartID:    cartID,
		ProductID: productID,
	})
	if err != nil {
		slog.Error("CartRepo: AddProduct: returned error", "err", err)
		return err
	}
	return nil
}

func (c *CartRepo) DeleteProduct(ctx context.Context, cartID, userID, productID uuid.UUID) error {
	slog.Info("CartRepo: DeleteProduct: got called", "cartID", cartID, "userID", userID, "productID", productID)
	cart, err := c.q.GetCart(ctx, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("CartRepo: DeleteProduct: returned error", "err", err)
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

	err = c.q.DeleteFromCart(ctx, rdb.DeleteFromCartParams{CartID: cartID, ProductID: productID})
	if err != nil {
		slog.Error("CartRepo: DeleteProduct: returned error", "err", err)
		return err
	}
	return nil
}

func (c *CartRepo) Order(ctx context.Context, cartID, userID uuid.UUID) (models.Order, error) {
	slog.Info("CartRepo: Order: got called", "cartID", cartID, "userID", userID)
	cart, err := c.q.GetCart(ctx, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("CartRepo: Order: returned error", "err", err)
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
		reqProducts, err := c.q.GetProductsForOrder(ctx, cartID)
		if err != nil {
			slog.Error("CartRepo: Order: returned error", "err", err)
			return models.Order{}, err
		}

		products := make(map[string][]uuid.UUID)
		for _, rp := range reqProducts {
			products[rp.ApiLink] = append(products[rp.ApiLink], rp.ProductID)
		}

		return models.Order{ID: cartID, UserID: userID, Products: products}, nil
	}

	defer tx.Rollback(ctx)
	err = c.q.OrderStartProcess(ctx, cartID)
	if err != nil {
		slog.Error("CartRepo: Order: returned error", "err", err)
		return models.Order{}, err
	}
	reqProducts, err := c.q.GetProductsForOrder(ctx, cartID)
	if err != nil {
		slog.Error("CartRepo: Order: returned error", "err", err)
		return models.Order{}, err
	}

	products := make(map[string][]uuid.UUID)
	for _, rp := range reqProducts {
		products[rp.ApiLink] = append(products[rp.ApiLink], rp.ProductID)
	}

	return models.Order{ID: cartID, UserID: userID, Products: products}, nil
}

func (c *CartRepo) FinishOrder(ctx context.Context, orderID uuid.UUID) error {
	slog.Info("CartRepo: FinishOrder: got called", "orderID", orderID)
	err := c.q.FinishOrder(ctx, orderID)
	if err != nil {
		slog.Error("CartRepo: FinishOrder: returned error", "err", err)
		return err
	}
	return nil
}
