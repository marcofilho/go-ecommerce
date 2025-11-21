package entity

import (
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

func NewProduct(name, description string, price float64, quantity int) *Product {
	return &Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
