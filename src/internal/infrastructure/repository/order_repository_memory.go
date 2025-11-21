package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
)

type OrderRepositoryMemory struct {
	orders map[uuid.UUID]*entity.Order
	mu     sync.RWMutex
}

func NewOrderRepositoryMemory() repository.OrderRepository {
	return &OrderRepositoryMemory{
		orders: make(map[uuid.UUID]*entity.Order),
	}
}

func (r *OrderRepositoryMemory) Create(ctx context.Context, order *entity.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; exists {
		return errors.New("order already exists")
	}

	r.orders[order.ID] = order
	return nil
}

func (r *OrderRepositoryMemory) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}

	return order, nil
}

func (r *OrderRepositoryMemory) GetAll(ctx context.Context, page, pageSize int, status *entity.OrderStatus, paymentStatus *entity.PaymentStatus) ([]*entity.Order, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var orders []*entity.Order
	for _, order := range r.orders {
		// Apply filters
		if status != nil && order.Status != *status {
			continue
		}
		if paymentStatus != nil && order.PaymentStatus != *paymentStatus {
			continue
		}
		orders = append(orders, order)
	}

	total := len(orders)

	// Apply pagination
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > total {
		return []*entity.Order{}, total, nil
	}

	if end > total {
		end = total
	}

	return orders[start:end], total, nil
}

func (r *OrderRepositoryMemory) Update(ctx context.Context, order *entity.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; !exists {
		return errors.New("order not found")
	}

	r.orders[order.ID] = order
	return nil
}
