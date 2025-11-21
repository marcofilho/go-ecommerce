package order

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
)

type UseCase struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewUseCase(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
) *UseCase {
	return &UseCase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (uc *UseCase) CreateOrder(customerID int, products []entity.Product) (*entity.Order, error) {
	for i, product := range products {
		product, err := uc.productRepo.GetByID(product.ID)
		if err != nil {
			return nil, errors.New("Product not found: " + product.ID.String())
		}

		if !product.IsAvailable(product.Quantity) {
			return nil, errors.New("Insufficient stock for product: " + product.Name)
		}

		products[i].Price = product.Price
	}

	order := &entity.Order{
		Customer_ID:    customerID,
		Products:       products,
		Status:         entity.Pending,
		Payment_Status: entity.Unpaid,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	order.CalculateTotal()

	if err := order.Validate(); err != nil {
		return nil, err
	}

	for _, item := range products {
		product, _ := uc.productRepo.GetByID(item.ID)
		if err := product.DecreaseStock(item.Quantity); err != nil {
			return nil, err
		}

		if err := uc.productRepo.Update(product); err != nil {
			return nil, err
		}
	}

	if err := uc.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (uc *UseCase) GetOrder(id uuid.UUID) (*entity.Order, error) {
	if id == uuid.Nil {
		return nil, errors.New("order ID is required")
	}

	return uc.orderRepo.GetByID(id)
}

func (uc *UseCase) ListOrders() ([]*entity.Order, error) {
	return uc.orderRepo.GetAll()
}
