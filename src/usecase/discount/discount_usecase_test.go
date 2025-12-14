package discount

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type mockDiscountRepository struct {
	mock.Mock
}

func (m *mockDiscountRepository) Create(ctx context.Context, discount *entity.Discount) error {
	args := m.Called(ctx, discount)
	return args.Error(0)
}

func (m *mockDiscountRepository) Update(ctx context.Context, discount *entity.Discount) error {
	args := m.Called(ctx, discount)
	return args.Error(0)
}

func (m *mockDiscountRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Discount, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Discount), args.Error(1)
}

func (m *mockDiscountRepository) GetByPromoCode(ctx context.Context, promoCode string) (*entity.Discount, error) {
	args := m.Called(ctx, promoCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Discount), args.Error(1)
}

func TestCreateDiscount_Success(t *testing.T) {
	repo := new(mockDiscountRepository)
	uc := NewUseCase(repo)

	discount := &entity.Discount{
		PromoCode:    "SAVE20",
		DiscountType: entity.Percentage,
		Value:        20.0,
		Active:       true,
	}

	repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Discount")).Return(nil)

	err := uc.CreateDiscount(context.Background(), discount)
	if err != nil {
		t.Errorf("CreateDiscount() unexpected error: %v", err)
	}

	repo.AssertExpectations(t)
}

func TestCreateDiscount_ValidationError(t *testing.T) {
	repo := new(mockDiscountRepository)
	uc := NewUseCase(repo)

	discount := &entity.Discount{
		PromoCode:    "", // Invalid: empty promo code
		DiscountType: entity.Percentage,
		Value:        20.0,
		Active:       true,
	}

	err := uc.CreateDiscount(context.Background(), discount)
	if err == nil {
		t.Error("CreateDiscount() expected validation error, got nil")
	}
}

func TestCreateDiscount_RepositoryError(t *testing.T) {
	repo := new(mockDiscountRepository)
	uc := NewUseCase(repo)

	discount := &entity.Discount{
		PromoCode:    "SAVE20",
		DiscountType: entity.Percentage,
		Value:        20.0,
		Active:       true,
	}

	repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Discount")).Return(errors.New("database error"))

	err := uc.CreateDiscount(context.Background(), discount)
	if err == nil {
		t.Error("CreateDiscount() expected repository error, got nil")
	}

	repo.AssertExpectations(t)
}

func TestGetDiscount_Success(t *testing.T) {
	repo := new(mockDiscountRepository)
	uc := NewUseCase(repo)

	discountID := uuid.New()
	expectedDiscount := &entity.Discount{
		ID:           discountID,
		PromoCode:    "SAVE20",
		DiscountType: entity.Percentage,
		Value:        20.0,
		Active:       true,
	}

	repo.On("GetByID", mock.Anything, discountID).Return(expectedDiscount, nil)

	result, err := uc.GetDiscountByID(context.Background(), discountID)
	if err != nil {
		t.Errorf("GetDiscountByID() unexpected error: %v", err)
	}
	if result.ID != discountID {
		t.Errorf("GetDiscountByID() returned wrong discount ID: got %v, want %v", result.ID, discountID)
	}

	repo.AssertExpectations(t)
}

func TestGetDiscount_NotFound(t *testing.T) {
	repo := new(mockDiscountRepository)
	uc := NewUseCase(repo)

	discountID := uuid.New()

	repo.On("GetByID", mock.Anything, discountID).Return(nil, errors.New("not found"))

	_, err := uc.GetDiscountByID(context.Background(), discountID)
	if err == nil {
		t.Error("GetDiscountByID() expected error for not found, got nil")
	}

	repo.AssertExpectations(t)
}

func TestUpdateDiscount_Success(t *testing.T) {
	repo := new(mockDiscountRepository)
	uc := NewUseCase(repo)

	discountID := uuid.New()
	updatedDiscount := &entity.Discount{
		ID:           discountID,
		PromoCode:    "SAVE20",
		DiscountType: entity.Percentage,
		Value:        25.0,  // Updated value
		Active:       false, // Deactivated
	}

	repo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Discount")).Return(nil)

	err := uc.UpdateDiscount(context.Background(), updatedDiscount)
	if err != nil {
		t.Errorf("UpdateDiscount() unexpected error: %v", err)
	}

	repo.AssertExpectations(t)
}

func TestUpdateDiscount_RepositoryError(t *testing.T) {
	repo := new(mockDiscountRepository)
	uc := NewUseCase(repo)

	updatedDiscount := &entity.Discount{
		PromoCode:    "SAVE20",
		DiscountType: entity.Percentage,
		Value:        25.0,
		Active:       false,
	}

	repo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Discount")).Return(errors.New("database error"))

	err := uc.UpdateDiscount(context.Background(), updatedDiscount)
	if err == nil {
		t.Error("UpdateDiscount() expected error for database error, got nil")
	}

	repo.AssertExpectations(t)
}

func TestUpdateDiscount_ValidationError(t *testing.T) {
	repo := new(mockDiscountRepository)
	uc := NewUseCase(repo)

	updatedDiscount := &entity.Discount{
		PromoCode:    "", // Invalid: empty promo code
		DiscountType: entity.Percentage,
		Value:        25.0,
		Active:       false,
	}

	err := uc.UpdateDiscount(context.Background(), updatedDiscount)
	if err == nil {
		t.Error("UpdateDiscount() expected validation error, got nil")
	}
}

func TestValidatePromoCode_Success(t *testing.T) {
	repo := new(mockDiscountRepository)
	uc := NewUseCase(repo)

	promoCode := "SAVE20"
	expectedDiscount := &entity.Discount{
		ID:           uuid.New(),
		PromoCode:    promoCode,
		DiscountType: entity.Percentage,
		Value:        20.0,
		Active:       true,
	}

	repo.On("GetByPromoCode", mock.Anything, promoCode).Return(expectedDiscount, nil)

	result, err := uc.GetDiscountByPromoCode(context.Background(), promoCode)
	if err != nil {
		t.Errorf("GetDiscountByPromoCode() unexpected error: %v", err)
	}
	if result.PromoCode != promoCode {
		t.Errorf("GetDiscountByPromoCode() returned wrong promo code: got %v, want %v", result.PromoCode, promoCode)
	}

	repo.AssertExpectations(t)
}

func TestValidatePromoCode_NotFound(t *testing.T) {
	repo := new(mockDiscountRepository)
	uc := NewUseCase(repo)

	promoCode := "INVALID"

	repo.On("GetByPromoCode", mock.Anything, promoCode).Return(nil, errors.New("not found"))

	_, err := uc.GetDiscountByPromoCode(context.Background(), promoCode)
	if err == nil {
		t.Error("GetDiscountByPromoCode() expected error for not found, got nil")
	}

	repo.AssertExpectations(t)
}

func TestValidatePromoCode_EmptyCode(t *testing.T) {
	repo := new(mockDiscountRepository)
	uc := NewUseCase(repo)

	repo.On("GetByPromoCode", mock.Anything, "").Return(nil, errors.New("not found"))

	_, err := uc.GetDiscountByPromoCode(context.Background(), "")
	if err == nil {
		t.Error("GetDiscountByPromoCode() expected error for empty promo code, got nil")
	}
}
