package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/dto"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	"github.com/marcofilho/go-ecommerce/src/usecase/order"
)

type mockOrderRepo struct {
	createFunc  func(ctx context.Context, order *entity.Order) error
	getByIDFunc func(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	getAllFunc  func(ctx context.Context, page, pageSize int, status *entity.OrderStatus, paymentStatus *entity.PaymentStatus) ([]*entity.Order, int, error)
	updateFunc  func(ctx context.Context, order *entity.Order) error
}

func (m *mockOrderRepo) Create(ctx context.Context, order *entity.Order) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, order)
	}
	return nil
}

func (m *mockOrderRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockOrderRepo) GetAll(ctx context.Context, page, pageSize int, status *entity.OrderStatus, paymentStatus *entity.PaymentStatus) ([]*entity.Order, int, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx, page, pageSize, status, paymentStatus)
	}
	return nil, 0, nil
}

func (m *mockOrderRepo) Update(ctx context.Context, order *entity.Order) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, order)
	}
	return nil
}

var _ repository.OrderRepository = (*mockOrderRepo)(nil)

func TestOrderHandler_CreateOrder_Success(t *testing.T) {
	productID := uuid.New()
	mockOrderRepo := &mockOrderRepo{}
	mockProductRepo := &mockProductRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
			return &entity.Product{
				ID: id, Name: "Laptop", Price: 999.99, Quantity: 10,
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}, nil
		},
		updateFunc: func(ctx context.Context, product *entity.Product) error {
			return nil
		},
	}

	handler := NewOrderHandler(newOrderUseCase(mockOrderRepo, mockProductRepo))

	reqBody := dto.CreateOrderRequest{
		CustomerID: 123,
		Products: []dto.OrderItemRequest{
			{ProductID: productID.String(), Quantity: 2},
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateOrder(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestOrderHandler_CreateOrder_InvalidJSON(t *testing.T) {
	handler := NewOrderHandler(newOrderUseCase(&mockOrderRepo{}, &mockProductRepo{}))

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer([]byte("invalid")))
	w := httptest.NewRecorder()

	handler.CreateOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestOrderHandler_CreateOrder_InvalidProductID(t *testing.T) {
	handler := NewOrderHandler(newOrderUseCase(&mockOrderRepo{}, &mockProductRepo{}))

	reqBody := dto.CreateOrderRequest{
		CustomerID: 123,
		Products: []dto.OrderItemRequest{
			{ProductID: "invalid-uuid", Quantity: 2},
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestOrderHandler_CreateOrder_UseCaseError(t *testing.T) {
	productID := uuid.New()
	mockOrderRepo := &mockOrderRepo{}
	mockProductRepo := &mockProductRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
			return nil, errors.New("not found")
		},
	}

	handler := NewOrderHandler(newOrderUseCase(mockOrderRepo, mockProductRepo))

	reqBody := dto.CreateOrderRequest{
		CustomerID: 123,
		Products: []dto.OrderItemRequest{
			{ProductID: productID.String(), Quantity: 2},
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestOrderHandler_GetOrder_Success(t *testing.T) {
	orderID := uuid.New()
	mockOrderRepo := &mockOrderRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
			return &entity.Order{
				ID:            id,
				CustomerID:    123,
				Status:        entity.Pending,
				PaymentStatus: entity.Unpaid,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}, nil
		},
	}

	handler := NewOrderHandler(newOrderUseCase(mockOrderRepo, &mockProductRepo{}))

	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
	req.SetPathValue("id", orderID.String())
	w := httptest.NewRecorder()

	handler.GetOrder(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestOrderHandler_GetOrder_InvalidID(t *testing.T) {
	handler := NewOrderHandler(newOrderUseCase(&mockOrderRepo{}, &mockProductRepo{}))

	req := httptest.NewRequest(http.MethodGet, "/orders/invalid-id", nil)
	req.SetPathValue("id", "invalid-id")
	w := httptest.NewRecorder()

	handler.GetOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestOrderHandler_GetOrder_NotFound(t *testing.T) {
	mockOrderRepo := &mockOrderRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
			return nil, errors.New("not found")
		},
	}

	handler := NewOrderHandler(newOrderUseCase(mockOrderRepo, &mockProductRepo{}))

	orderID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
	req.SetPathValue("id", orderID.String())
	w := httptest.NewRecorder()

	handler.GetOrder(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestOrderHandler_ListOrders_Success(t *testing.T) {
	mockOrderRepo := &mockOrderRepo{
		getAllFunc: func(ctx context.Context, page, pageSize int, status *entity.OrderStatus, paymentStatus *entity.PaymentStatus) ([]*entity.Order, int, error) {
			return []*entity.Order{
				{ID: uuid.New(), CustomerID: 1, Status: entity.Pending, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				{ID: uuid.New(), CustomerID: 2, Status: entity.Completed, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			}, 2, nil
		},
	}

	handler := NewOrderHandler(newOrderUseCase(mockOrderRepo, &mockProductRepo{}))

	req := httptest.NewRequest(http.MethodGet, "/orders?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	handler.ListOrders(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response dto.OrderListResponse
	json.NewDecoder(w.Body).Decode(&response)
	if len(response.Data) != 2 {
		t.Errorf("expected 2 orders, got %d", len(response.Data))
	}
}

func TestOrderHandler_ListOrders_WithFilters(t *testing.T) {
	mockOrderRepo := &mockOrderRepo{
		getAllFunc: func(ctx context.Context, page, pageSize int, status *entity.OrderStatus, paymentStatus *entity.PaymentStatus) ([]*entity.Order, int, error) {
			if status == nil {
				t.Error("expected status filter to be set")
			}
			if *status != entity.Pending {
				t.Errorf("expected status pending, got %s", *status)
			}
			return []*entity.Order{}, 0, nil
		},
	}

	handler := NewOrderHandler(newOrderUseCase(mockOrderRepo, &mockProductRepo{}))

	req := httptest.NewRequest(http.MethodGet, "/orders?status=pending&payment_status=unpaid", nil)
	w := httptest.NewRecorder()

	handler.ListOrders(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestOrderHandler_ListOrders_UseCaseError(t *testing.T) {
	mockOrderRepo := &mockOrderRepo{
		getAllFunc: func(ctx context.Context, page, pageSize int, status *entity.OrderStatus, paymentStatus *entity.PaymentStatus) ([]*entity.Order, int, error) {
			return nil, 0, errors.New("database error")
		},
	}

	handler := NewOrderHandler(newOrderUseCase(mockOrderRepo, &mockProductRepo{}))

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()

	handler.ListOrders(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestOrderHandler_UpdateOrderStatus_Success(t *testing.T) {
	orderID := uuid.New()
	mockOrderRepo := &mockOrderRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
			return &entity.Order{
				ID:         id,
				CustomerID: 123,
				Status:     entity.Pending,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}, nil
		},
		updateFunc: func(ctx context.Context, order *entity.Order) error {
			return nil
		},
	}

	handler := NewOrderHandler(newOrderUseCase(mockOrderRepo, &mockProductRepo{}))

	reqBody := dto.UpdateOrderStatusRequest{Status: string(entity.Completed)}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/orders/"+orderID.String()+"/status", bytes.NewBuffer(body))
	req.SetPathValue("id", orderID.String())
	w := httptest.NewRecorder()

	handler.UpdateOrderStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestOrderHandler_UpdateOrderStatus_InvalidID(t *testing.T) {
	handler := NewOrderHandler(newOrderUseCase(&mockOrderRepo{}, &mockProductRepo{}))

	reqBody := dto.UpdateOrderStatusRequest{Status: string(entity.Completed)}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/orders/invalid-id/status", bytes.NewBuffer(body))
	req.SetPathValue("id", "invalid-id")
	w := httptest.NewRecorder()

	handler.UpdateOrderStatus(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestOrderHandler_UpdateOrderStatus_InvalidJSON(t *testing.T) {
	orderID := uuid.New()
	handler := NewOrderHandler(newOrderUseCase(&mockOrderRepo{}, &mockProductRepo{}))

	req := httptest.NewRequest(http.MethodPut, "/orders/"+orderID.String()+"/status", bytes.NewBuffer([]byte("invalid")))
	req.SetPathValue("id", orderID.String())
	w := httptest.NewRecorder()

	handler.UpdateOrderStatus(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestOrderHandler_UpdateOrderStatus_UseCaseError(t *testing.T) {
	orderID := uuid.New()
	mockOrderRepo := &mockOrderRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
			return &entity.Order{
				ID:         id,
				CustomerID: 123,
				Status:     entity.Completed,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}, nil
		},
	}

	handler := NewOrderHandler(newOrderUseCase(mockOrderRepo, &mockProductRepo{}))

	reqBody := dto.UpdateOrderStatusRequest{Status: string(entity.Cancelled)}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/orders/"+orderID.String()+"/status", bytes.NewBuffer(body))
	req.SetPathValue("id", orderID.String())
	w := httptest.NewRecorder()

	handler.UpdateOrderStatus(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func newOrderUseCase(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) *order.UseCase {
	// Create a mock variant repo for testing
	variantRepo := &mockVariantRepo{}
	return order.NewUseCase(orderRepo, productRepo, variantRepo)
}

// Mock variant repository for testing
type mockVariantRepo struct{}

func (m *mockVariantRepo) Create(ctx context.Context, variant *entity.ProductVariant) error {
	return nil
}

func (m *mockVariantRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductVariant, error) {
	return nil, errors.New("variant not found")
}

func (m *mockVariantRepo) GetAll(ctx context.Context, page, pageSize int) ([]*entity.ProductVariant, int, error) {
	return nil, 0, nil
}

func (m *mockVariantRepo) GetAllByProductID(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]*entity.ProductVariant, int, error) {
	return nil, 0, nil
}

func (m *mockVariantRepo) Update(ctx context.Context, variant *entity.ProductVariant) error {
	return nil
}

func (m *mockVariantRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
