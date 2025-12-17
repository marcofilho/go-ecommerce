package discount

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
)

type DiscountService interface {
	CreateDiscount(ctx context.Context, discount *entity.Discount, productIDs, categoryIDs, userIDs []uuid.UUID, userUsageLimit *int) error
	UpdateDiscount(ctx context.Context, discount *entity.Discount, productIDs, categoryIDs, userIDs []uuid.UUID, userUsageLimit *int) error
	GetDiscountByID(ctx context.Context, id uuid.UUID) (*entity.Discount, error)
	GetDiscountByPromoCode(ctx context.Context, promoCode string) (*entity.Discount, error)
	ValidateDiscountForOrder(ctx context.Context, promoCode string, userID uuid.UUID, orderItems []OrderItemInfo, orderTotal float64) (*DiscountValidationResult, error)
	ApplyDiscount(ctx context.Context, discountID, userID uuid.UUID) error
}

// OrderItemInfo contains information about an item in an order for discount validation
type OrderItemInfo struct {
	ProductID   uuid.UUID
	CategoryIDs []uuid.UUID
	Quantity    int
	Price       float64
}

// DiscountValidationResult contains the result of validating a discount
type DiscountValidationResult struct {
	Valid             bool
	Discount          *entity.Discount
	ApplicableItems   []uuid.UUID
	EstimatedDiscount float64
	Message           string
}

type UseCase struct {
	repo        repository.DiscountRepository
	productRepo repository.ProductRepository
}

func NewUseCase(repo repository.DiscountRepository, productRepo repository.ProductRepository) *UseCase {
	return &UseCase{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (uc *UseCase) CreateDiscount(ctx context.Context, discount *entity.Discount, productIDs, categoryIDs, userIDs []uuid.UUID, userUsageLimit *int) error {
	if err := discount.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, discount); err != nil {
		return err
	}

	// Associate products, categories, and users
	if len(productIDs) > 0 {
		if err := uc.repo.AssociateProducts(ctx, discount.ID, productIDs); err != nil {
			return err
		}
	}

	if len(categoryIDs) > 0 {
		if err := uc.repo.AssociateCategories(ctx, discount.ID, categoryIDs); err != nil {
			return err
		}
	}

	if len(userIDs) > 0 {
		if err := uc.repo.AssociateUsers(ctx, discount.ID, userIDs, userUsageLimit); err != nil {
			return err
		}
	}

	return nil
}

func (uc *UseCase) UpdateDiscount(ctx context.Context, discount *entity.Discount, productIDs, categoryIDs, userIDs []uuid.UUID, userUsageLimit *int) error {
	if err := discount.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Update(ctx, discount); err != nil {
		return err
	}

	// Update associations
	if err := uc.repo.AssociateProducts(ctx, discount.ID, productIDs); err != nil {
		return err
	}

	if err := uc.repo.AssociateCategories(ctx, discount.ID, categoryIDs); err != nil {
		return err
	}

	if err := uc.repo.AssociateUsers(ctx, discount.ID, userIDs, userUsageLimit); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetDiscountByID(ctx context.Context, id uuid.UUID) (*entity.Discount, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetDiscountByPromoCode(ctx context.Context, promoCode string) (*entity.Discount, error) {
	return uc.repo.GetByPromoCode(ctx, promoCode)
}

// ValidateDiscountForOrder validates if a discount can be applied to an order
func (uc *UseCase) ValidateDiscountForOrder(ctx context.Context, promoCode string, userID uuid.UUID, orderItems []OrderItemInfo, orderTotal float64) (*DiscountValidationResult, error) {
	// Get discount with all relationships
	discount, err := uc.repo.GetByPromoCodeWithRelations(ctx, promoCode)
	if err != nil {
		return &DiscountValidationResult{
			Valid:   false,
			Message: "Invalid promo code",
		}, nil
	}

	// Check if discount is valid (active, dates, global usage limits)
	if !discount.IsValid() {
		return &DiscountValidationResult{
			Valid:    false,
			Discount: discount,
			Message:  "Discount is not valid or has expired",
		}, nil
	}

	// Check if discount is restricted to specific users
	if len(discount.Users) > 0 {
		isAuthorizedUser := false
		var userUsage *entity.DiscountUser

		for _, u := range discount.Users {
			if u.ID == userID {
				isAuthorizedUser = true
				// Get user-specific usage info
				userUsage, _ = uc.repo.GetUserUsage(ctx, discount.ID, userID)
				break
			}
		}

		if !isAuthorizedUser {
			return &DiscountValidationResult{
				Valid:    false,
				Discount: discount,
				Message:  "This discount is not available for your account",
			}, nil
		}

		// Check user-specific usage limits
		if userUsage != nil {
			if !discount.CanBeUsedBy(userUsage.UsageCount, userUsage.UsageLimit) {
				return &DiscountValidationResult{
					Valid:    false,
					Discount: discount,
					Message:  "You have reached the usage limit for this discount",
				}, nil
			}
		}
	}

	// Check if discount applies to at least one item in the order
	applicableItems := []uuid.UUID{}
	applicableTotal := 0.0

	for _, item := range orderItems {
		if discount.AppliesTo(item.ProductID, item.CategoryIDs) {
			applicableItems = append(applicableItems, item.ProductID)
			applicableTotal += item.Price * float64(item.Quantity)
		}
	}

	if len(applicableItems) == 0 {
		return &DiscountValidationResult{
			Valid:    false,
			Discount: discount,
			Message:  "Discount does not apply to any items in your order",
		}, nil
	}

	// Calculate estimated discount
	// If discount applies to specific items, use applicable total
	// Otherwise use full order total
	totalToDiscount := orderTotal
	if len(discount.Products) > 0 || len(discount.Categories) > 0 {
		totalToDiscount = applicableTotal
	}

	discountedTotal, err := discount.ApplyDiscount(totalToDiscount)
	if err != nil {
		return &DiscountValidationResult{
			Valid:    false,
			Discount: discount,
			Message:  err.Error(),
		}, nil
	}

	estimatedDiscount := totalToDiscount - discountedTotal

	return &DiscountValidationResult{
		Valid:             true,
		Discount:          discount,
		ApplicableItems:   applicableItems,
		EstimatedDiscount: estimatedDiscount,
		Message:           "Discount is valid and can be applied",
	}, nil
}

// ApplyDiscount increments usage counters when a discount is successfully applied
func (uc *UseCase) ApplyDiscount(ctx context.Context, discountID, userID uuid.UUID) error {
	// Increment global usage
	if err := uc.repo.IncrementUsage(ctx, discountID); err != nil {
		return err
	}

	// Increment user-specific usage
	if err := uc.repo.IncrementUserUsage(ctx, discountID, userID); err != nil {
		return err
	}

	return nil
}
