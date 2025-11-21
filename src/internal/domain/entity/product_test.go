package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name    string
		product Product
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid product",
			product: Product{
				Name:     "Laptop",
				Price:    999.99,
				Quantity: 10,
			},
			wantErr: false,
		},
		{
			name: "valid product with zero quantity",
			product: Product{
				Name:     "Laptop",
				Price:    999.99,
				Quantity: 0,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			product: Product{
				Name:     "",
				Price:    999.99,
				Quantity: 10,
			},
			wantErr: true,
			errMsg:  "Product name is required",
		},
		{
			name: "negative price",
			product: Product{
				Name:     "Laptop",
				Price:    -10.00,
				Quantity: 10,
			},
			wantErr: true,
			errMsg:  "Product price cannot be negative",
		},
		{
			name: "negative quantity",
			product: Product{
				Name:     "Laptop",
				Price:    999.99,
				Quantity: -5,
			},
			wantErr: true,
			errMsg:  "Product quantity cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestProduct_ValidateForCreation(t *testing.T) {
	tests := []struct {
		name    string
		product Product
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid product with stock",
			product: Product{
				Name:     "Laptop",
				Price:    999.99,
				Quantity: 10,
			},
			wantErr: false,
		},
		{
			name: "zero quantity for new product",
			product: Product{
				Name:     "Laptop",
				Price:    999.99,
				Quantity: 0,
			},
			wantErr: true,
			errMsg:  "Product quantity must be greater than 0 for new products",
		},
		{
			name: "invalid product - empty name",
			product: Product{
				Name:     "",
				Price:    999.99,
				Quantity: 10,
			},
			wantErr: true,
			errMsg:  "Product name is required",
		},
		{
			name: "invalid product - negative price",
			product: Product{
				Name:     "Laptop",
				Price:    -10.00,
				Quantity: 10,
			},
			wantErr: true,
			errMsg:  "Product price cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.ValidateForCreation()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateForCreation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("ValidateForCreation() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestProduct_IsAvailable(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		quantity int
		want     bool
	}{
		{
			name:     "enough stock",
			product:  Product{Quantity: 10},
			quantity: 5,
			want:     true,
		},
		{
			name:     "exact stock",
			product:  Product{Quantity: 10},
			quantity: 10,
			want:     true,
		},
		{
			name:     "insufficient stock",
			product:  Product{Quantity: 5},
			quantity: 10,
			want:     false,
		},
		{
			name:     "out of stock",
			product:  Product{Quantity: 0},
			quantity: 1,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.product.IsAvailable(tt.quantity); got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProduct_DecreaseStock(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		quantity int
		wantErr  bool
		wantQty  int
	}{
		{
			name:     "successful decrease",
			product:  Product{Quantity: 10},
			quantity: 3,
			wantErr:  false,
			wantQty:  7,
		},
		{
			name:     "decrease to zero",
			product:  Product{Quantity: 5},
			quantity: 5,
			wantErr:  false,
			wantQty:  0,
		},
		{
			name:     "insufficient stock",
			product:  Product{Quantity: 5},
			quantity: 10,
			wantErr:  true,
			wantQty:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.DecreaseStock(tt.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecreaseStock() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.product.Quantity != tt.wantQty {
				t.Errorf("DecreaseStock() quantity = %v, want %v", tt.product.Quantity, tt.wantQty)
			}
		})
	}
}

func TestProduct_IncreaseStock(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		quantity int
		wantErr  bool
		wantQty  int
	}{
		{
			name:     "successful increase",
			product:  Product{Quantity: 10},
			quantity: 5,
			wantErr:  false,
			wantQty:  15,
		},
		{
			name:     "increase from zero",
			product:  Product{Quantity: 0},
			quantity: 10,
			wantErr:  false,
			wantQty:  10,
		},
		{
			name:     "negative quantity",
			product:  Product{Quantity: 10},
			quantity: -5,
			wantErr:  true,
			wantQty:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.IncreaseStock(tt.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncreaseStock() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.product.Quantity != tt.wantQty {
				t.Errorf("IncreaseStock() quantity = %v, want %v", tt.product.Quantity, tt.wantQty)
			}
		})
	}
}

func TestProduct_BeforeCreate(t *testing.T) {
	t.Run("generates UUID if not set", func(t *testing.T) {
		product := &Product{}
		err := product.BeforeCreate(nil)
		if err != nil {
			t.Errorf("BeforeCreate() error = %v", err)
		}
		if product.ID == uuid.Nil {
			t.Error("BeforeCreate() did not generate UUID")
		}
	})

	t.Run("keeps existing UUID", func(t *testing.T) {
		existingID := uuid.New()
		product := &Product{ID: existingID}
		err := product.BeforeCreate(nil)
		if err != nil {
			t.Errorf("BeforeCreate() error = %v", err)
		}
		if product.ID != existingID {
			t.Error("BeforeCreate() changed existing UUID")
		}
	})
}
