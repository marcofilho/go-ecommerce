package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	"gorm.io/gorm"
)

type ProductRepositoryPostgres struct {
	db *gorm.DB
}

func NewProductRepositoryPostgres(db *gorm.DB) repository.ProductRepository {
	return &ProductRepositoryPostgres{
		db: db,
	}
}

func (r *ProductRepositoryPostgres) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *ProductRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).Preload("Categories").Preload("Variants").First(&product, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepositoryPostgres) GetAll(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error) {
	var products []*entity.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Product{})

	if inStockOnly {
		query = query.Where("quantity > ?", 0)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Preload("Categories").Preload("Variants").Offset(offset).Limit(pageSize).Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, int(total), nil
}

func (r *ProductRepositoryPostgres) Update(ctx context.Context, product *entity.Product) error {
	result := r.db.WithContext(ctx).Save(product)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("Product not found")
	}

	return nil
}

func (r *ProductRepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entity.Product{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("Product not found")
	}

	return nil
}
