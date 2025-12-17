package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountType string

const (
	Percentage DiscountType = "percentage"
	Amount     DiscountType = "amount"
)

type Discount struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	PromoCode         string         `gorm:"uniqueIndex;not null" json:"promo_code"`
	DiscountType      DiscountType   `gorm:"type:varchar(20);not null;default:'amount'" json:"discount_type"`
	Value             float64        `gorm:"type:decimal(10,2);not null" json:"value"`
	Active            bool           `gorm:"not null;default:true" json:"active"`
	MinPurchaseAmount *float64       `gorm:"type:decimal(10,2)" json:"min_purchase_amount,omitempty"`
	MaxDiscountAmount *float64       `gorm:"type:decimal(10,2)" json:"max_discount_amount,omitempty"`
	UsageLimit        *int           `json:"usage_limit,omitempty"`
	UsageCount        int            `gorm:"default:0" json:"usage_count"`
	ValidFrom         *time.Time     `json:"valid_from,omitempty"`
	ValidUntil        *time.Time     `json:"valid_until,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Products   []Product  `gorm:"many2many:discount_products;" json:"products,omitempty"`
	Categories []Category `gorm:"many2many:discount_categories;" json:"categories,omitempty"`
	Users      []User     `gorm:"many2many:discount_users;" json:"users,omitempty"`
}

// DiscountUser represents the junction table for discount-user relationships
type DiscountUser struct {
	DiscountID uuid.UUID `gorm:"primaryKey;type:uuid" json:"discount_id"`
	UserID     uuid.UUID `gorm:"primaryKey;type:uuid" json:"user_id"`
	UsageCount int       `gorm:"default:0" json:"usage_count"`
	UsageLimit *int      `json:"usage_limit,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	Discount   *Discount `gorm:"foreignKey:DiscountID" json:"-"`
	User       *User     `gorm:"foreignKey:UserID" json:"-"`
}

func (DiscountUser) TableName() string {
	return "discount_users"
}

func (d *Discount) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

func (d *Discount) Validate() error {
	if d.PromoCode == "" {
		return errors.New("promo code is required")
	}
	if d.Value <= 0 {
		return errors.New("discount value must be positive")
	}
	if d.DiscountType == Percentage && d.Value > 100 {
		return errors.New("percentage discount cannot exceed 100%")
	}
	if d.ValidFrom != nil && d.ValidUntil != nil && d.ValidFrom.After(*d.ValidUntil) {
		return errors.New("valid_from must be before valid_until")
	}
	return nil
}

// IsValid checks if the discount is currently valid
func (d *Discount) IsValid() bool {
	if !d.Active {
		return false
	}

	now := time.Now()

	if d.ValidFrom != nil && now.Before(*d.ValidFrom) {
		return false
	}

	if d.ValidUntil != nil && now.After(*d.ValidUntil) {
		return false
	}

	if d.UsageLimit != nil && d.UsageCount >= *d.UsageLimit {
		return false
	}

	return true
}

// CanBeUsedBy checks if a specific user can use this discount
func (d *Discount) CanBeUsedBy(userUsageCount int, userUsageLimit *int) bool {
	if !d.IsValid() {
		return false
	}

	if userUsageLimit != nil && userUsageCount >= *userUsageLimit {
		return false
	}

	return true
}

// AppliesTo checks if discount applies to a specific product based on product ID and category IDs
func (d *Discount) AppliesTo(productID uuid.UUID, categoryIDs []uuid.UUID) bool {
	// If no specific products or categories, applies to all
	if len(d.Products) == 0 && len(d.Categories) == 0 {
		return true
	}

	// Check if product is in discount's product list
	for _, p := range d.Products {
		if p.ID == productID {
			return true
		}
	}

	// Check if product's category is in discount's category list
	for _, dc := range d.Categories {
		for _, catID := range categoryIDs {
			if dc.ID == catID {
				return true
			}
		}
	}

	return false
}

// ApplyDiscount calculates the discounted amount
func (d *Discount) ApplyDiscount(total float64) (float64, error) {
	if d.MinPurchaseAmount != nil && total < *d.MinPurchaseAmount {
		return total, errors.New("minimum purchase amount not met")
	}

	var discountAmount float64

	switch d.DiscountType {
	case Percentage:
		discountAmount = total * (d.Value / 100)
		if d.MaxDiscountAmount != nil && discountAmount > *d.MaxDiscountAmount {
			discountAmount = *d.MaxDiscountAmount
		}
	case Amount:
		discountAmount = d.Value
	default:
		return total, errors.New("invalid discount type")
	}

	newTotal := total - discountAmount
	if newTotal < 0 {
		newTotal = 0
	}

	return newTotal, nil
}
