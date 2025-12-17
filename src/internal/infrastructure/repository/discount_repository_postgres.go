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

func (r *DiscountRepositoryPostgres) GetByPromoCodeWithRelations(ctx context.Context, promoCode string) (*entity.Discount, error) {
	var discount entity.Discount
	err := r.db.WithContext(ctx).
		Preload("Products").
		Preload("Categories").
		Preload("Users").
		Where("promo_code = ? AND active = ?", promoCode, true).
		First(&discount).Error
	if err != nil {
		return nil, err
	}
	return &discount, nil
}

func (r *DiscountRepositoryPostgres) GetUserUsage(ctx context.Context, discountID, userID uuid.UUID) (*entity.DiscountUser, error) {
	var usage entity.DiscountUser
	err := r.db.WithContext(ctx).
		Where("discount_id = ? AND user_id = ?", discountID, userID).
		First(&usage).Error
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

func (r *DiscountRepositoryPostgres) IncrementUsage(ctx context.Context, discountID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.Discount{}).
		Where("id = ?", discountID).
		UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

func (r *DiscountRepositoryPostgres) IncrementUserUsage(ctx context.Context, discountID, userID uuid.UUID) error {
	// First, try to increment if record exists
	result := r.db.WithContext(ctx).
		Model(&entity.DiscountUser{}).
		Where("discount_id = ? AND user_id = ?", discountID, userID).
		UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1))

	if result.Error != nil {
		return result.Error
	}

	// If no record was updated, create a new one with usage_count = 1
	if result.RowsAffected == 0 {
		usage := &entity.DiscountUser{
			DiscountID: discountID,
			UserID:     userID,
			UsageCount: 1,
		}
		return r.db.WithContext(ctx).Create(usage).Error
	}

	return nil
}

func (r *DiscountRepositoryPostgres) AssociateProducts(ctx context.Context, discountID uuid.UUID, productIDs []uuid.UUID) error {
	if len(productIDs) == 0 {
		return nil
	}

	var discount entity.Discount
	if err := r.db.WithContext(ctx).First(&discount, "id = ?", discountID).Error; err != nil {
		return err
	}

	var products []entity.Product
	if err := r.db.WithContext(ctx).Find(&products, "id IN ?", productIDs).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&discount).Association("Products").Replace(products)
}

func (r *DiscountRepositoryPostgres) AssociateCategories(ctx context.Context, discountID uuid.UUID, categoryIDs []uuid.UUID) error {
	if len(categoryIDs) == 0 {
		return nil
	}

	var discount entity.Discount
	if err := r.db.WithContext(ctx).First(&discount, "id = ?", discountID).Error; err != nil {
		return err
	}

	var categories []entity.Category
	if err := r.db.WithContext(ctx).Find(&categories, "id IN ?", categoryIDs).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&discount).Association("Categories").Replace(categories)
}

func (r *DiscountRepositoryPostgres) AssociateUsers(ctx context.Context, discountID uuid.UUID, userIDs []uuid.UUID, usageLimit *int) error {
	if len(userIDs) == 0 {
		return nil
	}

	// Clear existing associations
	if err := r.db.WithContext(ctx).Exec("DELETE FROM discount_users WHERE discount_id = ?", discountID).Error; err != nil {
		return err
	}

	// Create new associations
	for _, userID := range userIDs {
		usage := &entity.DiscountUser{
			DiscountID: discountID,
			UserID:     userID,
			UsageLimit: usageLimit,
		}
		if err := r.db.WithContext(ctx).Create(usage).Error; err != nil {
			return err
		}
	}

	return nil
}
