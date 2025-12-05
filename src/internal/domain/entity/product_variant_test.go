package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestProductVariant_GetPrice_WithOverride(t *testing.T) {
	overridePrice := 149.99
	variant := &ProductVariant{
		ID:             uuid.New(),
		ProductID:      uuid.New(),
		VariantName:    "Size",
		VariantValue:   "Large",
		Price_Override: &overridePrice,
		Quantity:       10,
	}

	price, err := variant.GetPrice()

	if err != nil {
		t.Fatalf("GetPrice() error = %v, want nil", err)
	}

	if price != overridePrice {
		t.Errorf("GetPrice() = %v, want %v", price, overridePrice)
	}
}

func TestProductVariant_GetPrice_WithoutOverride(t *testing.T) {
	basePrice := 99.99
	product := &Product{
		ID:    uuid.New(),
		Name:  "T-Shirt",
		Price: basePrice,
	}

	variant := &ProductVariant{
		ID:             uuid.New(),
		ProductID:      product.ID,
		VariantName:    "Color",
		VariantValue:   "Blue",
		Price_Override: nil,
		Quantity:       10,
		Product:        product,
	}

	price, err := variant.GetPrice()

	if err != nil {
		t.Fatalf("GetPrice() error = %v, want nil", err)
	}

	if price != basePrice {
		t.Errorf("GetPrice() = %v, want %v (base product price)", price, basePrice)
	}
}

func TestProductVariant_GetPrice_WithoutOverride_NoProduct(t *testing.T) {
	variant := &ProductVariant{
		ID:             uuid.New(),
		ProductID:      uuid.New(),
		VariantName:    "Color",
		VariantValue:   "Red",
		Price_Override: nil,
		Quantity:       10,
		Product:        nil, // Product not loaded
	}

	_, err := variant.GetPrice()

	if err == nil {
		t.Error("GetPrice() should return error when Product is not loaded and no override is set")
	}

	expectedError := "Product not loaded: cannot determine variant price"
	if err.Error() != expectedError {
		t.Errorf("GetPrice() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestProductVariant_GetPrice_ZeroOverride(t *testing.T) {
	zeroPrice := 0.0
	variant := &ProductVariant{
		ID:             uuid.New(),
		ProductID:      uuid.New(),
		VariantName:    "Sample",
		VariantValue:   "Free",
		Price_Override: &zeroPrice,
		Quantity:       5,
	}

	price, err := variant.GetPrice()

	if err != nil {
		t.Fatalf("GetPrice() error = %v, want nil", err)
	}

	if price != 0.0 {
		t.Errorf("GetPrice() = %v, want 0.0 (zero override should be valid)", price)
	}
}

func TestProductVariant_HasPriceOverride(t *testing.T) {
	tests := []struct {
		name           string
		priceOverride  *float64
		expectedResult bool
	}{
		{
			name:           "With override set",
			priceOverride:  floatPtr(99.99),
			expectedResult: true,
		},
		{
			name:           "Without override (nil)",
			priceOverride:  nil,
			expectedResult: false,
		},
		{
			name:           "With zero override",
			priceOverride:  floatPtr(0.0),
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variant := &ProductVariant{
				ID:             uuid.New(),
				ProductID:      uuid.New(),
				VariantName:    "Test",
				VariantValue:   "Value",
				Price_Override: tt.priceOverride,
				Quantity:       10,
			}

			result := variant.HasPriceOverride()

			if result != tt.expectedResult {
				t.Errorf("HasPriceOverride() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestProductVariant_ValidateForCreation_Success(t *testing.T) {
	overridePrice := 49.99
	variant := &ProductVariant{
		VariantName:    "Color",
		VariantValue:   "Red",
		Price_Override: &overridePrice,
		Quantity:       5,
	}

	err := variant.ValidateForCreation()

	if err != nil {
		t.Errorf("ValidateForCreation() error = %v, want nil", err)
	}
}

func TestProductVariant_ValidateForCreation_MissingVariantName(t *testing.T) {
	variant := &ProductVariant{
		VariantName:  "",
		VariantValue: "Red",
		Quantity:     5,
	}

	err := variant.ValidateForCreation()

	if err == nil {
		t.Error("ValidateForCreation() should return error for missing variant name")
	}

	expectedError := "Variant name is required"
	if err.Error() != expectedError {
		t.Errorf("ValidateForCreation() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestProductVariant_ValidateForCreation_MissingVariantValue(t *testing.T) {
	variant := &ProductVariant{
		VariantName:  "Color",
		VariantValue: "",
		Quantity:     5,
	}

	err := variant.ValidateForCreation()

	if err == nil {
		t.Error("ValidateForCreation() should return error for missing variant value")
	}

	expectedError := "Variant value is required"
	if err.Error() != expectedError {
		t.Errorf("ValidateForCreation() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestProductVariant_ValidateForCreation_NegativePriceOverride(t *testing.T) {
	negativePrice := -10.0
	variant := &ProductVariant{
		VariantName:    "Size",
		VariantValue:   "Large",
		Price_Override: &negativePrice,
		Quantity:       5,
	}

	err := variant.ValidateForCreation()

	if err == nil {
		t.Error("ValidateForCreation() should return error for negative price override")
	}

	expectedError := "Variant price override cannot be negative"
	if err.Error() != expectedError {
		t.Errorf("ValidateForCreation() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestProductVariant_ValidateForCreation_NegativeQuantity(t *testing.T) {
	variant := &ProductVariant{
		VariantName:  "Color",
		VariantValue: "Blue",
		Quantity:     -5,
	}

	err := variant.ValidateForCreation()

	if err == nil {
		t.Error("ValidateForCreation() should return error for negative quantity")
	}

	expectedError := "Variant quantity cannot be negative"
	if err.Error() != expectedError {
		t.Errorf("ValidateForCreation() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestProductVariant_ValidateForCreation_ZeroQuantity(t *testing.T) {
	variant := &ProductVariant{
		VariantName:  "Color",
		VariantValue: "Green",
		Quantity:     0,
	}

	err := variant.ValidateForCreation()

	if err == nil {
		t.Error("ValidateForCreation() should return error for zero quantity")
	}

	expectedError := "Variant quantity must be greater than 0 for new variants"
	if err.Error() != expectedError {
		t.Errorf("ValidateForCreation() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestProductVariant_ValidateForCreation_ZeroPriceOverrideIsValid(t *testing.T) {
	zeroPrice := 0.0
	variant := &ProductVariant{
		VariantName:    "Sample",
		VariantValue:   "Free",
		Price_Override: &zeroPrice,
		Quantity:       5,
	}

	err := variant.ValidateForCreation()

	if err != nil {
		t.Errorf("ValidateForCreation() error = %v, want nil (zero price override should be valid)", err)
	}
}

func TestProductVariant_BeforeCreate(t *testing.T) {
	variant := &ProductVariant{
		VariantName:  "Color",
		VariantValue: "Black",
		Quantity:     10,
	}

	err := variant.BeforeCreate(nil)

	if err != nil {
		t.Errorf("BeforeCreate() error = %v, want nil", err)
	}

	if variant.ID == uuid.Nil {
		t.Error("BeforeCreate() should generate UUID for ID")
	}
}

func TestProductVariant_BeforeCreate_PreservesExistingID(t *testing.T) {
	existingID := uuid.New()
	variant := &ProductVariant{
		ID:           existingID,
		VariantName:  "Size",
		VariantValue: "Medium",
		Quantity:     10,
	}

	err := variant.BeforeCreate(nil)

	if err != nil {
		t.Errorf("BeforeCreate() error = %v, want nil", err)
	}

	if variant.ID != existingID {
		t.Errorf("BeforeCreate() changed existing ID from %v to %v", existingID, variant.ID)
	}
}

// Helper function to create float pointer
func floatPtr(f float64) *float64 {
	return &f
}
