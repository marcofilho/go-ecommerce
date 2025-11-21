package product

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
)

type UseCase struct {
	repo repository.ProductRepository
}

func NewUseCase(repo repository.ProductRepository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (uc *UseCase) CreateProduct(name, description string, price float64, quantity int) (*entity.Product, error) {
	product := &entity.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *UseCase) GetProduct(id uuid.UUID) (*entity.Product, error) {
	if id == uuid.Nil {
		return nil, errors.New("Product ID is required")
	}
	return uc.repo.GetByID(id)
}

func (uc *UseCase) ListProducts() ([]*entity.Product, error) {
	return uc.repo.GetAll()
}

func (uc *UseCase) UpdateProduct(id uuid.UUID, name, description string, price float64, quantity int) (*entity.Product, error) {
	product, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	product.Name = name
	product.Description = description
	product.Price = price
	product.Quantity = quantity
	product.UpdatedAt = time.Now()

	if err := product.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *UseCase) DeleteProduct(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("Product ID is required")
	}

	return uc.repo.Delete(id)
}

func (uc *UseCase) UpdateStock(id uuid.UUID, quantity int) error {
	product, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}

	if quantity >= 0 {
		if err := product.IncreaseStock(quantity - product.Quantity); err != nil {
			return err
		}
	} else {
		return errors.New("Stock quantity cannot be negative")
	}

	return uc.repo.Update(product)
}
