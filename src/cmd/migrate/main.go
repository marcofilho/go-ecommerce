package main

import (
	"log"

	"github.com/marcofilho/go-ecommerce/src/internal/config"
	"github.com/marcofilho/go-ecommerce/src/internal/infrastructure/database"
)

func main() {
	log.Println("Running database migrations...")

	cfg := config.Load()

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Migrations completed successfully!")
}
