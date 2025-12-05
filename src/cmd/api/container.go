package main

import (
	"gorm.io/gorm"

	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/handler"
	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/middleware"
	"github.com/marcofilho/go-ecommerce/src/internal/config"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/auth"
	infraRepo "github.com/marcofilho/go-ecommerce/src/internal/infrastructure/repository"
	authUseCase "github.com/marcofilho/go-ecommerce/src/usecase/auth"
	orderUseCase "github.com/marcofilho/go-ecommerce/src/usecase/order"
	paymentUseCase "github.com/marcofilho/go-ecommerce/src/usecase/payment"
	productUseCase "github.com/marcofilho/go-ecommerce/src/usecase/product"
	productVariantUseCase "github.com/marcofilho/go-ecommerce/src/usecase/product_variant"
)

// Container holds all application dependencies
type Container struct {
	DB     *gorm.DB
	Config *config.Config

	// Repositories
	ProductRepo        repository.ProductRepository
	ProductVariantRepo repository.ProductVariantRepository
	OrderRepo          repository.OrderRepository
	WebhookRepo        repository.WebhookRepository
	UserRepo           repository.UserRepository

	// Infrastructure
	JWTProvider *auth.JWTProvider

	// Use Cases
	ProductUseCase        *productUseCase.UseCase
	ProductVariantUseCase *productVariantUseCase.UseCase
	OrderUseCase          *orderUseCase.UseCase
	PaymentUseCase        *paymentUseCase.PaymentUseCase
	AuthUseCase           *authUseCase.UseCase

	// Handlers
	ProductHandler        *handler.ProductHandler
	ProductVariantHandler *handler.ProductVariantHandler
	OrderHandler          *handler.OrderHandler
	PaymentHandler        *handler.PaymentHandler
	AuthHandler           *handler.AuthHandler

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware
}

// NewContainer creates and wires up all dependencies
func NewContainer(db *gorm.DB, cfg *config.Config) *Container {
	c := &Container{
		DB:     db,
		Config: cfg,
	}

	c.ProductRepo = infraRepo.NewProductRepositoryPostgres(db)
	c.ProductVariantRepo = infraRepo.NewProductVariantRepositoryPostgres(db)
	c.OrderRepo = infraRepo.NewOrderRepositoryPostgres(db)
	c.WebhookRepo = infraRepo.NewWebhookRepository(db)
	c.UserRepo = infraRepo.NewUserRepository(db)

	// Infrastructure
	c.JWTProvider = auth.NewJWTProvider(cfg.JWT.Secret, cfg.JWT.ExpirationHours)

	// Use Cases
	c.ProductUseCase = productUseCase.NewUseCase(c.ProductRepo)
	c.ProductVariantUseCase = productVariantUseCase.NewUseCase(c.ProductVariantRepo)
	c.OrderUseCase = orderUseCase.NewUseCase(c.OrderRepo, c.ProductRepo)
	c.PaymentUseCase = paymentUseCase.NewPaymentUseCase(c.OrderRepo, c.WebhookRepo)
	c.AuthUseCase = authUseCase.NewUseCase(c.UserRepo, c.JWTProvider)

	// Handlers
	c.ProductHandler = handler.NewProductHandler(c.ProductUseCase)
	c.ProductVariantHandler = handler.NewProductVariantHandler(c.ProductVariantUseCase)
	c.OrderHandler = handler.NewOrderHandler(c.OrderUseCase)
	c.PaymentHandler = handler.NewPaymentHandler(c.PaymentUseCase, cfg.Webhook.Secret)
	c.AuthHandler = handler.NewAuthHandler(c.AuthUseCase)

	// Middleware
	c.AuthMiddleware = middleware.NewAuthMiddleware(c.AuthUseCase)

	return c
}
