package main

import (
	"log"
	"net/http"

	"github.com/marcofilho/go-ecommerce/src/internal/config"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/database"
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
