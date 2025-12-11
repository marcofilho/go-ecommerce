package dto

import (
	"testing"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

func TestToProductResponse(t *testing.T) {
	id := uuid.New()
	product := &entity.Product{
		ID:          id,
		Name:        "Laptop",
		Description: "High-end gaming laptop",
		Price:       1299.99,
		Quantity:    5,
	}

	response := ToProductResponse(product)

	if response.ID != id.String() {
		t.Errorf("ToProductResponse() ID = %v, want %v", response.ID, id.String())
	}
	if response.Name != "Laptop" {
		t.Errorf("ToProductResponse() Name = %v, want Laptop", response.Name)
	}
	if response.Price != 1299.99 {
		t.Errorf("ToProductResponse() Price = %v, want 1299.99", response.Price)
	}
	if response.Quantity != 5 {
		t.Errorf("ToProductResponse() Quantity = %v, want 5", response.Quantity)
	}
}

func TestToProductListResponse(t *testing.T) {
	products := []*entity.Product{
		{
			ID:       uuid.New(),
			Name:     "Laptop",
			Price:    1299.99,
			Quantity: 5,
		},
		{
			ID:       uuid.New(),
			Name:     "Mouse",
			Price:    29.99,
			Quantity: 100,
		},
	}

	response := ToProductListResponse(products, 2, 1, 10)

	if len(response.Data) != 2 {
		t.Errorf("ToProductListResponse() length = %v, want 2", len(response.Data))
	}
	if response.Pagination.Total != 2 {
		t.Errorf("ToProductListResponse() Total = %v, want 2", response.Pagination.Total)
	}
	if response.Pagination.Page != 1 {
		t.Errorf("ToProductListResponse() Page = %v, want 1", response.Pagination.Page)
	}
	if response.Pagination.PageSize != 10 {
		t.Errorf("ToProductListResponse() PageSize = %v, want 10", response.Pagination.PageSize)
	}
	if response.Pagination.TotalPages != 1 {
		t.Errorf("ToProductListResponse() TotalPages = %v, want 1", response.Pagination.TotalPages)
	}
	if response.Data[0].Name != "Laptop" {
		t.Errorf("ToProductListResponse() Data[0].Name = %v, want Laptop", response.Data[0].Name)
	}
}

func TestToOrderResponse(t *testing.T) {
	orderID := uuid.New()
	productID := uuid.New()
	order := &entity.Order{
		ID:         orderID,
		CustomerID: 123,
		Products: []entity.OrderItem{
			{
				ID:         uuid.New(),
				ProductID:  productID,
				Quantity:   2,
				Price:      100.00,
				TotalPrice: 200.00,
			},
		},
		TotalPrice:    200.00,
		Status:        entity.Pending,
		PaymentStatus: entity.Unpaid,
	}

	response := ToOrderResponse(order)

	if response.ID != orderID.String() {
		t.Errorf("ToOrderResponse() ID = %v, want %v", response.ID, orderID.String())
	}
	if response.CustomerID != 123 {
		t.Errorf("ToOrderResponse() CustomerID = %v, want 123", response.CustomerID)
	}
	if response.TotalPrice != 200.00 {
		t.Errorf("ToOrderResponse() TotalPrice = %v, want 200.00", response.TotalPrice)
	}
	if response.Status != string(entity.Pending) {
		t.Errorf("ToOrderResponse() Status = %v, want pending", response.Status)
	}
	if len(response.Products) != 1 {
		t.Errorf("ToOrderResponse() Products length = %v, want 1", len(response.Products))
	}
	if response.Products[0].ProductID != productID.String() {
		t.Errorf("ToOrderResponse() Products[0].ProductID = %v, want %v", response.Products[0].ProductID, productID.String())
	}
	if response.Products[0].Subtotal != 200.00 {
		t.Errorf("ToOrderResponse() Products[0].Subtotal = %v, want 200.00", response.Products[0].Subtotal)
	}
}

func TestToOrderListResponse(t *testing.T) {
	orders := []*entity.Order{
		{
			ID:         uuid.New(),
			CustomerID: 1,
			TotalPrice: 100.00,
			Status:     entity.Pending,
		},
		{
			ID:         uuid.New(),
			CustomerID: 2,
			TotalPrice: 200.00,
			Status:     entity.Completed,
		},
	}

	response := ToOrderListResponse(orders, 2, 1, 10)

	if len(response.Data) != 2 {
		t.Errorf("ToOrderListResponse() length = %v, want 2", len(response.Data))
	}
	if response.Pagination.Total != 2 {
		t.Errorf("ToOrderListResponse() Total = %v, want 2", response.Pagination.Total)
	}
	if response.Pagination.Page != 1 {
		t.Errorf("ToOrderListResponse() Page = %v, want 1", response.Pagination.Page)
	}
	if response.Pagination.PageSize != 10 {
		t.Errorf("ToOrderListResponse() PageSize = %v, want 10", response.Pagination.PageSize)
	}
	if response.Pagination.TotalPages != 1 {
		t.Errorf("ToOrderListResponse() TotalPages = %v, want 1", response.Pagination.TotalPages)
	}
	if response.Data[0].CustomerID != 1 {
		t.Errorf("ToOrderListResponse() Data[0].CustomerID = %v, want 1", response.Data[0].CustomerID)
	}
}
