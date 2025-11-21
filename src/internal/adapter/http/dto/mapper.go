package dto

import (
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

// Product Mappers
func ToProductResponse(product *entity.Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ToProductListResponse(products []*entity.Product, total, page, pageSize int) PaginatedResponse[ProductResponse] {
	productResponses := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		productResponses = append(productResponses, ToProductResponse(product))
	}

	return PaginatedResponse[ProductResponse]{
		Data:     productResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

// Order Mappers
func ToOrderResponse(order *entity.Order) OrderResponse {
	products := make([]OrderItemResponse, 0, len(order.Products))
	for _, product := range order.Products {
		products = append(products, OrderItemResponse{
			ProductID: product.ProductID.String(),
			Quantity:  product.Quantity,
			Subtotal:  product.Subtotal(),
		})
	}

	return OrderResponse{
		ID:            order.ID.String(),
		CustomerID:    order.CustomerID,
		Products:      products,
		TotalPrice:    order.TotalPrice,
		Status:        string(order.Status),
		PaymentStatus: string(order.PaymentStatus),
		CreatedAt:     order.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     order.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ToOrderListResponse(orders []*entity.Order, total, page, pageSize int) PaginatedResponse[OrderResponse] {
	orderResponses := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, ToOrderResponse(order))
	}

	return PaginatedResponse[OrderResponse]{
		Data:     orderResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
