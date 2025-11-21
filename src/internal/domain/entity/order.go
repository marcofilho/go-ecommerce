package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/enum"
)

type Order struct {
	ID             uuid.UUID
	Customer_ID    int
	Products       []Product
	Total_Price    float64
	Status         enum.Status
	Payment_Status enum.PaymentStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewOrder(customerID int, products []Product, totalPrice float64) *Order {
	return &Order{
		ID:             uuid.New(),
		Customer_ID:    customerID,
		Products:       products,
		Total_Price:    totalPrice,
		Status:         enum.Pending,
		Payment_Status: enum.Unpaid,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
