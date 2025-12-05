package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
)

// ProductService defines the interface for product operations
type ProductService interface {
	CreateProduct(ctx context.Context, name, description string, price float64, quantity int) (*entity.Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	ListProducts(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, name, description string, price float64, quantity int) (*entity.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type UseCase struct {
	repo repository.ProductRepository
}

func NewUseCase(repo repository.ProductRepository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (uc *UseCase) CreateProduct(ctx context.Context, name, description string, price float64, quantity int) (*entity.Product, error) {
	product := &entity.Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := product.ValidateForCreation(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *UseCase) GetProduct(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) ListProducts(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return uc.repo.GetAll(ctx, page, pageSize, inStockOnly)
}

func (uc *UseCase) UpdateProduct(ctx context.Context, id uuid.UUID, name, description string, price float64, quantity int) (*entity.Product, error) {
	product, err := uc.repo.GetByID(ctx, id)
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

	if err := uc.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *UseCase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
