package productvariant

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
)

// ProductVariantService defines the interface for product variant operations
type ProductVariantService interface {
	CreateProductVariant(ctx context.Context, productID uuid.UUID, variantName, variantValue string, priceOverride *float64, quantity int) (*entity.ProductVariant, error)
	GetProductVariant(ctx context.Context, id uuid.UUID) (*entity.ProductVariant, error)
	ListProductVariants(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]*entity.ProductVariant, int, error)
	UpdateProductVariant(ctx context.Context, id uuid.UUID, variantName, variantValue string, priceOverride *float64, quantity int) (*entity.ProductVariant, error)
	DeleteProductVariant(ctx context.Context, id uuid.UUID) error
}

type UseCase struct {
	repo repository.ProductVariantRepository
}

func NewUseCase(repo repository.ProductVariantRepository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (uc *UseCase) CreateProductVariant(ctx context.Context, productID uuid.UUID, variantName, variantValue string, priceOverride *float64, quantity int) (*entity.ProductVariant, error) {
	productVariant := &entity.ProductVariant{
		ID:             uuid.New(),
		ProductID:      productID,
		VariantName:    variantName,
		VariantValue:   variantValue,
		Price_Override: priceOverride,
		Quantity:       quantity,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := productVariant.ValidateForCreation(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, productVariant); err != nil {
		return nil, err
	}

	return productVariant, nil
}

func (uc *UseCase) GetProductVariant(ctx context.Context, id uuid.UUID) (*entity.ProductVariant, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) ListProductVariants(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]*entity.ProductVariant, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return uc.repo.GetAllByProductID(ctx, productID, page, pageSize)
}

func (uc *UseCase) UpdateProductVariant(ctx context.Context, id uuid.UUID, variantName, variantValue string, priceOverride *float64, quantity int) (*entity.ProductVariant, error) {
	variant, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	variant.VariantName = variantName
	variant.VariantValue = variantValue
	variant.Price_Override = priceOverride
	variant.Quantity = quantity
	variant.UpdatedAt = time.Now()

	if err := variant.ValidateForCreation(); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, variant); err != nil {
		return nil, err
	}

	return variant, nil
}

func (uc *UseCase) DeleteProductVariant(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
