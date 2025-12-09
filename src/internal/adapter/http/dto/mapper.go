package dto

import (
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
)

// Product Mappers
func ToProductResponse(product *entity.Product) ProductResponse {
	// Map categories
	categories := make([]CategoryResponse, 0, len(product.Categories))
	for _, cat := range product.Categories {
		categories = append(categories, CategoryResponse{
			ID:   cat.ID.String(),
			Name: cat.Name,
		})
	}

	// Map variants
	variants := make([]ProductVariantResponse, 0, len(product.Variants))
	for _, variant := range product.Variants {
		variants = append(variants, ToProductVariantResponse(&variant))
	}

	return ProductResponse{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		Categories:  categories,
		Variants:    variants,
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

// ProductVariant Mappers
func ToProductVariantResponse(variant *entity.ProductVariant) ProductVariantResponse {
	price, _ := variant.GetPrice() // Ignoring error for response mapping

	return ProductVariantResponse{
		ID:            variant.ID.String(),
		ProductID:     variant.ProductID.String(),
		VariantName:   variant.VariantName,
		VariantValue:  variant.VariantValue,
		Price:         price,
		PriceOverride: variant.Price_Override,
		HasOverride:   variant.HasPriceOverride(),
		Quantity:      variant.Quantity,
		CreatedAt:     variant.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     variant.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ToProductVariantListResponse(variants []*entity.ProductVariant, total, page, pageSize int) PaginatedResponse[ProductVariantResponse] {
	variantResponses := make([]ProductVariantResponse, 0, len(variants))
	for _, variant := range variants {
		variantResponses = append(variantResponses, ToProductVariantResponse(variant))
	}

	return PaginatedResponse[ProductVariantResponse]{
		Data:     variantResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
