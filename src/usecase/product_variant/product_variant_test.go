package productvariant

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductVariantRepository is a mock implementation of ProductVariantRepository
type MockProductVariantRepository struct {
	mock.Mock
}

func (m *MockProductVariantRepository) Create(ctx context.Context, variant *entity.ProductVariant) error {
	args := m.Called(ctx, variant)
	return args.Error(0)
}

func (m *MockProductVariantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductVariant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ProductVariant), args.Error(1)
}

func (m *MockProductVariantRepository) GetAll(ctx context.Context, page, pageSize int) ([]*entity.ProductVariant, int, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*entity.ProductVariant), args.Int(1), args.Error(2)
}

func (m *MockProductVariantRepository) GetAllByProductID(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]*entity.ProductVariant, int, error) {
	args := m.Called(ctx, productID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*entity.ProductVariant), args.Int(1), args.Error(2)
}

func (m *MockProductVariantRepository) Update(ctx context.Context, variant *entity.ProductVariant) error {
	args := m.Called(ctx, variant)
	return args.Error(0)
}

func (m *MockProductVariantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateProductVariant(t *testing.T) {
	mockRepo := new(MockProductVariantRepository)
	useCase := NewUseCase(mockRepo)
	ctx := context.Background()

	productID := uuid.New()
	priceOverride := 39.99

	t.Run("Success - Create variant with price override", func(t *testing.T) {
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.ProductVariant")).Return(nil).Once()

		variant, err := useCase.CreateProductVariant(ctx, productID, "Size", "Large", &priceOverride, 50)

		assert.NoError(t, err)
		assert.NotNil(t, variant)
		assert.Equal(t, productID, variant.ProductID)
		assert.Equal(t, "Size", variant.VariantName)
		assert.Equal(t, "Large", variant.VariantValue)
		assert.Equal(t, &priceOverride, variant.Price_Override)
		assert.Equal(t, 50, variant.Quantity)
		assert.NotEqual(t, uuid.Nil, variant.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Create variant without price override", func(t *testing.T) {
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.ProductVariant")).Return(nil).Once()

		variant, err := useCase.CreateProductVariant(ctx, productID, "Color", "Blue", nil, 100)

		assert.NoError(t, err)
		assert.NotNil(t, variant)
		assert.Equal(t, productID, variant.ProductID)
		assert.Equal(t, "Color", variant.VariantName)
		assert.Equal(t, "Blue", variant.VariantValue)
		assert.Nil(t, variant.Price_Override)
		assert.Equal(t, 100, variant.Quantity)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Invalid variant name (empty)", func(t *testing.T) {
		variant, err := useCase.CreateProductVariant(ctx, productID, "", "Medium", nil, 30)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "Variant name is required")
	})

	t.Run("Failure - Invalid variant value (empty)", func(t *testing.T) {
		variant, err := useCase.CreateProductVariant(ctx, productID, "Size", "", nil, 30)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "Variant value is required")
	})

	t.Run("Failure - Invalid quantity (negative)", func(t *testing.T) {
		variant, err := useCase.CreateProductVariant(ctx, productID, "Size", "Small", nil, -10)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "Variant quantity cannot be negative")
	})

	t.Run("Failure - Invalid price override (negative)", func(t *testing.T) {
		negativePriceOverride := -10.00
		variant, err := useCase.CreateProductVariant(ctx, productID, "Size", "Medium", &negativePriceOverride, 20)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "Variant price override cannot be negative")
	})

	t.Run("Failure - Repository error", func(t *testing.T) {
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.ProductVariant")).Return(errors.New("database error")).Once()

		variant, err := useCase.CreateProductVariant(ctx, productID, "Color", "Red", nil, 25)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetProductVariant(t *testing.T) {
	mockRepo := new(MockProductVariantRepository)
	useCase := NewUseCase(mockRepo)
	ctx := context.Background()

	variantID := uuid.New()
	productID := uuid.New()
	priceOverride := 45.99

	t.Run("Success - Get existing variant", func(t *testing.T) {
		expectedVariant := &entity.ProductVariant{
			ID:             variantID,
			ProductID:      productID,
			VariantName:    "Material",
			VariantValue:   "Cotton",
			Price_Override: &priceOverride,
			Quantity:       75,
		}

		mockRepo.On("GetByID", ctx, variantID).Return(expectedVariant, nil).Once()

		variant, err := useCase.GetProductVariant(ctx, variantID)

		assert.NoError(t, err)
		assert.NotNil(t, variant)
		assert.Equal(t, variantID, variant.ID)
		assert.Equal(t, productID, variant.ProductID)
		assert.Equal(t, "Material", variant.VariantName)
		assert.Equal(t, "Cotton", variant.VariantValue)
		assert.Equal(t, 75, variant.Quantity)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Variant not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, variantID).Return(nil, errors.New("variant not found")).Once()

		variant, err := useCase.GetProductVariant(ctx, variantID)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "variant not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Repository error", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, variantID).Return(nil, errors.New("database connection error")).Once()

		variant, err := useCase.GetProductVariant(ctx, variantID)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "database connection error")
		mockRepo.AssertExpectations(t)
	})
}

func TestListProductVariants(t *testing.T) {
	mockRepo := new(MockProductVariantRepository)
	useCase := NewUseCase(mockRepo)
	ctx := context.Background()

	productID := uuid.New()
	priceOverride1 := 29.99
	priceOverride2 := 34.99

	t.Run("Success - List variants with pagination", func(t *testing.T) {
		expectedVariants := []*entity.ProductVariant{
			{
				ID:             uuid.New(),
				ProductID:      productID,
				VariantName:    "Size",
				VariantValue:   "Small",
				Price_Override: &priceOverride1,
				Quantity:       20,
			},
			{
				ID:             uuid.New(),
				ProductID:      productID,
				VariantName:    "Size",
				VariantValue:   "Large",
				Price_Override: &priceOverride2,
				Quantity:       30,
			},
		}

		mockRepo.On("GetAllByProductID", ctx, productID, 1, 10).Return(expectedVariants, 2, nil).Once()

		variants, total, err := useCase.ListProductVariants(ctx, productID, 1, 10)

		assert.NoError(t, err)
		assert.NotNil(t, variants)
		assert.Len(t, variants, 2)
		assert.Equal(t, 2, total)
		assert.Equal(t, "Small", variants[0].VariantValue)
		assert.Equal(t, "Large", variants[1].VariantValue)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Empty list", func(t *testing.T) {
		mockRepo.On("GetAllByProductID", ctx, productID, 1, 10).Return([]*entity.ProductVariant{}, 0, nil).Once()

		variants, total, err := useCase.ListProductVariants(ctx, productID, 1, 10)

		assert.NoError(t, err)
		assert.NotNil(t, variants)
		assert.Len(t, variants, 0)
		assert.Equal(t, 0, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Default page to 1 when invalid", func(t *testing.T) {
		mockRepo.On("GetAllByProductID", ctx, productID, 1, 10).Return([]*entity.ProductVariant{}, 0, nil).Once()

		_, _, err := useCase.ListProductVariants(ctx, productID, 0, 10)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Default pageSize to 10 when invalid (too small)", func(t *testing.T) {
		mockRepo.On("GetAllByProductID", ctx, productID, 1, 10).Return([]*entity.ProductVariant{}, 0, nil).Once()

		_, _, err := useCase.ListProductVariants(ctx, productID, 1, 0)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Default pageSize to 10 when invalid (too large)", func(t *testing.T) {
		mockRepo.On("GetAllByProductID", ctx, productID, 1, 10).Return([]*entity.ProductVariant{}, 0, nil).Once()

		_, _, err := useCase.ListProductVariants(ctx, productID, 1, 150)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Repository error", func(t *testing.T) {
		mockRepo.On("GetAllByProductID", ctx, productID, 1, 10).Return(nil, 0, errors.New("database error")).Once()

		variants, total, err := useCase.ListProductVariants(ctx, productID, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, variants)
		assert.Equal(t, 0, total)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateProductVariant(t *testing.T) {
	mockRepo := new(MockProductVariantRepository)
	useCase := NewUseCase(mockRepo)
	ctx := context.Background()

	variantID := uuid.New()
	productID := uuid.New()
	oldPriceOverride := 29.99
	newPriceOverride := 34.99

	t.Run("Success - Update variant with price override", func(t *testing.T) {
		existingVariant := &entity.ProductVariant{
			ID:             variantID,
			ProductID:      productID,
			VariantName:    "Size",
			VariantValue:   "Small",
			Price_Override: &oldPriceOverride,
			Quantity:       20,
		}

		mockRepo.On("GetByID", ctx, variantID).Return(existingVariant, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.ProductVariant")).Return(nil).Once()

		variant, err := useCase.UpdateProductVariant(ctx, variantID, "Size", "Medium", &newPriceOverride, 50)

		assert.NoError(t, err)
		assert.NotNil(t, variant)
		assert.Equal(t, variantID, variant.ID)
		assert.Equal(t, "Medium", variant.VariantValue)
		assert.Equal(t, &newPriceOverride, variant.Price_Override)
		assert.Equal(t, 50, variant.Quantity)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Update variant removing price override", func(t *testing.T) {
		existingVariant := &entity.ProductVariant{
			ID:             variantID,
			ProductID:      productID,
			VariantName:    "Size",
			VariantValue:   "Large",
			Price_Override: &oldPriceOverride,
			Quantity:       30,
		}

		mockRepo.On("GetByID", ctx, variantID).Return(existingVariant, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.ProductVariant")).Return(nil).Once()

		variant, err := useCase.UpdateProductVariant(ctx, variantID, "Size", "Large", nil, 35)

		assert.NoError(t, err)
		assert.NotNil(t, variant)
		assert.Nil(t, variant.Price_Override)
		assert.Equal(t, 35, variant.Quantity)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Variant not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, variantID).Return(nil, errors.New("variant not found")).Once()

		variant, err := useCase.UpdateProductVariant(ctx, variantID, "Size", "XL", nil, 10)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "variant not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Invalid variant name after update", func(t *testing.T) {
		existingVariant := &entity.ProductVariant{
			ID:             variantID,
			ProductID:      productID,
			VariantName:    "Size",
			VariantValue:   "Small",
			Price_Override: nil,
			Quantity:       20,
		}

		mockRepo.On("GetByID", ctx, variantID).Return(existingVariant, nil).Once()

		variant, err := useCase.UpdateProductVariant(ctx, variantID, "", "Medium", nil, 25)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "Variant name is required")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Invalid variant value after update", func(t *testing.T) {
		existingVariant := &entity.ProductVariant{
			ID:             variantID,
			ProductID:      productID,
			VariantName:    "Size",
			VariantValue:   "Small",
			Price_Override: nil,
			Quantity:       20,
		}

		mockRepo.On("GetByID", ctx, variantID).Return(existingVariant, nil).Once()

		variant, err := useCase.UpdateProductVariant(ctx, variantID, "Size", "", nil, 25)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "Variant value is required")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Invalid quantity after update", func(t *testing.T) {
		existingVariant := &entity.ProductVariant{
			ID:             variantID,
			ProductID:      productID,
			VariantName:    "Size",
			VariantValue:   "Small",
			Price_Override: nil,
			Quantity:       20,
		}

		mockRepo.On("GetByID", ctx, variantID).Return(existingVariant, nil).Once()

		variant, err := useCase.UpdateProductVariant(ctx, variantID, "Size", "Medium", nil, -5)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "Variant quantity cannot be negative")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Invalid price override after update", func(t *testing.T) {
		existingVariant := &entity.ProductVariant{
			ID:             variantID,
			ProductID:      productID,
			VariantName:    "Size",
			VariantValue:   "Small",
			Price_Override: nil,
			Quantity:       20,
		}

		negativePriceOverride := -15.00
		mockRepo.On("GetByID", ctx, variantID).Return(existingVariant, nil).Once()

		variant, err := useCase.UpdateProductVariant(ctx, variantID, "Size", "Medium", &negativePriceOverride, 25)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "Variant price override cannot be negative")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Repository update error", func(t *testing.T) {
		existingVariant := &entity.ProductVariant{
			ID:             variantID,
			ProductID:      productID,
			VariantName:    "Size",
			VariantValue:   "Small",
			Price_Override: nil,
			Quantity:       20,
		}

		mockRepo.On("GetByID", ctx, variantID).Return(existingVariant, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.ProductVariant")).Return(errors.New("database error")).Once()

		variant, err := useCase.UpdateProductVariant(ctx, variantID, "Size", "Medium", nil, 25)

		assert.Error(t, err)
		assert.Nil(t, variant)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteProductVariant(t *testing.T) {
	mockRepo := new(MockProductVariantRepository)
	useCase := NewUseCase(mockRepo)
	ctx := context.Background()

	variantID := uuid.New()

	t.Run("Success - Delete existing variant", func(t *testing.T) {
		mockRepo.On("Delete", ctx, variantID).Return(nil).Once()

		err := useCase.DeleteProductVariant(ctx, variantID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Variant not found", func(t *testing.T) {
		mockRepo.On("Delete", ctx, variantID).Return(errors.New("variant not found")).Once()

		err := useCase.DeleteProductVariant(ctx, variantID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "variant not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Repository error", func(t *testing.T) {
		mockRepo.On("Delete", ctx, variantID).Return(errors.New("database error")).Once()

		err := useCase.DeleteProductVariant(ctx, variantID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}
