package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Cancelled OrderStatus = "cancelled"
	Completed OrderStatus = "completed"
)

type PaymentStatus string

const (
	Unpaid PaymentStatus = "unpaid"
	Paid   PaymentStatus = "paid"
	Failed PaymentStatus = "failed"
)

type Order struct {
	ID             uuid.UUID
	Customer_ID    int
	Products       []Product
	Total_Price    float64
	Status         OrderStatus
	Payment_Status PaymentStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (o *Order) Validate() error {
	if o.Customer_ID <= 0 {
		return errors.New("customer ID is required")
	}
	if len(o.Products) == 0 {
		return errors.New("order must have at least one item")
	}
	for _, product := range o.Products {
		if err := product.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (o *Order) CalculateTotal() {
	total := 0.0
	for _, item := range o.Products {
		total += item.Price * float64(item.Quantity)
	}
	o.Total_Price = total
}
