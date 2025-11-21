package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	Quantity    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
