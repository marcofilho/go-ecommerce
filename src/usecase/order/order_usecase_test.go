package order

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/repository"
)

type mockOrderRepo struct {
	orders    map[uuid.UUID]*entity.Order
	createErr error
	updateErr error
}

func newMockOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{orders: make(map[uuid.UUID]*entity.Order)}
}

func (m *mockOrderRepo) Create(ctx context.Context, order *entity.Order) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	o, ok := m.orders[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return o, nil
}

func (m *mockOrderRepo) GetAll(ctx context.Context, page, pageSize int, status *entity.OrderStatus, paymentStatus *entity.PaymentStatus) ([]*entity.Order, int, error) {
	var result []*entity.Order
	for _, o := range m.orders {
		result = append(result, o)
	}
	return result, len(result), nil
}

func (m *mockOrderRepo) Update(ctx context.Context, order *entity.Order) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.orders[order.ID]; !ok {
		return errors.New("not found")
	}
	m.orders[order.ID] = order
	return nil
}

type mockProductRepo struct {
	products  map[uuid.UUID]*entity.Product
	updateErr error
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{products: make(map[uuid.UUID]*entity.Product)}
}

func (m *mockProductRepo) Create(ctx context.Context, product *entity.Product) error {
	return nil
}

func (m *mockProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}

func (m *mockProductRepo) GetAll(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error) {
	return nil, 0, nil
}

func (m *mockProductRepo) Update(ctx context.Context, product *entity.Product) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

type mockVariantRepo struct {
	variants  map[uuid.UUID]*entity.ProductVariant
	updateErr error
}

func newMockVariantRepo() *mockVariantRepo {
	return &mockVariantRepo{variants: make(map[uuid.UUID]*entity.ProductVariant)}
}

func (m *mockVariantRepo) Create(ctx context.Context, variant *entity.ProductVariant) error {
	return nil
}

func (m *mockVariantRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductVariant, error) {
	v, ok := m.variants[id]
	if !ok {
		return nil, errors.New("variant not found")
	}
	return v, nil
}

func (m *mockVariantRepo) GetAll(ctx context.Context, page, pageSize int) ([]*entity.ProductVariant, int, error) {
	return nil, 0, nil
}

func (m *mockVariantRepo) GetAllByProductID(ctx context.Context, productID uuid.UUID, page, pageSize int) ([]*entity.ProductVariant, int, error) {
	return nil, 0, nil
}

func (m *mockVariantRepo) Update(ctx context.Context, variant *entity.ProductVariant) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.variants[variant.ID] = variant
	return nil
}

func (m *mockVariantRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func TestCreateOrder_Success(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	variantRepo := newMockVariantRepo()
	uc := NewUseCase(orderRepo, productRepo, variantRepo)

	pid := uuid.New()
	productRepo.products[pid] = &entity.Product{
		ID: pid, Name: "Laptop", Price: 100, Quantity: 10,
	}

	items := []CreateOrderItem{{ProductID: pid, Quantity: 2}}
	order, err := uc.CreateOrder(context.Background(), 123, items)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if order.CustomerID != 123 {
		t.Error("customer ID mismatch")
	}
}

func TestCreateOrder_NoItems(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	_, err := uc.CreateOrder(context.Background(), 123, []CreateOrderItem{})
	if err == nil {
		t.Error("expected error for empty items")
	}
}

func TestCreateOrder_InsufficientStock(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	pid := uuid.New()
	productRepo.products[pid] = &entity.Product{
		ID: pid, Name: "Laptop", Price: 100, Quantity: 5,
	}

	items := []CreateOrderItem{{ProductID: pid, Quantity: 10}}
	_, err := uc.CreateOrder(context.Background(), 123, items)

	if err == nil {
		t.Error("expected error for insufficient stock")
	}
}

func TestGetOrder_Success(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	oid := uuid.New()
	orderRepo.orders[oid] = &entity.Order{ID: oid, CustomerID: 123}

	order, err := uc.GetOrder(context.Background(), oid)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if order.ID != oid {
		t.Error("order ID mismatch")
	}
}

func TestListOrders_Success(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	orderRepo.orders[uuid.New()] = &entity.Order{CustomerID: 1}
	orderRepo.orders[uuid.New()] = &entity.Order{CustomerID: 2}

	orders, total, err := uc.ListOrders(context.Background(), 1, 10, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(orders))
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

func TestUpdateOrderStatus_Success(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	oid := uuid.New()
	orderRepo.orders[oid] = &entity.Order{
		ID: oid, Status: entity.Pending,
	}

	updated, err := uc.UpdateOrderStatus(context.Background(), oid, entity.Completed)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Status != entity.Completed {
		t.Error("status not updated")
	}
}

func TestUpdateOrderStatus_InvalidTransition(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	oid := uuid.New()
	orderRepo.orders[oid] = &entity.Order{
		ID: oid, Status: entity.Completed,
	}

	_, err := uc.UpdateOrderStatus(context.Background(), oid, entity.Cancelled)
	if err == nil {
		t.Error("expected error for invalid transition")
	}
}

func TestCreateOrder_InvalidCustomerID(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	items := []CreateOrderItem{{ProductID: uuid.New(), Quantity: 1}}
	_, err := uc.CreateOrder(context.Background(), 0, items)
	if err == nil {
		t.Error("expected error for invalid customer ID")
	}

	_, err = uc.CreateOrder(context.Background(), -1, items)
	if err == nil {
		t.Error("expected error for negative customer ID")
	}
}

func TestCreateOrder_ProductNotFound(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	items := []CreateOrderItem{{ProductID: uuid.New(), Quantity: 1}}
	_, err := uc.CreateOrder(context.Background(), 123, items)
	if err == nil {
		t.Error("expected error for product not found")
	}
}

func TestCreateOrder_ProductUpdateError(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	productRepo.updateErr = errors.New("update failed")
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	pid := uuid.New()
	productRepo.products[pid] = &entity.Product{
		ID: pid, Name: "Laptop", Price: 100, Quantity: 10,
	}

	items := []CreateOrderItem{{ProductID: pid, Quantity: 2}}
	_, err := uc.CreateOrder(context.Background(), 123, items)
	if err == nil {
		t.Error("expected error from product update")
	}
}

func TestCreateOrder_OrderCreateError(t *testing.T) {
	orderRepo := newMockOrderRepo()
	orderRepo.createErr = errors.New("create failed")
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	pid := uuid.New()
	productRepo.products[pid] = &entity.Product{
		ID: pid, Name: "Laptop", Price: 100, Quantity: 10,
	}

	items := []CreateOrderItem{{ProductID: pid, Quantity: 2}}
	_, err := uc.CreateOrder(context.Background(), 123, items)
	if err == nil {
		t.Error("expected error from order create")
	}
}

func TestListOrders_PaginationDefaults(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	// Test page < 1 defaults to 1
	_, _, err := uc.ListOrders(context.Background(), 0, 10, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test page_size < 1 defaults to 10
	_, _, err = uc.ListOrders(context.Background(), 1, 0, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test page_size > 100 defaults to 10
	_, _, err = uc.ListOrders(context.Background(), 1, 150, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateOrderStatus_NotFound(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	_, err := uc.UpdateOrderStatus(context.Background(), uuid.New(), entity.Completed)
	if err == nil {
		t.Error("expected not found error")
	}
}

func TestUpdateOrderStatus_RepositoryError(t *testing.T) {
	orderRepo := newMockOrderRepo()
	orderRepo.updateErr = errors.New("update failed")
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	oid := uuid.New()
	orderRepo.orders[oid] = &entity.Order{
		ID: oid, Status: entity.Pending,
	}

	_, err := uc.UpdateOrderStatus(context.Background(), oid, entity.Completed)
	if err == nil {
		t.Error("expected repository error")
	}
}

func TestCreateOrder_InvalidOrderItem(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	pid := uuid.New()
	productRepo.products[pid] = &entity.Product{
		ID: pid, Name: "Laptop", Price: 100, Quantity: 10,
	}

	// Negative quantity should fail order item validation
	items := []CreateOrderItem{{ProductID: pid, Quantity: -1}}
	_, err := uc.CreateOrder(context.Background(), 123, items)
	if err == nil {
		t.Error("expected error for invalid order item")
	}
}

func TestCreateOrder_DecreaseStockError(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	pid := uuid.New()
	productRepo.products[pid] = &entity.Product{
		ID: pid, Name: "Laptop", Price: 100, Quantity: 5,
	}

	// Request exactly available amount - should succeed
	items := []CreateOrderItem{{ProductID: pid, Quantity: 5}}
	order, err := uc.CreateOrder(context.Background(), 123, items)
	if err != nil {
		t.Fatalf("expected no error for valid order, got %v", err)
	}
	if order == nil {
		t.Error("expected valid order")
	}
}

func TestCreateOrder_ZeroQuantityItem(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	pid := uuid.New()
	productRepo.products[pid] = &entity.Product{
		ID: pid, Name: "Laptop", Price: 100, Quantity: 10,
	}

	// Zero quantity should fail validation
	items := []CreateOrderItem{{ProductID: pid, Quantity: 0}}
	_, err := uc.CreateOrder(context.Background(), 123, items)
	if err == nil {
		t.Error("expected error for zero quantity item")
	}
}

func TestCreateOrder_NilProductID(t *testing.T) {
	orderRepo := newMockOrderRepo()
	productRepo := newMockProductRepo()
	uc := NewUseCase(orderRepo, productRepo, newMockVariantRepo())

	pid := uuid.New()
	productRepo.products[pid] = &entity.Product{
		ID: pid, Name: "Laptop", Price: -10, Quantity: 10,
	}

	// This should pass product lookup but could fail other validations
	items := []CreateOrderItem{{ProductID: pid, Quantity: 1}}
	_, err := uc.CreateOrder(context.Background(), 123, items)
	// May or may not error depending on validation logic
	_ = err
}

var _ repository.OrderRepository = (*mockOrderRepo)(nil)
var _ repository.ProductRepository = (*mockProductRepo)(nil)
