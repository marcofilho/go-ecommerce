package entity

import (
	"errors"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	OrderID    uuid.UUID  `gorm:"type:uuid;not null"`
	ProductID  uuid.UUID  `gorm:"type:uuid;not null"`
	VariantID  *uuid.UUID `gorm:"type:uuid"`
	Quantity   int        `gorm:"not null"`
	Price      float64    `gorm:"type:decimal(10,2);not null"`
	TotalPrice float64    `gorm:"type:decimal(10,2);not null"`
}

func (oi *OrderItem) Validate() error {
	if oi.ID == uuid.Nil {
		return errors.New("Order item ID is required")
	}
	if oi.ProductID == uuid.Nil {
		return errors.New("Product ID is required")
	}
	if oi.Quantity <= 0 {
		return errors.New("Quantity must be greater than 0")
	}
	if oi.Price < 0 {
		return errors.New("Price cannot be negative")
	}
	if oi.TotalPrice < 0 {
		return errors.New("Total price cannot be negative")
	}
	return nil
}

func (oi *OrderItem) CalculateTotal() {
	oi.TotalPrice = oi.Price * float64(oi.Quantity)
}

func (oi *OrderItem) Subtotal() float64 {
	return oi.TotalPrice
}
