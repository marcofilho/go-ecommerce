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

func (m *mockDiscountRepository) GetByPromoCodeWithRelations(ctx context.Context, promoCode string) (*entity.Discount, error) {
	args := m.Called(ctx, promoCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Discount), args.Error(1)
}

func (m *mockDiscountRepository) GetUserUsage(ctx context.Context, discountID, userID uuid.UUID) (*entity.DiscountUser, error) {
	args := m.Called(ctx, discountID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.DiscountUser), args.Error(1)
}

func (m *mockDiscountRepository) IncrementUsage(ctx context.Context, discountID uuid.UUID) error {
	args := m.Called(ctx, discountID)
	return args.Error(0)
}

func (m *mockDiscountRepository) IncrementUserUsage(ctx context.Context, discountID, userID uuid.UUID) error {
	args := m.Called(ctx, discountID, userID)
	return args.Error(0)
}

func (m *mockDiscountRepository) AssociateProducts(ctx context.Context, discountID uuid.UUID, productIDs []uuid.UUID) error {
	args := m.Called(ctx, discountID, productIDs)
	return args.Error(0)
}

func (m *mockDiscountRepository) AssociateCategories(ctx context.Context, discountID uuid.UUID, categoryIDs []uuid.UUID) error {
	args := m.Called(ctx, discountID, categoryIDs)
	return args.Error(0)
}

func (m *mockDiscountRepository) AssociateUsers(ctx context.Context, discountID uuid.UUID, userIDs []uuid.UUID, usageLimit *int) error {
	args := m.Called(ctx, discountID, userIDs, usageLimit)
	return args.Error(0)
}

// Mock product repository
type mockProductRepository struct {
	mock.Mock
}

func (m *mockProductRepository) Create(ctx context.Context, product *entity.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *mockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *mockProductRepository) GetAll(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error) {
	args := m.Called(ctx, page, pageSize, inStockOnly)
	return args.Get(0).([]*entity.Product), args.Int(1), args.Error(2)
}

func (m *mockProductRepository) Update(ctx context.Context, product *entity.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *mockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateDiscount_Success(t *testing.T) {
	repo := new(mockDiscountRepository)
	productRepo := new(mockProductRepository)
	uc := NewUseCase(repo, productRepo)

	discount := &entity.Discount{
		PromoCode:    "SAVE20",
		DiscountType: entity.Percentage,
		Value:        20.0,
		Active:       true,
	}

	productIDs := []uuid.UUID{uuid.New()}
	categoryIDs := []uuid.UUID{uuid.New()}
	userIDs := []uuid.UUID{uuid.New()}
	var userLimit *int

	repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Discount")).Return(nil)
	repo.On("AssociateProducts", mock.Anything, mock.Anything, productIDs).Return(nil)
	repo.On("AssociateCategories", mock.Anything, mock.Anything, categoryIDs).Return(nil)
	repo.On("AssociateUsers", mock.Anything, mock.Anything, userIDs, userLimit).Return(nil)

	err := uc.CreateDiscount(context.Background(), discount, productIDs, categoryIDs, userIDs, userLimit)
	if err != nil {
		t.Errorf("CreateDiscount() unexpected error: %v", err)
	}

	repo.AssertExpectations(t)
}

func TestCreateDiscount_ValidationError(t *testing.T) {
	repo := new(mockDiscountRepository)
	productRepo := new(mockProductRepository)
	uc := NewUseCase(repo, productRepo)

	discount := &entity.Discount{
		PromoCode:    "", // Invalid: empty promo code
		DiscountType: entity.Percentage,
		Value:        20.0,
		Active:       true,
	}

	err := uc.CreateDiscount(context.Background(), discount, nil, nil, nil, nil)
	if err == nil {
		t.Error("CreateDiscount() expected validation error, got nil")
	}
}

func TestCreateDiscount_RepositoryError(t *testing.T) {
	repo := new(mockDiscountRepository)
	productRepo := new(mockProductRepository)
	uc := NewUseCase(repo, productRepo)

	discount := &entity.Discount{
		PromoCode:    "SAVE20",
		DiscountType: entity.Percentage,
		Value:        20.0,
		Active:       true,
	}

	repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Discount")).Return(errors.New("database error"))

	err := uc.CreateDiscount(context.Background(), discount, nil, nil, nil, nil)
	if err == nil {
		t.Error("CreateDiscount() expected repository error, got nil")
	}

	repo.AssertExpectations(t)
}

func TestGetDiscount_Success(t *testing.T) {
	repo := new(mockDiscountRepository)
	productRepo := new(mockProductRepository)
	uc := NewUseCase(repo, productRepo)

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
	productRepo := new(mockProductRepository)
	uc := NewUseCase(repo, productRepo)

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
	productRepo := new(mockProductRepository)
	uc := NewUseCase(repo, productRepo)

	discountID := uuid.New()
	updatedDiscount := &entity.Discount{
		ID:           discountID,
		PromoCode:    "SAVE20",
		DiscountType: entity.Percentage,
		Value:        25.0,  // Updated value
		Active:       false, // Deactivated
	}

	repo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Discount")).Return(nil)
	repo.On("AssociateProducts", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	repo.On("AssociateCategories", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	repo.On("AssociateUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := uc.UpdateDiscount(context.Background(), updatedDiscount, nil, nil, nil, nil)
	if err != nil {
		t.Errorf("UpdateDiscount() unexpected error: %v", err)
	}

	repo.AssertExpectations(t)
}

func TestUpdateDiscount_RepositoryError(t *testing.T) {
	repo := new(mockDiscountRepository)
	productRepo := new(mockProductRepository)
	uc := NewUseCase(repo, productRepo)

	updatedDiscount := &entity.Discount{
		PromoCode:    "SAVE20",
		DiscountType: entity.Percentage,
		Value:        25.0,
		Active:       false,
	}

	repo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Discount")).Return(errors.New("database error"))

	err := uc.UpdateDiscount(context.Background(), updatedDiscount, nil, nil, nil, nil)
	if err == nil {
		t.Error("UpdateDiscount() expected error for database error, got nil")
	}

	repo.AssertExpectations(t)
}

func TestUpdateDiscount_ValidationError(t *testing.T) {
	repo := new(mockDiscountRepository)
	productRepo := new(mockProductRepository)
	uc := NewUseCase(repo, productRepo)

	updatedDiscount := &entity.Discount{
		PromoCode:    "", // Invalid: empty promo code
		DiscountType: entity.Percentage,
		Value:        25.0,
		Active:       false,
	}

	err := uc.UpdateDiscount(context.Background(), updatedDiscount, nil, nil, nil, nil)
	if err == nil {
		t.Error("UpdateDiscount() expected validation error, got nil")
	}
}

func TestValidatePromoCode_Success(t *testing.T) {
	repo := new(mockDiscountRepository)
	productRepo := new(mockProductRepository)
	uc := NewUseCase(repo, productRepo)

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
	productRepo := new(mockProductRepository)
	uc := NewUseCase(repo, productRepo)

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
	productRepo := new(mockProductRepository)
	uc := NewUseCase(repo, productRepo)

	repo.On("GetByPromoCode", mock.Anything, "").Return(nil, errors.New("not found"))

	_, err := uc.GetDiscountByPromoCode(context.Background(), "")
	if err == nil {
		t.Error("GetDiscountByPromoCode() expected error for empty promo code, got nil")
	}
}
