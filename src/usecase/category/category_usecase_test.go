package category

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

// MockCategoryRepository is a mock implementation of repository.CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetAll(ctx context.Context, page, pageSize int) ([]*entity.Category, int, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*entity.Category), args.Get(1).(int), args.Error(2)
}

func (m *MockCategoryRepository) GetByName(ctx context.Context, name string) (*entity.Category, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryRepository) AssignCategoryToProduct(ctx context.Context, productID, categoryID uuid.UUID) error {
	args := m.Called(ctx, productID, categoryID)
	return args.Error(0)
}

func (m *MockCategoryRepository) RemoveCategoryFromProduct(ctx context.Context, productID, categoryID uuid.UUID) error {
	args := m.Called(ctx, productID, categoryID)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetProductCategories(ctx context.Context, productID uuid.UUID) ([]*entity.Category, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).([]*entity.Category), args.Error(1)
}

func TestUseCase_CreateCategory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		name := "Electronics"

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(c *entity.Category) bool {
			return c.Name == name
		})).Return(nil)

		result, err := useCase.CreateCategory(context.Background(), name)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, name, result.Name)
		assert.NotEqual(t, uuid.Nil, result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation Error - Empty Name", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		result, err := useCase.CreateCategory(context.Background(), "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Category name is required")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		name := "Electronics"

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(c *entity.Category) bool {
			return c.Name == name
		})).Return(errors.New("database error"))

		result, err := useCase.CreateCategory(context.Background(), name)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestUseCase_GetCategory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		categoryID := uuid.New()
		expectedCategory := &entity.Category{
			ID:   categoryID,
			Name: "Electronics",
		}

		mockRepo.On("GetByID", mock.Anything, categoryID).Return(expectedCategory, nil)

		result, err := useCase.GetCategory(context.Background(), categoryID)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		categoryID := uuid.New()

		mockRepo.On("GetByID", mock.Anything, categoryID).Return(nil, errors.New("category not found"))

		result, err := useCase.GetCategory(context.Background(), categoryID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUseCase_ListCategories(t *testing.T) {
	t.Run("Success - Default Pagination", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		expectedCategories := []*entity.Category{
			{ID: uuid.New(), Name: "Electronics"},
			{ID: uuid.New(), Name: "Clothing"},
		}
		expectedTotal := 2

		mockRepo.On("GetAll", mock.Anything, 1, 10).Return(expectedCategories, expectedTotal, nil)

		categories, total, err := useCase.ListCategories(context.Background(), 0, 0)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
		assert.Equal(t, expectedTotal, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Custom Pagination", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		expectedCategories := []*entity.Category{
			{ID: uuid.New(), Name: "Electronics"},
		}
		expectedTotal := 10

		mockRepo.On("GetAll", mock.Anything, 2, 5).Return(expectedCategories, expectedTotal, nil)

		categories, total, err := useCase.ListCategories(context.Background(), 2, 5)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
		assert.Equal(t, expectedTotal, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - Max Page Size Limit", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		expectedCategories := []*entity.Category{}
		expectedTotal := 0

		// Should limit to 100
		mockRepo.On("GetAll", mock.Anything, 1, 10).Return(expectedCategories, expectedTotal, nil)

		categories, total, err := useCase.ListCategories(context.Background(), 1, 200)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
		assert.Equal(t, expectedTotal, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		mockRepo.On("GetAll", mock.Anything, 1, 10).Return([]*entity.Category{}, 0, errors.New("database error"))

		categories, total, err := useCase.ListCategories(context.Background(), 1, 10)

		assert.Error(t, err)
		assert.Empty(t, categories)
		assert.Equal(t, 0, total)
		mockRepo.AssertExpectations(t)
	})
}

func TestUseCase_UpdateCategory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		categoryID := uuid.New()
		existingCategory := &entity.Category{
			ID:   categoryID,
			Name: "Old Name",
		}
		newName := "Updated Electronics"

		mockRepo.On("GetByID", mock.Anything, categoryID).Return(existingCategory, nil)
		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *entity.Category) bool {
			return c.ID == categoryID && c.Name == newName
		})).Return(nil)

		result, err := useCase.UpdateCategory(context.Background(), categoryID, newName)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, newName, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		categoryID := uuid.New()
		existingCategory := &entity.Category{
			ID:   categoryID,
			Name: "Old Name",
		}

		mockRepo.On("GetByID", mock.Anything, categoryID).Return(existingCategory, nil)

		result, err := useCase.UpdateCategory(context.Background(), categoryID, "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Category name is required")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("Category Not Found", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		categoryID := uuid.New()

		mockRepo.On("GetByID", mock.Anything, categoryID).Return(nil, errors.New("not found"))

		result, err := useCase.UpdateCategory(context.Background(), categoryID, "New Name")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		categoryID := uuid.New()
		existingCategory := &entity.Category{
			ID:   categoryID,
			Name: "Old Name",
		}

		mockRepo.On("GetByID", mock.Anything, categoryID).Return(existingCategory, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))

		result, err := useCase.UpdateCategory(context.Background(), categoryID, "New Name")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUseCase_DeleteCategory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		categoryID := uuid.New()

		mockRepo.On("Delete", mock.Anything, categoryID).Return(nil)

		err := useCase.DeleteCategory(context.Background(), categoryID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		categoryID := uuid.New()

		mockRepo.On("Delete", mock.Anything, categoryID).Return(errors.New("database error"))

		err := useCase.DeleteCategory(context.Background(), categoryID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUseCase_AssignCategoryToProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		productID := uuid.New()
		categoryID := uuid.New()

		mockRepo.On("AssignCategoryToProduct", mock.Anything, productID, categoryID).Return(nil)

		err := useCase.AssignCategoryToProduct(context.Background(), productID, categoryID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		productID := uuid.New()
		categoryID := uuid.New()

		mockRepo.On("AssignCategoryToProduct", mock.Anything, productID, categoryID).Return(errors.New("already assigned"))

		err := useCase.AssignCategoryToProduct(context.Background(), productID, categoryID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUseCase_RemoveCategoryFromProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		productID := uuid.New()
		categoryID := uuid.New()

		mockRepo.On("RemoveCategoryFromProduct", mock.Anything, productID, categoryID).Return(nil)

		err := useCase.RemoveCategoryFromProduct(context.Background(), productID, categoryID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		productID := uuid.New()
		categoryID := uuid.New()

		mockRepo.On("RemoveCategoryFromProduct", mock.Anything, productID, categoryID).Return(errors.New("not found"))

		err := useCase.RemoveCategoryFromProduct(context.Background(), productID, categoryID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUseCase_GetProductCategories(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		productID := uuid.New()
		expectedCategories := []*entity.Category{
			{ID: uuid.New(), Name: "Electronics"},
			{ID: uuid.New(), Name: "Computers"},
		}

		mockRepo.On("GetProductCategories", mock.Anything, productID).Return(expectedCategories, nil)

		categories, err := useCase.GetProductCategories(context.Background(), productID)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		useCase := NewUseCase(mockRepo)

		productID := uuid.New()

		mockRepo.On("GetProductCategories", mock.Anything, productID).Return([]*entity.Category{}, errors.New("product not found"))

		categories, err := useCase.GetProductCategories(context.Background(), productID)

		assert.Error(t, err)
		assert.Empty(t, categories)
		mockRepo.AssertExpectations(t)
	})
}
