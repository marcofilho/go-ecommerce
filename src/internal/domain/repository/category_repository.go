package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	GetAll(ctx context.Context, page, pageSize int) ([]*entity.Category, int, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByName(ctx context.Context, name string) (*entity.Category, error)

	// Product-Category relationship methods
	AssignCategoryToProduct(ctx context.Context, productID, categoryID uuid.UUID) error
	RemoveCategoryFromProduct(ctx context.Context, productID, categoryID uuid.UUID) error
	GetProductCategories(ctx context.Context, productID uuid.UUID) ([]*entity.Category, error)
}
