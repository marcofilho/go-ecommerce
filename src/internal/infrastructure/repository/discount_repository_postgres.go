package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"gorm.io/gorm"
)

type DiscountRepositoryPostgres struct {
	db *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) *DiscountRepositoryPostgres {
	return &DiscountRepositoryPostgres{db: db}
}

func (r *DiscountRepositoryPostgres) Create(ctx context.Context, discount *entity.Discount) error {
	return r.db.WithContext(ctx).Create(discount).Error
}

func (r *DiscountRepositoryPostgres) Update(ctx context.Context, discount *entity.Discount) error {
	return r.db.WithContext(ctx).Save(discount).Error
}

func (r *DiscountRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.Discount, error) {
	var discount entity.Discount
	err := r.db.WithContext(ctx).First(&discount, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &discount, nil
}

func (r *DiscountRepositoryPostgres) GetByPromoCode(ctx context.Context, promoCode string) (*entity.Discount, error) {
	var discount entity.Discount
	err := r.db.WithContext(ctx).Where("promo_code = ? AND active = ?", promoCode, true).First(&discount).Error
	if err != nil {
		return nil, err
	}
	return &discount, nil
}
