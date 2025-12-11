package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductVariant struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProductID      uuid.UUID `gorm:"type:uuid;not null;index"`
	VariantName    string    `gorm:"size:255;not null"`
	VariantValue   string    `gorm:"size:255;not null"`
	Price_Override *float64  `gorm:"type:decimal(10,2)"` // Pointer to distinguish between 0 and unset
	Quantity       int       `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	Product *Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

func (p *ProductVariant) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// GetPrice returns the effective price for this variant.
// If price_override is set, returns the override value.
// Otherwise, returns the base product price.
func (pv *ProductVariant) GetPrice() (float64, error) {
	if pv.Price_Override != nil {
		return *pv.Price_Override, nil
	}

	if pv.Product == nil {
		return 0, errors.New("Product not loaded: cannot determine variant price")
	}

	return pv.Product.Price, nil
}

// HasPriceOverride returns true if this variant has a custom price
func (pv *ProductVariant) HasPriceOverride() bool {
	return pv.Price_Override != nil
}

func (p *ProductVariant) ValidateForCreation() error {
	if p.VariantName == "" {
		return errors.New("Variant name is required")
	}
	if p.VariantValue == "" {
		return errors.New("Variant value is required")
	}
	if p.Price_Override != nil && *p.Price_Override < 0 {
		return errors.New("Variant price override cannot be negative")
	}
	if p.Quantity < 0 {
		return errors.New("Variant quantity cannot be negative")
	}
	if p.Quantity == 0 {
		return errors.New("Variant quantity must be greater than 0 for new variants")
	}
	return nil
}

// IsAvailable checks if the variant has enough stock
func (pv *ProductVariant) IsAvailable(quantity int) bool {
	return pv.Quantity >= quantity
}

// DecreaseStock reduces the variant's quantity
func (pv *ProductVariant) DecreaseStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("Quantity to decrease must be positive")
	}
	if !pv.IsAvailable(quantity) {
		return errors.New("Insufficient variant stock")
	}
	pv.Quantity -= quantity
	return nil
}

// IncreaseStock adds to the variant's quantity
func (pv *ProductVariant) IncreaseStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("Quantity to increase must be positive")
	}
	pv.Quantity += quantity
	return nil
}
