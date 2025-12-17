package order

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/audit"
)

type CreateOrderItem struct {
	ProductID uuid.UUID
	VariantID *uuid.UUID // Optional: if ordering a specific variant
	Quantity  int
}

type OrderService interface {
	CreateOrder(ctx context.Context, customerID int, items []CreateOrderItem, promoCode string) (*entity.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	ListOrders(ctx context.Context, page, pageSize int, status *entity.OrderStatus, paymentStatus *entity.PaymentStatus) ([]*entity.Order, int, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, newStatus entity.OrderStatus) (*entity.Order, error)
}

type Services interface {
	GetAuditService() audit.AuditService
}

type UseCase struct {
	orderRepo    repository.OrderRepository
	productRepo  repository.ProductRepository
	variantRepo  repository.ProductVariantRepository
	discountRepo repository.DiscountRepository
	services     Services
}

func NewUseCase(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, variantRepo repository.ProductVariantRepository, discountRepo repository.DiscountRepository, services Services) *UseCase {
	return &UseCase{
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		variantRepo:  variantRepo,
		discountRepo: discountRepo,
		services:     services,
	}
}

func (uc *UseCase) CreateOrder(ctx context.Context, customerID int, items []CreateOrderItem, promoCode string) (*entity.Order, error) {
	if customerID <= 0 {
		return nil, errors.New("Invalid customer ID")
	}

	if len(items) == 0 {
		return nil, errors.New("Order must have at least one item")
	}

	var orderItems []entity.OrderItem
	for _, item := range items {
		// Check if ordering a specific variant
		if item.VariantID != nil {
			// Order with variant: decrement variant stock
			variant, err := uc.variantRepo.GetByID(ctx, *item.VariantID)
			if err != nil {
				return nil, errors.New("Product variant not found: " + item.VariantID.String())
			}

			// Verify variant belongs to the specified product
			if variant.ProductID != item.ProductID {
				return nil, errors.New("Variant does not belong to the specified product")
			}

			if !variant.IsAvailable(item.Quantity) {
				return nil, errors.New("Insufficient stock for product variant")
			}

			// Get price from variant (uses override or base product price)
			price, err := variant.GetPrice()
			if err != nil {
				return nil, err
			}

			orderItem := entity.OrderItem{
				ID:        uuid.New(),
				ProductID: item.ProductID,
				VariantID: item.VariantID,
				Quantity:  item.Quantity,
				Price:     price,
			}

			orderItem.CalculateTotal()

			if err := orderItem.Validate(); err != nil {
				return nil, err
			}

			orderItems = append(orderItems, orderItem)

			// Decrease variant stock
			if err := variant.DecreaseStock(item.Quantity); err != nil {
				return nil, err
			}

			if err := uc.variantRepo.Update(ctx, variant); err != nil {
				return nil, err
			}
		} else {
			// Order without variant: decrement base product stock
			product, err := uc.productRepo.GetByID(ctx, item.ProductID)
			if err != nil {
				return nil, errors.New("Product not found: " + item.ProductID.String())
			}

			if !product.IsAvailable(item.Quantity) {
				return nil, errors.New("Insufficient stock for product: " + product.Name)
			}

			orderItem := entity.OrderItem{
				ID:        uuid.New(),
				ProductID: product.ID,
				VariantID: nil,
				Quantity:  item.Quantity,
				Price:     product.Price,
			}

			orderItem.CalculateTotal()

			if err := orderItem.Validate(); err != nil {
				return nil, err
			}

			orderItems = append(orderItems, orderItem)

			// Decrease base product stock
			if err := product.DecreaseStock(item.Quantity); err != nil {
				return nil, err
			}

			if err := uc.productRepo.Update(ctx, product); err != nil {
				return nil, err
			}
		}
	}

	order := &entity.Order{
		ID:            uuid.New(),
		CustomerID:    customerID,
		Products:      orderItems,
		Status:        entity.Pending,
		PaymentStatus: entity.Unpaid,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Calculate order total
	order.CalculateTotal()

	// Apply discount if promo code provided (with enhanced validation)
	if promoCode != "" {
		// Get discount with all relationships (products, categories, users)
		discount, err := uc.discountRepo.GetByPromoCodeWithRelations(ctx, promoCode)
		if err != nil {
			return nil, errors.New("invalid or inactive promo code")
		}

		// Validate discount is active and not expired
		if !discount.IsValid() {
			return nil, errors.New("discount is not valid, expired, or has reached usage limit")
		}

		// Check if user is authorized (if discount is user-specific)
		userID := uuid.MustParse("00000000-0000-0000-0000-000000000000") // TODO: Get from auth context
		if len(discount.Users) > 0 {
			isAuthorized := false
			var userUsage *entity.DiscountUser

			for _, u := range discount.Users {
				if u.ID == userID {
					isAuthorized = true
					// Check per-user usage limits
					userUsage, _ = uc.discountRepo.GetUserUsage(ctx, discount.ID, userID)
					break
				}
			}

			if !isAuthorized {
				return nil, errors.New("this discount is not available for your account")
			}

			if userUsage != nil && !discount.CanBeUsedBy(userUsage.UsageCount, userUsage.UsageLimit) {
				return nil, errors.New("you have reached the usage limit for this discount")
			}
		}

		// Check which items the discount applies to
		applicableTotal := 0.0
		applicableItems := []string{}

		for _, item := range orderItems {
			// Get product with categories
			product, err := uc.productRepo.GetByID(ctx, item.ProductID)
			if err != nil {
				continue
			}

			// Extract category IDs
			categoryIDs := make([]uuid.UUID, 0, len(product.Categories))
			for _, cat := range product.Categories {
				categoryIDs = append(categoryIDs, cat.ID)
			}

			// Check if discount applies to this product
			if discount.AppliesTo(item.ProductID, categoryIDs) {
				applicableTotal += item.TotalPrice
				applicableItems = append(applicableItems, item.ProductID.String())
			}
		}

		// If discount is product/category specific, check if any items match
		if len(discount.Products) > 0 || len(discount.Categories) > 0 {
			if len(applicableItems) == 0 {
				return nil, errors.New("discount does not apply to any items in your order")
			}
			// Apply discount only to applicable items
			discountedTotal, err := discount.ApplyDiscount(applicableTotal)
			if err != nil {
				return nil, err
			}
			discountAmount := applicableTotal - discountedTotal
			order.TotalPrice -= discountAmount
		} else {
			// Site-wide discount: apply to entire order
			discountedTotal, err := discount.ApplyDiscount(order.TotalPrice)
			if err != nil {
				return nil, err
			}
			order.TotalPrice = discountedTotal
		}

		// Increment usage counters
		if err := uc.discountRepo.IncrementUsage(ctx, discount.ID); err != nil {
			return nil, err
		}
		if err := uc.discountRepo.IncrementUserUsage(ctx, discount.ID, userID); err != nil {
			// Log but don't fail the order
		}
	}

	if err := order.Validate(); err != nil {
		return nil, err
	}

	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (uc *UseCase) GetOrder(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	return uc.orderRepo.GetByID(ctx, id)
}

func (uc *UseCase) ListOrders(ctx context.Context, page, pageSize int, status *entity.OrderStatus, paymentStatus *entity.PaymentStatus) ([]*entity.Order, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return uc.orderRepo.GetAll(ctx, page, pageSize, status, paymentStatus)
}

func (uc *UseCase) UpdateOrderStatus(ctx context.Context, id uuid.UUID, newStatus entity.OrderStatus) (*entity.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store original state for audit
	originalStatus := order.Status

	if err := order.UpdateStatus(newStatus); err != nil {
		return nil, err
	}

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Log order status update
	uc.services.GetAuditService().LogChange(ctx, nil, "UPDATE_STATUS", "Order", order.ID,
		map[string]interface{}{"status": originalStatus},
		map[string]interface{}{"status": newStatus})

	return order, nil
}
