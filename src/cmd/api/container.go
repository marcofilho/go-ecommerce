package main

import (
	"gorm.io/gorm"

	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/handler"
	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/middleware"
	"github.com/marcofilho/go-ecommerce/src/internal/config"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/audit"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/auth"
	infraRepo "github.com/marcofilho/go-ecommerce/src/internal/infrastructure/repository"
	authUseCase "github.com/marcofilho/go-ecommerce/src/usecase/auth"
	categoryUseCase "github.com/marcofilho/go-ecommerce/src/usecase/category"
	orderUseCase "github.com/marcofilho/go-ecommerce/src/usecase/order"
	paymentUseCase "github.com/marcofilho/go-ecommerce/src/usecase/payment"
	productUseCase "github.com/marcofilho/go-ecommerce/src/usecase/product"
	productVariantUseCase "github.com/marcofilho/go-ecommerce/src/usecase/product_variant"
)

// Services holds common infrastructure services
type Services struct {
	audit audit.AuditService
}

func (s *Services) GetAuditService() audit.AuditService {
	return s.audit
}

// Container holds all application dependencies
type Container struct {
	DB     *gorm.DB
	Config *config.Config

	// Repositories
	ProductRepo        repository.ProductRepository
	ProductVariantRepo repository.ProductVariantRepository
	CategoryRepo       repository.CategoryRepository
	OrderRepo          repository.OrderRepository
	WebhookRepo        repository.WebhookRepository
	UserRepo           repository.UserRepository
	AuditLogRepo       repository.AuditLogRepository

	// Infrastructure
	JWTProvider *auth.JWTProvider
	Services    *Services

	// Use Cases
	ProductUseCase        *productUseCase.UseCase
	ProductVariantUseCase *productVariantUseCase.UseCase
	CategoryUseCase       *categoryUseCase.UseCase
	OrderUseCase          *orderUseCase.UseCase
	PaymentUseCase        *paymentUseCase.PaymentUseCase
	AuthUseCase           *authUseCase.UseCase

	// Handlers
	ProductHandler        *handler.ProductHandler
	ProductVariantHandler *handler.ProductVariantHandler
	CategoryHandler       *handler.CategoryHandler
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
	c.CategoryRepo = infraRepo.NewCategoryRepository(db)
	c.OrderRepo = infraRepo.NewOrderRepositoryPostgres(db)
	c.WebhookRepo = infraRepo.NewWebhookRepository(db)
	c.UserRepo = infraRepo.NewUserRepository(db)
	c.AuditLogRepo = infraRepo.NewAuditLogRepository(db)

	// Infrastructure Services
	c.JWTProvider = auth.NewJWTProvider(cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	c.Services = &Services{
		audit: audit.NewAuditService(c.AuditLogRepo),
	}

	// Use Cases
	c.ProductUseCase = productUseCase.NewUseCase(c.ProductRepo, c.Services)
	c.ProductVariantUseCase = productVariantUseCase.NewUseCase(c.ProductVariantRepo)
	c.CategoryUseCase = categoryUseCase.NewUseCase(c.CategoryRepo)
	c.OrderUseCase = orderUseCase.NewUseCase(c.OrderRepo, c.ProductRepo, c.ProductVariantRepo, c.Services)
	c.PaymentUseCase = paymentUseCase.NewPaymentUseCase(c.OrderRepo, c.WebhookRepo, c.Services)
	c.AuthUseCase = authUseCase.NewUseCase(c.UserRepo, c.JWTProvider)

	// Handlers
	c.ProductHandler = handler.NewProductHandler(c.ProductUseCase)
	c.ProductVariantHandler = handler.NewProductVariantHandler(c.ProductVariantUseCase)
	c.CategoryHandler = handler.NewCategoryHandler(c.CategoryUseCase)
	c.OrderHandler = handler.NewOrderHandler(c.OrderUseCase)
	c.PaymentHandler = handler.NewPaymentHandler(c.PaymentUseCase, cfg.Webhook.Secret)
	c.AuthHandler = handler.NewAuthHandler(c.AuthUseCase)

	// Middleware
	c.AuthMiddleware = middleware.NewAuthMiddleware(c.AuthUseCase)

	return c
}
