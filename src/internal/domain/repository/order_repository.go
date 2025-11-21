package repository

import (
	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

type OrderRepository interface {
	Create(order *entity.Order) error
	GetByID(id uuid.UUID) (*entity.Order, error)
	GetAll() ([]*entity.Order, error)
	Update(order *entity.Order) error
	Delete(id uuid.UUID) error
}
