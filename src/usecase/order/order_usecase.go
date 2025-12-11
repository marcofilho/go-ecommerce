package order

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
)

type CreateOrderItem struct {
	ProductID uuid.UUID
	VariantID *uuid.UUID // Optional: if ordering a specific variant
	Quantity  int
}

type OrderService interface {
	CreateOrder(ctx context.Context, customerID int, items []CreateOrderItem) (*entity.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	ListOrders(ctx context.Context, page, pageSize int, status *entity.OrderStatus, paymentStatus *entity.PaymentStatus) ([]*entity.Order, int, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, newStatus entity.OrderStatus) (*entity.Order, error)
}

type UseCase struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	variantRepo repository.ProductVariantRepository
}

func NewUseCase(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, variantRepo repository.ProductVariantRepository) *UseCase {
	return &UseCase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		variantRepo: variantRepo,
	}
}

func (uc *UseCase) CreateOrder(ctx context.Context, customerID int, items []CreateOrderItem) (*entity.Order, error) {
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

	order.CalculateTotal()

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

	if err := order.UpdateStatus(newStatus); err != nil {
		return nil, err
	}

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
