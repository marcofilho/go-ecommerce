package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
)

type ProductRepositoryMemory struct {
	products map[uuid.UUID]*entity.Product
	mu       sync.RWMutex
}

func NewProductRepositoryMemory() repository.ProductRepository {
	return &ProductRepositoryMemory{
		products: make(map[uuid.UUID]*entity.Product),
	}
}

func (r *ProductRepositoryMemory) Create(ctx context.Context, product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; exists {
		return errors.New("product already exists")
	}

	r.products[product.ID] = product
	return nil
}

func (r *ProductRepositoryMemory) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}

	return product, nil
}

func (r *ProductRepositoryMemory) GetAll(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var products []*entity.Product
	for _, product := range r.products {
		if inStockOnly && product.Quantity <= 0 {
			continue
		}
		products = append(products, product)
	}

	total := len(products)

	// Apply pagination
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > total {
		return []*entity.Product{}, total, nil
	}

	if end > total {
		end = total
	}

	return products[start:end], total, nil
}

func (r *ProductRepositoryMemory) Update(ctx context.Context, product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return errors.New("product not found")
	}

	r.products[product.ID] = product
	return nil
}

func (r *ProductRepositoryMemory) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[id]; !exists {
		return errors.New("product not found")
	}

	delete(r.products, id)
	return nil
}
