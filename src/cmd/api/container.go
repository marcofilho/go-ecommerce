package main

import (
	"gorm.io/gorm"

	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/handler"
	"github.com/marcofilho/go-ecommerce/src/internal/config"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	infraRepo "github.com/marcofilho/go-ecommerce/src/internal/infrastructure/repository"
	orderUseCase "github.com/marcofilho/go-ecommerce/src/usecase/order"
	paymentUseCase "github.com/marcofilho/go-ecommerce/src/usecase/payment"
	productUseCase "github.com/marcofilho/go-ecommerce/src/usecase/product"
)

// Container holds all application dependencies
type Container struct {
	DB     *gorm.DB
	Config *config.Config

	// Repositories
	ProductRepo repository.ProductRepository
	OrderRepo   repository.OrderRepository
	WebhookRepo repository.WebhookRepository

	// Use Cases
	ProductUseCase *productUseCase.UseCase
	OrderUseCase   *orderUseCase.UseCase
	PaymentUseCase *paymentUseCase.PaymentUseCase

	// Handlers
	ProductHandler *handler.ProductHandler
	OrderHandler   *handler.OrderHandler
	PaymentHandler *handler.PaymentHandler
}

// NewContainer creates and wires up all dependencies
func NewContainer(db *gorm.DB, cfg *config.Config) *Container {
	c := &Container{
		DB:     db,
		Config: cfg,
	}

	c.ProductRepo = infraRepo.NewProductRepositoryPostgres(db)
	c.OrderRepo = infraRepo.NewOrderRepositoryPostgres(db)
	c.WebhookRepo = infraRepo.NewWebhookRepository(db)

	c.ProductUseCase = productUseCase.NewUseCase(c.ProductRepo)
	c.OrderUseCase = orderUseCase.NewUseCase(c.OrderRepo, c.ProductRepo)
	c.PaymentUseCase = paymentUseCase.NewPaymentUseCase(c.OrderRepo, c.WebhookRepo)

	c.ProductHandler = handler.NewProductHandler(c.ProductUseCase)
	c.OrderHandler = handler.NewOrderHandler(c.OrderUseCase)
	c.PaymentHandler = handler.NewPaymentHandler(c.PaymentUseCase, cfg.Webhook.Secret)

	return c
}
