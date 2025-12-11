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
	mockServices "github.com/marcofilho/go-ecommerce/src/internal/testing"
	"github.com/marcofilho/go-ecommerce/src/usecase/product"
)

type mockProductRepo struct {
	createFunc  func(ctx context.Context, product *entity.Product) error
	getByIDFunc func(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	getAllFunc  func(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error)
	updateFunc  func(ctx context.Context, product *entity.Product) error
	deleteFunc  func(ctx context.Context, id uuid.UUID) error
}

func (m *mockProductRepo) Create(ctx context.Context, prod *entity.Product) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, prod)
	}
	return nil
}

func (m *mockProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockProductRepo) GetAll(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx, page, pageSize, inStockOnly)
	}
	return nil, 0, nil
}

func (m *mockProductRepo) Update(ctx context.Context, prod *entity.Product) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, prod)
	}
	return nil
}

func (m *mockProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

var _ repository.ProductRepository = (*mockProductRepo)(nil)

func TestProductHandler_CreateProduct_Success(t *testing.T) {
	mockRepo := &mockProductRepo{
		createFunc: func(ctx context.Context, prod *entity.Product) error {
			return nil
		},
	}

	uc := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	reqBody := dto.ProductRequest{
		Name:        "Laptop",
		Description: "Gaming laptop",
		Price:       999.99,
		Quantity:    10,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	uc.CreateProduct(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var response dto.ProductResponse
	json.NewDecoder(w.Body).Decode(&response)
	if response.Name != "Laptop" {
		t.Errorf("expected name Laptop, got %s", response.Name)
	}
}

func TestProductHandler_CreateProduct_InvalidJSON(t *testing.T) {
	mockRepo := &mockProductRepo{}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_CreateProduct_UseCaseError(t *testing.T) {
	mockRepo := &mockProductRepo{
		createFunc: func(ctx context.Context, prod *entity.Product) error {
			return errors.New("validation error")
		},
	}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	reqBody := dto.ProductRequest{Name: "", Price: -10, Quantity: 0}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_GetProduct_Success(t *testing.T) {
	productID := uuid.New()
	mockRepo := &mockProductRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
			return &entity.Product{
				ID:        id,
				Name:      "Laptop",
				Price:     999.99,
				Quantity:  10,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	req := httptest.NewRequest(http.MethodGet, "/products/"+productID.String(), nil)
	req.SetPathValue("id", productID.String())
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response dto.ProductResponse
	json.NewDecoder(w.Body).Decode(&response)
	if response.ID != productID.String() {
		t.Errorf("expected ID %s, got %s", productID.String(), response.ID)
	}
}

func TestProductHandler_GetProduct_InvalidID(t *testing.T) {
	mockRepo := &mockProductRepo{}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	req := httptest.NewRequest(http.MethodGet, "/products/invalid-id", nil)
	req.SetPathValue("id", "invalid-id")
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_GetProduct_NotFound(t *testing.T) {
	mockRepo := &mockProductRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
			return nil, errors.New("not found")
		},
	}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	productID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/products/"+productID.String(), nil)
	req.SetPathValue("id", productID.String())
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestProductHandler_ListProducts_Success(t *testing.T) {
	mockRepo := &mockProductRepo{
		getAllFunc: func(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error) {
			return []*entity.Product{
				{ID: uuid.New(), Name: "P1", Price: 100, Quantity: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				{ID: uuid.New(), Name: "P2", Price: 200, Quantity: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			}, 2, nil
		},
	}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	req := httptest.NewRequest(http.MethodGet, "/products?page=1&page_size=10&in_stock_only=true", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response dto.ProductListResponse
	json.NewDecoder(w.Body).Decode(&response)
	if len(response.Data) != 2 {
		t.Errorf("expected 2 products, got %d", len(response.Data))
	}
}

func TestProductHandler_ListProducts_InStockOnlyFalse(t *testing.T) {
	mockRepo := &mockProductRepo{
		getAllFunc: func(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error) {
			if inStockOnly {
				t.Error("expected inStockOnly to be false")
			}
			return []*entity.Product{}, 0, nil
		},
	}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	req := httptest.NewRequest(http.MethodGet, "/products?in_stock_only=false", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestProductHandler_ListProducts_UseCaseError(t *testing.T) {
	mockRepo := &mockProductRepo{
		getAllFunc: func(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error) {
			return nil, 0, errors.New("database error")
		},
	}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestProductHandler_UpdateProduct_Success(t *testing.T) {
	productID := uuid.New()
	mockRepo := &mockProductRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
			return &entity.Product{
				ID:        id,
				Name:      "Old",
				Price:     100,
				Quantity:  5,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
		updateFunc: func(ctx context.Context, prod *entity.Product) error {
			return nil
		},
	}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	reqBody := dto.ProductRequest{
		Name:        "Updated Laptop",
		Description: "Updated description",
		Price:       1299.99,
		Quantity:    20,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/products/"+productID.String(), bytes.NewBuffer(body))
	req.SetPathValue("id", productID.String())
	w := httptest.NewRecorder()

	handler.UpdateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response dto.ProductResponse
	json.NewDecoder(w.Body).Decode(&response)
	if response.Name != "Updated Laptop" {
		t.Errorf("expected name 'Updated Laptop', got %s", response.Name)
	}
}

func TestProductHandler_UpdateProduct_InvalidID(t *testing.T) {
	mockRepo := &mockProductRepo{}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	reqBody := dto.ProductRequest{Name: "Updated"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/products/invalid-id", bytes.NewBuffer(body))
	req.SetPathValue("id", "invalid-id")
	w := httptest.NewRecorder()

	handler.UpdateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_UpdateProduct_InvalidJSON(t *testing.T) {
	productID := uuid.New()
	mockRepo := &mockProductRepo{}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	req := httptest.NewRequest(http.MethodPut, "/products/"+productID.String(), bytes.NewBuffer([]byte("invalid")))
	req.SetPathValue("id", productID.String())
	w := httptest.NewRecorder()

	handler.UpdateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_UpdateProduct_UseCaseError(t *testing.T) {
	productID := uuid.New()
	mockRepo := &mockProductRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
			return nil, errors.New("not found")
		},
	}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	reqBody := dto.ProductRequest{Name: "Test"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/products/"+productID.String(), bytes.NewBuffer(body))
	req.SetPathValue("id", productID.String())
	w := httptest.NewRecorder()

	handler.UpdateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_DeleteProduct_Success(t *testing.T) {
	productID := uuid.New()
	mockRepo := &mockProductRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
			return &entity.Product{ID: productID, Name: "Test Product", Price: 100}, nil
		},
		deleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return nil
		},
	}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	req := httptest.NewRequest(http.MethodDelete, "/products/"+productID.String(), nil)
	req.SetPathValue("id", productID.String())
	w := httptest.NewRecorder()

	handler.DeleteProduct(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestProductHandler_DeleteProduct_InvalidID(t *testing.T) {
	mockRepo := &mockProductRepo{}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	req := httptest.NewRequest(http.MethodDelete, "/products/invalid-id", nil)
	req.SetPathValue("id", "invalid-id")
	w := httptest.NewRecorder()

	handler.DeleteProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestProductHandler_DeleteProduct_NotFound(t *testing.T) {
	productID := uuid.New()
	mockRepo := &mockProductRepo{
		deleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return errors.New("not found")
		},
	}
	handler := NewProductHandler(product.NewUseCase(mockRepo, &mockServices.MockServices{}))

	req := httptest.NewRequest(http.MethodDelete, "/products/"+productID.String(), nil)
	req.SetPathValue("id", productID.String())
	w := httptest.NewRecorder()

	handler.DeleteProduct(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
