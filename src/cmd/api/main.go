package main

import (
	"log"
	"net/http"

	_ "github.com/marcofilho/go-ecommerce/docs"
	"github.com/marcofilho/go-ecommerce/src/internal/config"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/database"
)

// @title Go E-Commerce API
// @version 1.0
// @description RESTful API for managing products and orders in an e-commerce system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email marco@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http

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

	// Initialize dependency container
	container := NewContainer(db)

	// Setup routes
	mux := SetupRoutes(container)

	// Start server
	serverAddr := ":" + cfg.Server.Port
	log.Printf("Server starting on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatal(err)
	}
}
