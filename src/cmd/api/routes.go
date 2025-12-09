package main

import (
	"net/http"

	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/middleware"
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

	// Product routes
	// Public: Anyone can view products
	mux.HandleFunc("GET /api/products", c.ProductHandler.ListProducts)
	mux.HandleFunc("GET /api/products/{id}", c.ProductHandler.GetProduct)

	// Admin only: Create, update, delete products
	mux.Handle("POST /api/products", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionCreateProduct)(
			http.HandlerFunc(c.ProductHandler.CreateProduct),
		),
	))
	mux.Handle("PUT /api/products/{id}", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionUpdateProduct)(
			http.HandlerFunc(c.ProductHandler.UpdateProduct),
		),
	))
	mux.Handle("DELETE /api/products/{id}", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionDeleteProduct)(
			http.HandlerFunc(c.ProductHandler.DeleteProduct),
		),
	))

	// Product Variant routes
	// Public: View product variants for a product
	mux.HandleFunc("GET /api/products/{id}/variants", c.ProductVariantHandler.ListProductVariants)

	// Admin only: Create product variant for a product
	mux.Handle("POST /api/products/{id}/variants", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionCreateProduct)(
			http.HandlerFunc(c.ProductVariantHandler.CreateProductVariant),
		),
	))

	// Admin only: Update and delete product variants
	mux.Handle("PUT /api/variants/{variant_id}", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionUpdateProduct)(
			http.HandlerFunc(c.ProductVariantHandler.UpdateProductVariant),
		),
	))
	mux.Handle("DELETE /api/variants/{variant_id}", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionDeleteProduct)(
			http.HandlerFunc(c.ProductVariantHandler.DeleteProductVariant),
		),
	))

	// Category routes
	// Public: List categories
	mux.HandleFunc("GET /api/categories", c.CategoryHandler.ListCategories)

	// Admin only: Create categories
	mux.Handle("POST /api/categories", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionCreateProduct)(
			http.HandlerFunc(c.CategoryHandler.CreateCategory),
		),
	))

	// Product-Category relationship routes
	// Public: Get product categories
	mux.HandleFunc("GET /api/products/{id}/categories", c.CategoryHandler.GetProductCategories)

	// Admin only: Assign category to product
	mux.Handle("POST /api/products/{id}/categories", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionCreateProduct)(
			http.HandlerFunc(c.CategoryHandler.AssignCategoryToProduct),
		),
	))

	// Admin only: Remove category from product
	mux.Handle("DELETE /api/products/{id}/categories/{category_id}", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionDeleteProduct)(
			http.HandlerFunc(c.CategoryHandler.RemoveCategoryFromProduct),
		),
	))

	// Order routes
	// Authenticated users: Create and view orders
	mux.Handle("POST /api/orders", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionCreateOrder)(
			http.HandlerFunc(c.OrderHandler.CreateOrder),
		),
	))
	mux.Handle("GET /api/orders", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionListOrders)(
			http.HandlerFunc(c.OrderHandler.ListOrders),
		),
	))
	mux.Handle("GET /api/orders/{id}", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionViewOrder)(
			http.HandlerFunc(c.OrderHandler.GetOrder),
		),
	))

	// Admin only: Update order status
	mux.Handle("PUT /api/orders/{id}/status", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionUpdateOrderStatus)(
			http.HandlerFunc(c.OrderHandler.UpdateOrderStatus),
		),
	))

	// Payment webhook routes
	mux.HandleFunc("POST /api/payment-webhook", c.PaymentHandler.PaymentWebhookHandler) // Public - external integration

	// Admin only: View webhook history
	mux.Handle("GET /api/orders/{id}/payment-history", c.AuthMiddleware.Authenticate(
		c.AuthMiddleware.RequirePermission(middleware.PermissionViewWebhookHistory)(
			http.HandlerFunc(c.PaymentHandler.GetWebhookHistoryHandler),
		),
	))

	return mux
}
