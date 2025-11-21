package repository

import (
	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

type ProductRepository interface {
	Create(product *entity.Product) error
	GetByID(id uuid.UUID) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id uuid.UUID) error
	GetAll() ([]*entity.Product, error)
}
