package entity

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestDiscount_Validate(t *testing.T) {
	tests := []struct {
		name     string
		discount Discount
		wantErr  bool
	}{
		{
			name: "Valid discount",
			discount: Discount{
				PromoCode:    "SAVE20",
				DiscountType: Percentage,
				Value:        20.0,
				Active:       true,
			},
			wantErr: false,
		},
		{
			name: "Empty promo code",
			discount: Discount{
				PromoCode:    "",
				DiscountType: Percentage,
				Value:        20.0,
				Active:       true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.discount.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscount_ApplyDiscount_Percentage(t *testing.T) {
	tests := []struct {
		name          string
		discount      Discount
		total         float64
		expectedTotal float64
		wantErr       bool
	}{
		{
			name: "20% off $100",
			discount: Discount{
				DiscountType: Percentage,
				Value:        20.0,
			},
			total:         100.0,
			expectedTotal: 80.0,
			wantErr:       false,
		},
		{
			name: "50% off $200",
			discount: Discount{
				DiscountType: Percentage,
				Value:        50.0,
			},
			total:         200.0,
			expectedTotal: 100.0,
			wantErr:       false,
		},
		{
			name: "10% off $150",
			discount: Discount{
				DiscountType: Percentage,
				Value:        10.0,
			},
			total:         150.0,
			expectedTotal: 135.0,
			wantErr:       false,
		},
		{
			name: "100% off (free)",
			discount: Discount{
				DiscountType: Percentage,
				Value:        100.0,
			},
			total:         100.0,
			expectedTotal: 0.0,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.discount.ApplyDiscount(tt.total)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyDiscount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expectedTotal {
				t.Errorf("ApplyDiscount() = %v, want %v", result, tt.expectedTotal)
			}
		})
	}
}

func TestDiscount_ApplyDiscount_Amount(t *testing.T) {
	tests := []struct {
		name          string
		discount      Discount
		total         float64
		expectedTotal float64
		wantErr       bool
	}{
		{
			name: "$15 off $100",
			discount: Discount{
				DiscountType: Amount,
				Value:        15.0,
			},
			total:         100.0,
			expectedTotal: 85.0,
			wantErr:       false,
		},
		{
			name: "$50 off $200",
			discount: Discount{
				DiscountType: Amount,
				Value:        50.0,
			},
			total:         200.0,
			expectedTotal: 150.0,
			wantErr:       false,
		},
		{
			name: "$10 off $10 (free)",
			discount: Discount{
				DiscountType: Amount,
				Value:        10.0,
			},
			total:         10.0,
			expectedTotal: 0.0,
			wantErr:       false,
		},
		{
			name: "$100 off $50 (negative)",
			discount: Discount{
				DiscountType: Amount,
				Value:        100.0,
			},
			total:         50.0,
			expectedTotal: -50.0,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.discount.ApplyDiscount(tt.total)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyDiscount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expectedTotal {
				t.Errorf("ApplyDiscount() = %v, want %v", result, tt.expectedTotal)
			}
		})
	}
}

func TestDiscount_ApplyDiscount_InvalidType(t *testing.T) {
	discount := Discount{
		DiscountType: "invalid",
		Value:        20.0,
	}

	result, err := discount.ApplyDiscount(100.0)
	if err == nil {
		t.Error("ApplyDiscount() expected error for invalid discount type, got nil")
	}
	if result != 100.0 {
		t.Errorf("ApplyDiscount() with invalid type should return original total, got %v", result)
	}
}

func TestDiscount_BeforeCreate(t *testing.T) {
	t.Run("Should generate UUID if not set", func(t *testing.T) {
		discount := &Discount{
			PromoCode: "TEST",
		}

		err := discount.BeforeCreate(&gorm.DB{})
		if err != nil {
			t.Errorf("BeforeCreate() unexpected error: %v", err)
		}

		if discount.ID == uuid.Nil {
			t.Error("BeforeCreate() should generate UUID")
		}
	})

	t.Run("Should not override existing UUID", func(t *testing.T) {
		existingID := uuid.New()
		discount := &Discount{
			ID:        existingID,
			PromoCode: "TEST",
		}

		err := discount.BeforeCreate(&gorm.DB{})
		if err != nil {
			t.Errorf("BeforeCreate() unexpected error: %v", err)
		}

		if discount.ID != existingID {
			t.Error("BeforeCreate() should not override existing UUID")
		}
	})
}
