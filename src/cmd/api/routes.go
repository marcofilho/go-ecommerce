package main

import (
	"net/http"

	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRoutes configures all application routes
func SetupRoutes(c *Container) *http.ServeMux {
	mux := http.NewServeMux()

	// Swagger documentation
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Auth routes (public)
	mux.HandleFunc("POST /api/auth/register", c.AuthHandler.Register)
	mux.HandleFunc("POST /api/auth/login", c.AuthHandler.Login)

	// Product routes (protected - require authentication)
	mux.Handle("POST /api/products", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequireRole(entity.RoleAdmin)(
			http.HandlerFunc(c.ProductHandler.CreateProduct),
		),
	))
	mux.HandleFunc("GET /api/products", c.ProductHandler.ListProducts)    // Public
	mux.HandleFunc("GET /api/products/{id}", c.ProductHandler.GetProduct) // Public
	mux.Handle("PUT /api/products/{id}", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequireRole(entity.RoleAdmin)(
			http.HandlerFunc(c.ProductHandler.UpdateProduct),
		),
	))
	mux.Handle("DELETE /api/products/{id}", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequireRole(entity.RoleAdmin)(
			http.HandlerFunc(c.ProductHandler.DeleteProduct),
		),
	))

	// Order routes (protected - authenticated users only)
	mux.Handle("POST /api/orders", c.AuthMiddleware.Authenticate(
		http.HandlerFunc(c.OrderHandler.CreateOrder),
	))
	mux.Handle("GET /api/orders", c.AuthMiddleware.Authenticate(
		http.HandlerFunc(c.OrderHandler.ListOrders),
	))
	mux.Handle("GET /api/orders/{id}", c.AuthMiddleware.Authenticate(
		http.HandlerFunc(c.OrderHandler.GetOrder),
	))
	mux.Handle("PUT /api/orders/{id}/status", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequireRole(entity.RoleAdmin)(
			http.HandlerFunc(c.OrderHandler.UpdateOrderStatus),
		),
	))

	// Payment webhook routes (public - external integration)
	mux.HandleFunc("POST /api/payment-webhook", c.PaymentHandler.PaymentWebhookHandler)
	mux.Handle("GET /api/orders/{id}/payment-history", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequireRole(entity.RoleAdmin)(
			http.HandlerFunc(c.PaymentHandler.GetWebhookHistoryHandler),
		),
	))

	return mux
}
