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

func TestProduct_HasVariants(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		expected bool
	}{
		{
			name:     "no variants",
			product:  Product{Variants: []ProductVariant{}},
			expected: false,
		},
		{
			name:     "nil variants",
			product:  Product{Variants: nil},
			expected: false,
		},
		{
			name: "has variants",
			product: Product{
				Variants: []ProductVariant{
					{VariantName: "Size", VariantValue: "Large"},
				},
			},
			expected: true,
		},
		{
			name: "multiple variants",
			product: Product{
				Variants: []ProductVariant{
					{VariantName: "Size", VariantValue: "Small"},
					{VariantName: "Size", VariantValue: "Large"},
					{VariantName: "Color", VariantValue: "Red"},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.product.HasVariants()
			if result != tt.expected {
				t.Errorf("HasVariants() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProduct_GetTotalVariantStock(t *testing.T) {
	tests := []struct {
		name     string
		product  Product
		expected int
	}{
		{
			name:     "no variants",
			product:  Product{Variants: []ProductVariant{}},
			expected: 0,
		},
		{
			name: "single variant",
			product: Product{
				Variants: []ProductVariant{
					{VariantName: "Size", VariantValue: "Large", Quantity: 10},
				},
			},
			expected: 10,
		},
		{
			name: "multiple variants",
			product: Product{
				Variants: []ProductVariant{
					{VariantName: "Size", VariantValue: "Small", Quantity: 5},
					{VariantName: "Size", VariantValue: "Medium", Quantity: 8},
					{VariantName: "Size", VariantValue: "Large", Quantity: 12},
				},
			},
			expected: 25,
		},
		{
			name: "variants with zero quantity",
			product: Product{
				Variants: []ProductVariant{
					{VariantName: "Color", VariantValue: "Red", Quantity: 0},
					{VariantName: "Color", VariantValue: "Blue", Quantity: 15},
				},
			},
			expected: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.product.GetTotalVariantStock()
			if result != tt.expected {
				t.Errorf("GetTotalVariantStock() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProduct_GetVariantByNameValue(t *testing.T) {
	product := Product{
		Variants: []ProductVariant{
			{ID: uuid.New(), VariantName: "Size", VariantValue: "Small", Quantity: 5},
			{ID: uuid.New(), VariantName: "Size", VariantValue: "Large", Quantity: 10},
			{ID: uuid.New(), VariantName: "Color", VariantValue: "Red", Quantity: 8},
		},
	}

	tests := []struct {
		name         string
		variantName  string
		variantValue string
		expectFound  bool
		expectedQty  int
	}{
		{
			name:         "find existing variant",
			variantName:  "Size",
			variantValue: "Large",
			expectFound:  true,
			expectedQty:  10,
		},
		{
			name:         "find another variant",
			variantName:  "Color",
			variantValue: "Red",
			expectFound:  true,
			expectedQty:  8,
		},
		{
			name:         "variant name not found",
			variantName:  "Weight",
			variantValue: "Heavy",
			expectFound:  false,
		},
		{
			name:         "variant value not found",
			variantName:  "Size",
			variantValue: "XL",
			expectFound:  false,
		},
		{
			name:         "case sensitive - wrong case",
			variantName:  "size",
			variantValue: "large",
			expectFound:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := product.GetVariantByNameValue(tt.variantName, tt.variantValue)

			if tt.expectFound {
				if result == nil {
					t.Errorf("GetVariantByNameValue() = nil, expected variant to be found")
					return
				}
				if result.Quantity != tt.expectedQty {
					t.Errorf("GetVariantByNameValue() quantity = %v, want %v", result.Quantity, tt.expectedQty)
				}
			} else {
				if result != nil {
					t.Errorf("GetVariantByNameValue() = %v, expected nil", result)
				}
			}
		})
	}
}

func TestProduct_GetVariantByNameValue_EmptyProduct(t *testing.T) {
	product := Product{Variants: []ProductVariant{}}

	result := product.GetVariantByNameValue("Size", "Large")

	if result != nil {
		t.Errorf("GetVariantByNameValue() on empty variants = %v, want nil", result)
	}
}
