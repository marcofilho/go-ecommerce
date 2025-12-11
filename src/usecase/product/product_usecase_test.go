package product

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
	mockServices "github.com/marcofilho/go-ecommerce/src/internal/testing"
)

type mockProductRepository struct {
	products     map[uuid.UUID]*entity.Product
	createErr    error
	updateErr    error
	deleteErr    error
	getByIDErr   error
	getAllErr    error
	getAllResult []*entity.Product
	getAllTotal  int
}

func newMockRepo() *mockProductRepository {
	return &mockProductRepository{
		products: make(map[uuid.UUID]*entity.Product),
	}
}

func (m *mockProductRepository) Create(ctx context.Context, product *entity.Product) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	p, ok := m.products[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}

func (m *mockProductRepository) GetAll(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error) {
	if m.getAllErr != nil {
		return nil, 0, m.getAllErr
	}
	if m.getAllResult != nil {
		return m.getAllResult, m.getAllTotal, nil
	}
	var result []*entity.Product
	for _, p := range m.products {
		if !inStockOnly || p.Quantity > 0 {
			result = append(result, p)
		}
	}
	return result, len(result), nil
}

func (m *mockProductRepository) Update(ctx context.Context, product *entity.Product) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.products[product.ID]; !ok {
		return errors.New("not found")
	}
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.products[id]; !ok {
		return errors.New("not found")
	}
	delete(m.products, id)
	return nil
}

func TestCreateProduct_Success(t *testing.T) {
	repo := newMockRepo()
	uc := NewUseCase(repo, &mockServices.MockServices{})

	product, err := uc.CreateProduct(context.Background(), "Laptop", "Gaming", 999.99, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if product.Name != "Laptop" {
		t.Errorf("expected name Laptop, got %s", product.Name)
	}
}

func TestCreateProduct_ValidationError(t *testing.T) {
	repo := newMockRepo()
	uc := NewUseCase(repo, &mockServices.MockServices{})

	_, err := uc.CreateProduct(context.Background(), "", "Desc", 100, 10)
	if err == nil {
		t.Error("expected validation error for empty name")
	}
}

func TestGetProduct_Success(t *testing.T) {
	repo := newMockRepo()
	uc := NewUseCase(repo, &mockServices.MockServices{})

	id := uuid.New()
	repo.products[id] = &entity.Product{ID: id, Name: "Test"}

	product, err := uc.GetProduct(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if product.ID != id {
		t.Error("product ID mismatch")
	}
}

func TestListProducts_Success(t *testing.T) {
	repo := newMockRepo()
	uc := NewUseCase(repo, &mockServices.MockServices{})

	repo.getAllResult = []*entity.Product{
		{ID: uuid.New(), Name: "P1", Quantity: 5},
		{ID: uuid.New(), Name: "P2", Quantity: 10},
	}
	repo.getAllTotal = 2

	products, total, err := uc.ListProducts(context.Background(), 1, 10, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

func TestUpdateProduct_Success(t *testing.T) {
	repo := newMockRepo()
	uc := NewUseCase(repo, &mockServices.MockServices{})

	id := uuid.New()
	repo.products[id] = &entity.Product{ID: id, Name: "Old", Price: 100, Quantity: 5}

	updated, err := uc.UpdateProduct(context.Background(), id, "New", "Updated", 200, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Name != "New" {
		t.Errorf("expected name New, got %s", updated.Name)
	}
}

func TestDeleteProduct_Success(t *testing.T) {
	repo := newMockRepo()
	uc := NewUseCase(repo, &mockServices.MockServices{})

	id := uuid.New()
	repo.products[id] = &entity.Product{ID: id}

	err := uc.DeleteProduct(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, ok := repo.products[id]; ok {
		t.Error("product should be deleted")
	}
}

func TestCreateProduct_RepositoryError(t *testing.T) {
	repo := newMockRepo()
	repo.createErr = errors.New("database error")
	uc := NewUseCase(repo, &mockServices.MockServices{})

	_, err := uc.CreateProduct(context.Background(), "Laptop", "Gaming", 999.99, 10)
	if err == nil {
		t.Error("expected error from repository")
	}
}

func TestCreateProduct_ZeroQuantityError(t *testing.T) {
	repo := newMockRepo()
	uc := NewUseCase(repo, &mockServices.MockServices{})

	_, err := uc.CreateProduct(context.Background(), "Laptop", "Gaming", 999.99, 0)
	if err == nil {
		t.Error("expected validation error for zero quantity")
	}
}

func TestListProducts_PaginationDefaults(t *testing.T) {
	repo := newMockRepo()
	uc := NewUseCase(repo, &mockServices.MockServices{})

	// Test page < 1 defaults to 1
	_, _, err := uc.ListProducts(context.Background(), 0, 10, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test page_size < 1 defaults to 10
	_, _, err = uc.ListProducts(context.Background(), 1, 0, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test page_size > 100 defaults to 10
	_, _, err = uc.ListProducts(context.Background(), 1, 150, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateProduct_NotFound(t *testing.T) {
	repo := newMockRepo()
	uc := NewUseCase(repo, &mockServices.MockServices{})

	id := uuid.New()
	_, err := uc.UpdateProduct(context.Background(), id, "New", "Updated", 200, 10)
	if err == nil {
		t.Error("expected not found error")
	}
}

func TestUpdateProduct_ValidationError(t *testing.T) {
	repo := newMockRepo()
	uc := NewUseCase(repo, &mockServices.MockServices{})

	id := uuid.New()
	repo.products[id] = &entity.Product{ID: id, Name: "Old", Price: 100, Quantity: 5}

	_, err := uc.UpdateProduct(context.Background(), id, "", "Updated", 200, 10)
	if err == nil {
		t.Error("expected validation error for empty name")
	}
}

func TestUpdateProduct_RepositoryError(t *testing.T) {
	repo := newMockRepo()
	repo.updateErr = errors.New("database error")
	uc := NewUseCase(repo, &mockServices.MockServices{})

	id := uuid.New()
	repo.products[id] = &entity.Product{ID: id, Name: "Old", Price: 100, Quantity: 5}

	_, err := uc.UpdateProduct(context.Background(), id, "New", "Updated", 200, 10)
	if err == nil {
		t.Error("expected repository error")
	}
}

var _ repository.ProductRepository = (*mockProductRepository)(nil)
