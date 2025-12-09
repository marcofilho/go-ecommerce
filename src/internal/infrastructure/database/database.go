package database

import (
	"fmt"

	"github.com/marcofilho/go-ecommerce/src/internal/config"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	// AutoMigrate creates tables and indexes
	// Order matters: tables with foreign keys must come after their references
	return db.AutoMigrate(
		&entity.User{},            // No dependencies
		&entity.Category{},        // No dependencies
		&entity.Product{},         // No dependencies
		&entity.ProductVariant{},  // Foreign key to Product
		&entity.ProductCategory{}, // Foreign key to Product and Category (junction table)
		&entity.Order{},           // Foreign key to User (CustomerID)
		&entity.OrderItem{},       // Foreign key to Order and Product
		&entity.WebhookLog{},      // Foreign key to Order
	)
}
