package category

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, name string) (*entity.Category, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	ListCategories(ctx context.Context, page, pageSize int) ([]*entity.Category, int, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, name string) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error

	// Product-Category relationship operations
	AssignCategoryToProduct(ctx context.Context, productID, categoryID uuid.UUID) error
	RemoveCategoryFromProduct(ctx context.Context, productID, categoryID uuid.UUID) error
	GetProductCategories(ctx context.Context, productID uuid.UUID) ([]*entity.Category, error)
}

type UseCase struct {
	repo repository.CategoryRepository
}

func NewUseCase(repo repository.CategoryRepository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (uc *UseCase) CreateCategory(ctx context.Context, name string) (*entity.Category, error) {
	category := &entity.Category{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := category.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (uc *UseCase) GetCategory(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) ListCategories(ctx context.Context, page, pageSize int) ([]*entity.Category, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return uc.repo.GetAll(ctx, page, pageSize)
}

func (uc *UseCase) UpdateCategory(ctx context.Context, id uuid.UUID, name string) (*entity.Category, error) {
	category, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Name = name
	category.UpdatedAt = time.Now()

	if err := category.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (uc *UseCase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UseCase) AssignCategoryToProduct(ctx context.Context, productID, categoryID uuid.UUID) error {
	return uc.repo.AssignCategoryToProduct(ctx, productID, categoryID)
}

func (uc *UseCase) RemoveCategoryFromProduct(ctx context.Context, productID, categoryID uuid.UUID) error {
	return uc.repo.RemoveCategoryFromProduct(ctx, productID, categoryID)
}

func (uc *UseCase) GetProductCategories(ctx context.Context, productID uuid.UUID) ([]*entity.Category, error) {
	return uc.repo.GetProductCategories(ctx, productID)
}
