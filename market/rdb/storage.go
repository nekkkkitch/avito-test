package repo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"sync"

	"market/models"

	"github.com/google/uuid"
)

var errDuplicateOrder = errors.New("duplicate order id")

type MemoryRepo struct {
	mu       sync.RWMutex
	products map[uuid.UUID]models.Product
	orders   map[uuid.UUID]models.Order
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		products: make(map[uuid.UUID]models.Product),
		orders:   make(map[uuid.UUID]models.Order),
	}
}

func (r *MemoryRepo) ListProducts(ctx context.Context) ([]models.Product, error) {

	slog.Info("MemoryRepo: ListProducts: got called")
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]models.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}
	sort.Slice(products, func(i, j int) bool {
		return products[i].Title < products[j].Title
	})
	return products, nil
}

func (r *MemoryRepo) GetProductByID(ctx context.Context, id uuid.UUID) (models.Product, error) {

	slog.Info("MemoryRepo: GetProductByID: got called", "id", id)
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[id]
	if !ok {
		err := fmt.Errorf("product %s not found", id)
		slog.Error("MemoryRepo: GetProductByID: returned error", "err", err)
		return models.Product{}, err
	}
	return product, nil
}

func (r *MemoryRepo) SaveProduct(ctx context.Context, product models.Product) error {
	slog.Info("MemoryRepo: SaveProduct: got called", "product", product)
	r.mu.Lock()
	defer r.mu.Unlock()

	r.products[product.ID] = product
	return nil
}

func (r *MemoryRepo) DeleteProduct(ctx context.Context, id uuid.UUID) error {

	slog.Info("MemoryRepo: DeleteProduct: got called", "id", id)
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		err := fmt.Errorf("product %s not found", id)
		slog.Error("MemoryRepo: DeleteProduct: returned error", "err", err)
		return err
	}
	delete(r.products, id)
	return nil
}

func (r *MemoryRepo) SaveOrder(ctx context.Context, order models.Order) error {

	slog.Info("MemoryRepo: SaveOrder: got called", "order", order)
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.ID]; ok {
		err := errDuplicateOrder
		slog.Error("MemoryRepo: SaveOrder: returned error", "err", err)
		return err
	}

	r.orders[order.ID] = order
	slog.Info("", "r.orders", r.orders)
	return nil
}

func (r *MemoryRepo) ListOrders(ctx context.Context) ([]models.Order, error) {
	slog.Info("MemoryRepo: ListOrders: got called")
	r.mu.RLock()
	defer r.mu.RUnlock()

	slog.Info("", "r.orders", r.orders)
	orders := make([]models.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}
	return orders, nil
}
