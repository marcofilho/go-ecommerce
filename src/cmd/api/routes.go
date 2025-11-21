package main

import "net/http"

// SetupRoutes configures all application routes
func SetupRoutes(c *Container) *http.ServeMux {
	mux := http.NewServeMux()

	// Product routes
	mux.HandleFunc("POST /api/products", c.ProductHandler.CreateProduct)
	mux.HandleFunc("GET /api/products", c.ProductHandler.ListProducts)
	mux.HandleFunc("GET /api/products/{id}", c.ProductHandler.GetProduct)
	mux.HandleFunc("PUT /api/products/{id}", c.ProductHandler.UpdateProduct)
	mux.HandleFunc("DELETE /api/products/{id}", c.ProductHandler.DeleteProduct)

	// Order routes
	mux.HandleFunc("POST /api/orders", c.OrderHandler.CreateOrder)
	mux.HandleFunc("GET /api/orders", c.OrderHandler.ListOrders)
	mux.HandleFunc("GET /api/orders/{id}", c.OrderHandler.GetOrder)
	mux.HandleFunc("PUT /api/orders/{id}/status", c.OrderHandler.UpdateOrderStatus)

	return mux
}
