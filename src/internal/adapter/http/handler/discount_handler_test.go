package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/dto"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/usecase/discount"
	"github.com/stretchr/testify/mock"
)

// Mock discount service
type mockDiscountService struct {
	mock.Mock
}

func (m *mockDiscountService) CreateDiscount(ctx context.Context, discount *entity.Discount, productIDs, categoryIDs, userIDs []uuid.UUID, userUsageLimit *int) error {
	args := m.Called(ctx, discount, productIDs, categoryIDs, userIDs, userUsageLimit)
	return args.Error(0)
}

func (m *mockDiscountService) GetDiscountByID(ctx context.Context, id uuid.UUID) (*entity.Discount, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Discount), args.Error(1)
}

func (m *mockDiscountService) UpdateDiscount(ctx context.Context, discount *entity.Discount, productIDs, categoryIDs, userIDs []uuid.UUID, userUsageLimit *int) error {
	args := m.Called(ctx, discount, productIDs, categoryIDs, userIDs, userUsageLimit)
	return args.Error(0)
}

func (m *mockDiscountService) GetDiscountByPromoCode(ctx context.Context, promoCode string) (*entity.Discount, error) {
	args := m.Called(ctx, promoCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Discount), args.Error(1)
}

func (m *mockDiscountService) ValidateDiscountForOrder(ctx context.Context, promoCode string, userID uuid.UUID, orderItems []discount.OrderItemInfo, orderTotal float64) (*discount.DiscountValidationResult, error) {
	args := m.Called(ctx, promoCode, userID, orderItems, orderTotal)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*discount.DiscountValidationResult), args.Error(1)
}

func (m *mockDiscountService) ApplyDiscount(ctx context.Context, discountID, userID uuid.UUID) error {
	args := m.Called(ctx, discountID, userID)
	return args.Error(0)
}

func TestDiscountHandler_CreateDiscount_Success(t *testing.T) {
	service := new(mockDiscountService)
	handler := NewDiscountHandler(service)

	request := dto.DiscountRequest{
		PromoCode:    "SAVE20",
		DiscountType: "percentage",
		Value:        20.0,
		Active:       true,
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/api/discounts", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	service.On("CreateDiscount", mock.Anything, mock.AnythingOfType("*entity.Discount"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	handler.CreateDiscount(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	service.AssertExpectations(t)
}

func TestDiscountHandler_CreateDiscount_InvalidJSON(t *testing.T) {
	service := new(mockDiscountService)
	handler := NewDiscountHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/api/discounts", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	handler.CreateDiscount(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestDiscountHandler_CreateDiscount_UseCaseError(t *testing.T) {
	service := new(mockDiscountService)
	handler := NewDiscountHandler(service)

	request := dto.DiscountRequest{
		PromoCode:    "SAVE20",
		DiscountType: "percentage",
		Value:        20.0,
		Active:       true,
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/api/discounts", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	service.On("CreateDiscount", mock.Anything, mock.AnythingOfType("*entity.Discount"), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("service error"))

	handler.CreateDiscount(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	service.AssertExpectations(t)
}

func TestDiscountHandler_GetDiscount_Success(t *testing.T) {
	service := new(mockDiscountService)
	handler := NewDiscountHandler(service)

	discountID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/discounts/"+discountID.String(), nil)
	req.SetPathValue("id", discountID.String())
	rec := httptest.NewRecorder()

	discount := &entity.Discount{
		ID:           discountID,
		PromoCode:    "SAVE20",
		DiscountType: entity.Percentage,
		Value:        20.0,
		Active:       true,
	}

	service.On("GetDiscountByID", mock.Anything, discountID).Return(discount, nil)

	handler.GetDiscount(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	service.AssertExpectations(t)
}

func TestDiscountHandler_GetDiscount_InvalidID(t *testing.T) {
	service := new(mockDiscountService)
	handler := NewDiscountHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/discounts/invalid-id", nil)
	req.SetPathValue("id", "invalid-id")
	rec := httptest.NewRecorder()

	handler.GetDiscount(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestDiscountHandler_GetDiscount_NotFound(t *testing.T) {
	service := new(mockDiscountService)
	handler := NewDiscountHandler(service)

	discountID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/discounts/"+discountID.String(), nil)
	req.SetPathValue("id", discountID.String())
	rec := httptest.NewRecorder()

	service.On("GetDiscountByID", mock.Anything, discountID).Return(nil, errors.New("not found"))

	handler.GetDiscount(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	service.AssertExpectations(t)
}

func TestDiscountHandler_ValidatePromoCode_Success(t *testing.T) {
	service := new(mockDiscountService)
	handler := NewDiscountHandler(service)

	request := dto.ApplyDiscountRequest{
		PromoCode: "SAVE20",
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/api/discounts/validate", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	discount := &entity.Discount{
		ID:           uuid.New(),
		PromoCode:    "SAVE20",
		DiscountType: entity.Percentage,
		Value:        20.0,
		Active:       true,
	}

	service.On("GetDiscountByPromoCode", mock.Anything, "SAVE20").Return(discount, nil)

	handler.ValidatePromoCode(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	service.AssertExpectations(t)
}

func TestDiscountHandler_ValidatePromoCode_NotFound(t *testing.T) {
	service := new(mockDiscountService)
	handler := NewDiscountHandler(service)

	request := dto.ApplyDiscountRequest{
		PromoCode: "INVALID",
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/api/discounts/validate", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	service.On("GetDiscountByPromoCode", mock.Anything, "INVALID").Return(nil, errors.New("not found"))

	handler.ValidatePromoCode(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	service.AssertExpectations(t)
}
