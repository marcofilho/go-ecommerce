package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

type DiscountRepository interface {
	Create(ctx context.Context, discount *entity.Discount) error
	Update(ctx context.Context, discount *entity.Discount) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Discount, error)
	GetByPromoCode(ctx context.Context, promoCode string) (*entity.Discount, error)
	GetByPromoCodeWithRelations(ctx context.Context, promoCode string) (*entity.Discount, error)
	GetUserUsage(ctx context.Context, discountID, userID uuid.UUID) (*entity.DiscountUser, error)
	IncrementUsage(ctx context.Context, discountID uuid.UUID) error
	IncrementUserUsage(ctx context.Context, discountID, userID uuid.UUID) error
	AssociateProducts(ctx context.Context, discountID uuid.UUID, productIDs []uuid.UUID) error
	AssociateCategories(ctx context.Context, discountID uuid.UUID, categoryIDs []uuid.UUID) error
	AssociateUsers(ctx context.Context, discountID uuid.UUID, userIDs []uuid.UUID, usageLimit *int) error
}
