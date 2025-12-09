package entity

import (
	"github.com/google/uuid"
)

// ProductCategory represents a many-to-many relationship between products and categories
type ProductCategory struct {
	ProductID  uuid.UUID `gorm:"type:uuid;primaryKey;index:idx_product_category"`
	CategoryID uuid.UUID `gorm:"type:uuid;primaryKey;index:idx_product_category"`

	// Foreign key relationships
	Product  Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Category Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
}
