package main

import (
	"gorm.io/gorm"

	handler "github.com/marcofilho/go-ecommerce/src/internal/adapter/http"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	infraRepo "github.com/marcofilho/go-ecommerce/src/internal/infrastructure/repository"
	orderUseCase "github.com/marcofilho/go-ecommerce/src/usecase/order"
	productUseCase "github.com/marcofilho/go-ecommerce/src/usecase/product"
)

// Container holds all application dependencies
type Container struct {
	DB *gorm.DB

	// Repositories
	ProductRepo repository.ProductRepository
	OrderRepo   repository.OrderRepository

	// Use Cases
	ProductUseCase *productUseCase.UseCase
	OrderUseCase   *orderUseCase.UseCase

	// Handlers
	ProductHandler *handler.ProductHandler
	OrderHandler   *handler.OrderHandler
}

// NewContainer creates and wires up all dependencies
func NewContainer(db *gorm.DB) *Container {
	c := &Container{DB: db}

	// Initialize repositories
	c.ProductRepo = infraRepo.NewProductRepositoryPostgres(db)
	c.OrderRepo = infraRepo.NewOrderRepositoryPostgres(db)

	// Initialize use cases
	c.ProductUseCase = productUseCase.NewUseCase(c.ProductRepo)
	c.OrderUseCase = orderUseCase.NewUseCase(c.OrderRepo, c.ProductRepo)

	// Initialize handlers
	c.ProductHandler = handler.NewProductHandler(c.ProductUseCase)
	c.OrderHandler = handler.NewOrderHandler(c.OrderUseCase)

	return c
}
