package discount

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
)

type DiscountService interface {
	CreateDiscount(ctx context.Context, discount *entity.Discount) error
	UpdateDiscount(ctx context.Context, discount *entity.Discount) error
	GetDiscountByID(ctx context.Context, id uuid.UUID) (*entity.Discount, error)
	GetDiscountByPromoCode(ctx context.Context, promoCode string) (*entity.Discount, error)
}

type UseCase struct {
	repo repository.DiscountRepository
}

func NewUseCase(repo repository.DiscountRepository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (uc *UseCase) CreateDiscount(ctx context.Context, discount *entity.Discount) error {
	if err := discount.Validate(); err != nil {
		return err
	}

	return uc.repo.Create(ctx, discount)
}

func (uc *UseCase) UpdateDiscount(ctx context.Context, discount *entity.Discount) error {
	if err := discount.Validate(); err != nil {
		return err
	}

	return uc.repo.Update(ctx, discount)
}

func (uc *UseCase) GetDiscountByID(ctx context.Context, id uuid.UUID) (*entity.Discount, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetDiscountByPromoCode(ctx context.Context, promoCode string) (*entity.Discount, error) {
	return uc.repo.GetByPromoCode(ctx, promoCode)
}
