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
}
