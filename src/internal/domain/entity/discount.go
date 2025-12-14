package entity

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountType string

const (
	Percentage DiscountType = "percentage"
	Amount     DiscountType = "amount"
)

type Discount struct {
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey"`
	PromoCode    string       `gorm:"uniqueIndex;not null"`
	DiscountType DiscountType `gorm:"type:varchar(20);not null;default:'amount'"`
	Value        float64      `gorm:"type:decimal(10,2);not null"`
	Active       bool         `gorm:"not null;default:true"`
}

func (d *Discount) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

func (d *Discount) Validate() error {
	if d.PromoCode == "" {
		return errors.New("Promo code is required")
	}
	return nil
}

func (d *Discount) IsActive() bool {
	return d.Active
}

func (d *Discount) Activate() {
	d.Active = true
}

func (d *Discount) Deactivate() {
	d.Active = false
}

func (d *Discount) ApplyDiscount(total float64) (float64, error) {
	if !d.Active {
		return 0, errors.New("discount is not active")
	}

	switch d.DiscountType {
	case Percentage:
		discountAmount := total * (d.Value / 100)
		return total - discountAmount, nil
	case Amount:
		return total - d.Value, nil
	default:
		return total, errors.New("invalid discount type")
	}
}
