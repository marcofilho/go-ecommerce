package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	"gorm.io/gorm"
)

type OrderRepositoryPostgres struct {
	db *gorm.DB
}

func NewOrderRepositoryPostgres(db *gorm.DB) repository.OrderRepository {
	return &OrderRepositoryPostgres{
		db: db,
	}
}

func (r *OrderRepositoryPostgres) Create(ctx context.Context, order *entity.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *OrderRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	err := r.db.WithContext(ctx).Preload("Items").First(&order, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepositoryPostgres) GetAll(ctx context.Context, page, pageSize int, status *entity.OrderStatus, paymentStatus *entity.PaymentStatus) ([]*entity.Order, int, error) {
	var orders []*entity.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Order{})

	// Apply filters
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if paymentStatus != nil {
		query = query.Where("payment_status = ?", *paymentStatus)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and preload items
	offset := (page - 1) * pageSize
	err := query.Preload("Items").Offset(offset).Limit(pageSize).Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}

func (r *OrderRepositoryPostgres) Update(ctx context.Context, order *entity.Order) error {
	result := r.db.WithContext(ctx).Save(order)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil
}
