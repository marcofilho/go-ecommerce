package dto

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Product DTOs
type ProductRequest struct {
	Name        string  `json:"name" example:"Laptop"`
	Description string  `json:"description" example:"High-performance laptop"`
	Price       float64 `json:"price" example:"999.99"`
	Quantity    int     `json:"quantity" example:"50"`
}

type ProductResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Price       float64                  `json:"price"`
	Quantity    int                      `json:"quantity"`
	Categories  []CategoryResponse       `json:"categories,omitempty"`
	Variants    []ProductVariantResponse `json:"variants,omitempty"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
}

// Order DTOs
type CreateOrderRequest struct {
	CustomerID int                `json:"customer_id" example:"123"`
	Products   []OrderItemRequest `json:"products"`
	PromoCode  string             `json:"promo_code,omitempty" example:"SUMMER2024"` // Optional promo code
}

type OrderItemRequest struct {
	ProductID string  `json:"product_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	VariantID *string `json:"variant_id,omitempty" example:"660e8400-e29b-41d4-a716-446655440000"` // Optional: order specific variant
	Quantity  int     `json:"quantity" example:"2"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" example:"completed"`
}

type OrderItemResponse struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}

type OrderResponse struct {
	ID            string              `json:"id"`
	CustomerID    int                 `json:"customer_id"`
	Products      []OrderItemResponse `json:"products"`
	TotalPrice    float64             `json:"total_price"`
	Status        string              `json:"status"`
	PaymentStatus string              `json:"payment_status"`
	CreatedAt     string              `json:"created_at"`
	UpdatedAt     string              `json:"updated_at"`
}

// ProductVariant DTOs
type ProductVariantRequest struct {
	ProductID     string   `json:"product_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	VariantName   string   `json:"variant_name" example:"Color"`
	VariantValue  string   `json:"variant_value" example:"Red"`
	PriceOverride *float64 `json:"price_override,omitempty" example:"99.99"` // Optional price override
	Quantity      int      `json:"quantity" example:"10"`
}

type ProductVariantResponse struct {
	ID            string   `json:"id"`
	ProductID     string   `json:"product_id"`
	VariantName   string   `json:"variant_name"`
	VariantValue  string   `json:"variant_value"`
	Price         float64  `json:"price"`                    // Effective price (override or base product price)
	PriceOverride *float64 `json:"price_override,omitempty"` // The override value if set
	HasOverride   bool     `json:"has_override"`             // Indicates if price is overridden
	Quantity      int      `json:"quantity"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

// Category DTOs
type CategoryRequest struct {
	Name string `json:"name" example:"Electronics"`
}

type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AssignCategoryRequest struct {
	CategoryID string `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// Auth DTOs
type AuthResponse struct {
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	ExpiresAt string `json:"expires_at"`
}

// Discount DTOs
type DiscountRequest struct {
	PromoCode         string   `json:"promo_code" example:"SUMMER2024"`
	DiscountType      string   `json:"discount_type" example:"percentage"` // "percentage" or "amount"
	Value             float64  `json:"value" example:"15.0"`               // Percentage (0-100) or fixed amount
	Active            bool     `json:"active" example:"true"`
	MinPurchaseAmount *float64 `json:"min_purchase_amount,omitempty" example:"50.00"`
	MaxDiscountAmount *float64 `json:"max_discount_amount,omitempty" example:"100.00"`
	UsageLimit        *int     `json:"usage_limit,omitempty" example:"100"`
	ValidFrom         *string  `json:"valid_from,omitempty" example:"2025-06-01T00:00:00Z"`
	ValidUntil        *string  `json:"valid_until,omitempty" example:"2025-12-31T23:59:59Z"`
	ProductIDs        []string `json:"product_ids,omitempty"`                  // Empty = applies to all products
	CategoryIDs       []string `json:"category_ids,omitempty"`                 // Empty = applies to all categories
	UserIDs           []string `json:"user_ids,omitempty"`                     // Empty = applies to all users
	UserUsageLimit    *int     `json:"user_usage_limit,omitempty" example:"5"` // Per-user usage limit
}

type DiscountResponse struct {
	ID                string             `json:"id"`
	PromoCode         string             `json:"promo_code"`
	DiscountType      string             `json:"discount_type"` // "percentage" or "amount"
	Value             float64            `json:"value"`
	Active            bool               `json:"active"`
	MinPurchaseAmount *float64           `json:"min_purchase_amount,omitempty"`
	MaxDiscountAmount *float64           `json:"max_discount_amount,omitempty"`
	UsageLimit        *int               `json:"usage_limit,omitempty"`
	UsageCount        int                `json:"usage_count"`
	ValidFrom         *string            `json:"valid_from,omitempty"`
	ValidUntil        *string            `json:"valid_until,omitempty"`
	Products          []ProductResponse  `json:"products,omitempty"`
	Categories        []CategoryResponse `json:"categories,omitempty"`
	Users             []UserSummary      `json:"users,omitempty"`
	CreatedAt         string             `json:"created_at"`
	UpdatedAt         string             `json:"updated_at"`
}

type UserSummary struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ApplyDiscountRequest struct {
	PromoCode string `json:"promo_code" example:"SUMMER2024"`
}

type ValidateDiscountRequest struct {
	PromoCode  string                     `json:"promo_code" example:"SUMMER2024"`
	OrderItems []ValidateDiscountItemInfo `json:"order_items"`
}

type ValidateDiscountItemInfo struct {
	ProductID string  `json:"product_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Quantity  int     `json:"quantity" example:"2"`
	Price     float64 `json:"price" example:"99.99"`
}

type ValidateDiscountResponse struct {
	Valid             bool              `json:"valid"`
	Discount          *DiscountResponse `json:"discount,omitempty"`
	ApplicableItems   []string          `json:"applicable_items,omitempty"`
	EstimatedDiscount float64           `json:"estimated_discount"`
	Message           string            `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Type aliases for backward compatibility and cleaner Swagger docs
type ProductListResponse = PaginatedResponse[ProductResponse]
type OrderListResponse = PaginatedResponse[OrderResponse]
type ProductVariantListResponse = PaginatedResponse[ProductVariantResponse]
type CategoryListResponse = PaginatedResponse[CategoryResponse]
