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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/dto"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

// MockCategoryService is a mock implementation of category.CategoryService
type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) CreateCategory(ctx context.Context, name string) (*entity.Category, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *MockCategoryService) GetCategory(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *MockCategoryService) ListCategories(ctx context.Context, page, pageSize int) ([]*entity.Category, int, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*entity.Category), args.Get(1).(int), args.Error(2)
}

func (m *MockCategoryService) UpdateCategory(ctx context.Context, id uuid.UUID, name string) (*entity.Category, error) {
	args := m.Called(ctx, id, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *MockCategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryService) AssignCategoryToProduct(ctx context.Context, productID, categoryID uuid.UUID) error {
	args := m.Called(ctx, productID, categoryID)
	return args.Error(0)
}

func (m *MockCategoryService) RemoveCategoryFromProduct(ctx context.Context, productID, categoryID uuid.UUID) error {
	args := m.Called(ctx, productID, categoryID)
	return args.Error(0)
}

func (m *MockCategoryService) GetProductCategories(ctx context.Context, productID uuid.UUID) ([]*entity.Category, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).([]*entity.Category), args.Error(1)
}

func TestCategoryHandler_CreateCategory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		categoryID := uuid.New()
		expectedCategory := &entity.Category{
			ID:   categoryID,
			Name: "Electronics",
		}

		reqBody := dto.CategoryRequest{
			Name: "Electronics",
		}
		body, _ := json.Marshal(reqBody)

		mockService.On("CreateCategory", mock.Anything, "Electronics").Return(expectedCategory, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.CreateCategory(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response dto.CategoryResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, categoryID.String(), response.ID)
		assert.Equal(t, "Electronics", response.Name)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.CreateCategory(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "CreateCategory")
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		reqBody := dto.CategoryRequest{
			Name: "Electronics",
		}
		body, _ := json.Marshal(reqBody)

		mockService.On("CreateCategory", mock.Anything, "Electronics").Return(nil, errors.New("database error"))

		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.CreateCategory(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCategoryHandler_ListCategories(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		categories := []*entity.Category{
			{ID: uuid.New(), Name: "Electronics"},
			{ID: uuid.New(), Name: "Clothing"},
		}

		mockService.On("ListCategories", mock.Anything, 1, 10).Return(categories, 2, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/categories?page=1&page_size=10", nil)
		w := httptest.NewRecorder()

		handler.ListCategories(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data     []dto.CategoryResponse `json:"data"`
			Total    int                    `json:"total"`
			Page     int                    `json:"page"`
			PageSize int                    `json:"page_size"`
		}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response.Data, 2)
		assert.Equal(t, 2, response.Total)

		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		mockService.On("ListCategories", mock.Anything, 1, 10).Return([]*entity.Category{}, 0, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
		w := httptest.NewRecorder()

		handler.ListCategories(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCategoryHandler_AssignCategoryToProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		productID := uuid.New()
		categoryID := uuid.New()

		reqBody := dto.AssignCategoryRequest{
			CategoryID: categoryID.String(),
		}
		body, _ := json.Marshal(reqBody)

		mockService.On("AssignCategoryToProduct", mock.Anything, productID, categoryID).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/api/products/"+productID.String()+"/categories", bytes.NewReader(body))
		req.SetPathValue("id", productID.String())
		w := httptest.NewRecorder()

		handler.AssignCategoryToProduct(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Product ID", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		reqBody := dto.AssignCategoryRequest{
			CategoryID: uuid.New().String(),
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/products/invalid/categories", bytes.NewReader(body))
		req.SetPathValue("id", "invalid")
		w := httptest.NewRecorder()

		handler.AssignCategoryToProduct(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "AssignCategoryToProduct")
	})

	t.Run("Invalid Category ID", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		productID := uuid.New()

		reqBody := dto.AssignCategoryRequest{
			CategoryID: "invalid",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/products/"+productID.String()+"/categories", bytes.NewReader(body))
		req.SetPathValue("id", productID.String())
		w := httptest.NewRecorder()

		handler.AssignCategoryToProduct(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "AssignCategoryToProduct")
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		productID := uuid.New()
		categoryID := uuid.New()

		reqBody := dto.AssignCategoryRequest{
			CategoryID: categoryID.String(),
		}
		body, _ := json.Marshal(reqBody)

		mockService.On("AssignCategoryToProduct", mock.Anything, productID, categoryID).Return(errors.New("already assigned"))

		req := httptest.NewRequest(http.MethodPost, "/api/products/"+productID.String()+"/categories", bytes.NewReader(body))
		req.SetPathValue("id", productID.String())
		w := httptest.NewRecorder()

		handler.AssignCategoryToProduct(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCategoryHandler_RemoveCategoryFromProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		productID := uuid.New()
		categoryID := uuid.New()

		mockService.On("RemoveCategoryFromProduct", mock.Anything, productID, categoryID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/products/"+productID.String()+"/categories/"+categoryID.String(), nil)
		req.SetPathValue("id", productID.String())
		req.SetPathValue("category_id", categoryID.String())
		w := httptest.NewRecorder()

		handler.RemoveCategoryFromProduct(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Product ID", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		categoryID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/api/products/invalid/categories/"+categoryID.String(), nil)
		req.SetPathValue("id", "invalid")
		req.SetPathValue("category_id", categoryID.String())
		w := httptest.NewRecorder()

		handler.RemoveCategoryFromProduct(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "RemoveCategoryFromProduct")
	})

	t.Run("Invalid Category ID", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		productID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/api/products/"+productID.String()+"/categories/invalid", nil)
		req.SetPathValue("id", productID.String())
		req.SetPathValue("category_id", "invalid")
		w := httptest.NewRecorder()

		handler.RemoveCategoryFromProduct(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "RemoveCategoryFromProduct")
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		productID := uuid.New()
		categoryID := uuid.New()

		mockService.On("RemoveCategoryFromProduct", mock.Anything, productID, categoryID).Return(errors.New("not found"))

		req := httptest.NewRequest(http.MethodDelete, "/api/products/"+productID.String()+"/categories/"+categoryID.String(), nil)
		req.SetPathValue("id", productID.String())
		req.SetPathValue("category_id", categoryID.String())
		w := httptest.NewRecorder()

		handler.RemoveCategoryFromProduct(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCategoryHandler_GetProductCategories(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		productID := uuid.New()
		categories := []*entity.Category{
			{ID: uuid.New(), Name: "Electronics"},
			{ID: uuid.New(), Name: "Computers"},
		}

		mockService.On("GetProductCategories", mock.Anything, productID).Return(categories, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/products/"+productID.String()+"/categories", nil)
		req.SetPathValue("id", productID.String())
		w := httptest.NewRecorder()

		handler.GetProductCategories(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []dto.CategoryResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Product ID", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/api/products/invalid/categories", nil)
		req.SetPathValue("id", "invalid")
		w := httptest.NewRecorder()

		handler.GetProductCategories(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "GetProductCategories")
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockCategoryService)
		handler := NewCategoryHandler(mockService)

		productID := uuid.New()

		mockService.On("GetProductCategories", mock.Anything, productID).Return([]*entity.Category{}, errors.New("product not found"))

		req := httptest.NewRequest(http.MethodGet, "/api/products/"+productID.String()+"/categories", nil)
		req.SetPathValue("id", productID.String())
		w := httptest.NewRecorder()

		handler.GetProductCategories(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}
