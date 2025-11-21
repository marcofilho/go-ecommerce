package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey"`
	CustomerID    int           `gorm:"not null"`
	Products      []OrderItem   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	TotalPrice    float64       `gorm:"type:decimal(10,2);not null"`
	Status        OrderStatus   `gorm:"type:varchar(20);not null;default:'pending'"`
	PaymentStatus PaymentStatus `gorm:"type:varchar(20);not null;default:'unpaid'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (o *Order) Validate() error {
	if o.CustomerID <= 0 {
		return errors.New("customer ID is required")
	}
	if len(o.Products) == 0 {
		return errors.New("Order must have at least one product")
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
		total += item.Subtotal()
	}

	o.TotalPrice = total
}

func (o *Order) CanTransitionTo(newStatus OrderStatus) error {
	if o.Status == Pending {
		if newStatus == Completed || newStatus == Cancelled {
			return nil
		}
	}

	return errors.New("Invalid status transition")
}

func (o *Order) UpdateStatus(newStatus OrderStatus) error {
	if err := o.CanTransitionTo(newStatus); err != nil {
		return err
	}

	o.Status = newStatus
	o.UpdatedAt = time.Now()

	return nil
}
