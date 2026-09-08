package repo

import (
	"context"
	"errors"
	"fmt"
	"market/models"
	"sort"
	"sync"

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
	_ = ctx
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
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[id]
	if !ok {
		return models.Product{}, fmt.Errorf("product %s not found", id)
	}
	return product, nil
}

func (r *MemoryRepo) SaveProduct(ctx context.Context, product models.Product) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	r.products[product.ID] = product
	return nil
}

func (r *MemoryRepo) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		return fmt.Errorf("product %s not found", id)
	}
	delete(r.products, id)
	return nil
}

func (r *MemoryRepo) SaveOrder(ctx context.Context, order models.Order) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.ID]; ok {
		return errDuplicateOrder
	}

	r.orders[order.ID] = order
	return nil
}
