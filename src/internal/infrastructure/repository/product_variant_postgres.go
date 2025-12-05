package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	"gorm.io/gorm"
)

type ProductVariantRepositoryPostgres struct {
	db *gorm.DB
}

func NewProductVariantRepositoryPostgres(db *gorm.DB) repository.ProductVariantRepository {
	return &ProductVariantRepositoryPostgres{
		db: db,
	}
}

func (r *ProductVariantRepositoryPostgres) Create(ctx context.Context, productVariant *entity.ProductVariant) error {
	return r.db.WithContext(ctx).Create(productVariant).Error
}

func (r *ProductVariantRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductVariant, error) {
	var productVariant entity.ProductVariant
	err := r.db.WithContext(ctx).Preload("Product").First(&productVariant, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Product variant not found")
		}
		return nil, err
	}

	return &productVariant, nil
}

func (r *ProductVariantRepositoryPostgres) GetAll(ctx context.Context, page, pageSize int) ([]*entity.ProductVariant, int, error) {
	var productVariants []*entity.ProductVariant
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ProductVariant{}).Preload("Product")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&productVariants).Error

	if err != nil {
		return nil, 0, err
	}

	return productVariants, int(total), nil
}

func (r *ProductVariantRepositoryPostgres) GetAllByProductID(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]*entity.ProductVariant, int, error) {
	var productVariants []*entity.ProductVariant
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ProductVariant{}).Preload("Product").Where("product_id = ?", productID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&productVariants).Error

	if err != nil {
		return nil, 0, err
	}

	return productVariants, int(total), nil
}

func (r *ProductVariantRepositoryPostgres) Update(ctx context.Context, productVariant *entity.ProductVariant) error {
	result := r.db.WithContext(ctx).Save(productVariant)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("Product variant not found")
	}

	return nil
}

func (r *ProductVariantRepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entity.ProductVariant{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("Product variant not found")
	}

	return nil
}
