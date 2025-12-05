package entity

import (
	"errors"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID        uint       `gorm:"primaryKey"`
	OrderID   uuid.UUID  `gorm:"type:uuid;not null"`
	ProductID uuid.UUID  `gorm:"type:uuid;not null"`
	VariantID *uuid.UUID `gorm:"type:uuid"` // Optional: if ordering a specific variant
	Quantity  int        `gorm:"not null"`
	Price     float64    `gorm:"type:decimal(10,2);not null"`
}

func (oi *OrderItem) Validate() error {
	if oi.ProductID == uuid.Nil {
		return errors.New("Product ID is required")
	}
	if oi.Quantity <= 0 {
		return errors.New("Quantity must be greater than 0")
	}
	if oi.Price < 0 {
		return errors.New("Price cannot be negative")
	}
	return nil
}

func (oi *OrderItem) Subtotal() float64 {
	return oi.Price * float64(oi.Quantity)
}
