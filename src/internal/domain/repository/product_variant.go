package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

type ProductVariantRepository interface {
	Create(ctx context.Context, productVariant *entity.ProductVariant) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductVariant, error)
	GetAll(ctx context.Context, page, pageSize int) ([]*entity.ProductVariant, int, error)
	GetAllByProductID(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]*entity.ProductVariant, int, error)
	Update(ctx context.Context, productVariant *entity.ProductVariant) error
	Delete(ctx context.Context, id uuid.UUID) error
}
