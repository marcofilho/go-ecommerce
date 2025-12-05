package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"type:decimal(10,2);not null"`
	Quantity    int       `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relations (not stored in DB, loaded via GORM preload)
	Variants []ProductVariant `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return errors.New("Product name is required")
	}
	if p.Price < 0 {
		return errors.New("Product price cannot be negative")
	}
	if p.Quantity < 0 {
		return errors.New("Product quantity cannot be negative")
	}

	return nil
}

func (p *Product) ValidateForCreation() error {
	if err := p.Validate(); err != nil {
		return err
	}
	if p.Quantity == 0 {
		return errors.New("Product quantity must be greater than 0 for new products")
	}
	return nil
}

func (p *Product) IsAvailable(quantity int) bool {
	return p.Quantity >= quantity
}

func (p *Product) DecreaseStock(quantity int) error {
	if !p.IsAvailable(quantity) {
		return errors.New("Insufficient stock")
	}

	p.Quantity -= quantity
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Product) IncreaseStock(quantity int) error {
	if quantity < 0 {
		return errors.New("Quantity must be positive")
	}

	p.Quantity += quantity
	p.UpdatedAt = time.Now()

	return nil
}

// HasVariants returns true if the product has any variants
func (p *Product) HasVariants() bool {
	return len(p.Variants) > 0
}

// GetTotalVariantStock returns the total stock across all variants
func (p *Product) GetTotalVariantStock() int {
	if !p.HasVariants() {
		return 0
	}

	total := 0
	for _, variant := range p.Variants {
		total += variant.Quantity
	}
	return total
}

// GetVariantByNameValue finds a variant by name and value
func (p *Product) GetVariantByNameValue(name, value string) *ProductVariant {
	for i := range p.Variants {
		if p.Variants[i].VariantName == name && p.Variants[i].VariantValue == value {
			return &p.Variants[i]
		}
	}
	return nil
}
