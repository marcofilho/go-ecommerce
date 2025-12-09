package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"gorm.io/gorm"
)

type CategoryRepositoryPostgres struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepositoryPostgres {
	return &CategoryRepositoryPostgres{db: db}
}

func (r *CategoryRepositoryPostgres) Create(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *CategoryRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	var category entity.Category
	err := r.db.WithContext(ctx).Preload("Products").First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepositoryPostgres) GetAll(ctx context.Context, page, pageSize int) ([]*entity.Category, int, error) {
	var categories []*entity.Category
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&entity.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Order("name ASC").
		Find(&categories).Error

	if err != nil {
		return nil, 0, err
	}

	return categories, int(total), nil
}

func (r *CategoryRepositoryPostgres) Update(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *CategoryRepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Category{}, "id = ?", id).Error
}

func (r *CategoryRepositoryPostgres) GetByName(ctx context.Context, name string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepositoryPostgres) AssignCategoryToProduct(ctx context.Context, productID, categoryID uuid.UUID) error {
	// Get product and category to ensure they exist
	var product entity.Product
	if err := r.db.WithContext(ctx).First(&product, "id = ?", productID).Error; err != nil {
		return err
	}

	var category entity.Category
	if err := r.db.WithContext(ctx).First(&category, "id = ?", categoryID).Error; err != nil {
		return err
	}

	// Add the association
	return r.db.WithContext(ctx).Model(&product).Association("Categories").Append(&category)
}

func (r *CategoryRepositoryPostgres) RemoveCategoryFromProduct(ctx context.Context, productID, categoryID uuid.UUID) error {
	var product entity.Product
	if err := r.db.WithContext(ctx).First(&product, "id = ?", productID).Error; err != nil {
		return err
	}

	var category entity.Category
	if err := r.db.WithContext(ctx).First(&category, "id = ?", categoryID).Error; err != nil {
		return err
	}

	// Remove the association
	return r.db.WithContext(ctx).Model(&product).Association("Categories").Delete(&category)
}

func (r *CategoryRepositoryPostgres) GetProductCategories(ctx context.Context, productID uuid.UUID) ([]*entity.Category, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).Preload("Categories").First(&product, "id = ?", productID).Error
	if err != nil {
		return nil, err
	}

	return convertCategoriesToPointers(product.Categories), nil
}

func convertCategoriesToPointers(categories []entity.Category) []*entity.Category {
	result := make([]*entity.Category, len(categories))
	for i := range categories {
		result[i] = &categories[i]
	}
	return result
}
