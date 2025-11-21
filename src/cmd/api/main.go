package main

import (
	"log"
	"net/http"

	handler "github.com/marcofilho/go-ecommerce/src/internal/adapter/http"
	"github.com/marcofilho/go-ecommerce/src/internal/config"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/database"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/repository"
	orderUseCase "github.com/marcofilho/go-ecommerce/src/usecase/order"
	productUseCase "github.com/marcofilho/go-ecommerce/src/usecase/product"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	productRepo := repository.NewProductRepositoryPostgres(db)
	orderRepo := repository.NewOrderRepositoryPostgres(db)

	// Initialize use cases
	productUC := productUseCase.NewUseCase(productRepo)
	orderUC := orderUseCase.NewUseCase(orderRepo, productRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productUC)
	orderHandler := handler.NewOrderHandler(orderUC)

	// Setup routes
	mux := http.NewServeMux()

	// Product routes
	mux.HandleFunc("POST /api/products", productHandler.CreateProduct)
	mux.HandleFunc("GET /api/products", productHandler.ListProducts)
	mux.HandleFunc("GET /api/products/{id}", productHandler.GetProduct)
	mux.HandleFunc("PUT /api/products/{id}", productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /api/products/{id}", productHandler.DeleteProduct)

	// Order routes
	mux.HandleFunc("POST /api/orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /api/orders", orderHandler.ListOrders)
	mux.HandleFunc("GET /api/orders/{id}", orderHandler.GetOrder)
	mux.HandleFunc("PUT /api/orders/{id}/status", orderHandler.UpdateOrderStatus)

	// Start server
	serverAddr := ":" + cfg.Server.Port
	log.Printf("Server starting on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatal(err)
	}
}
