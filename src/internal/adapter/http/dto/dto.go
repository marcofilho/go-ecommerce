package dto

// Generic paginated response
type PaginatedResponse[T any] struct {
	Data     []T `json:"data"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// Product DTOs
type ProductRequest struct {
	Name        string  `json:"name" example:"Laptop"`
	Description string  `json:"description" example:"High-performance laptop"`
	Price       float64 `json:"price" example:"999.99"`
	Quantity    int     `json:"quantity" example:"50"`
}

type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Order DTOs
type CreateOrderRequest struct {
	CustomerID int                `json:"customer_id" example:"123"`
	Products   []OrderItemRequest `json:"products"`
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

// Auth DTOs
type AuthResponse struct {
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	ExpiresAt string `json:"expires_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ProductListResponse struct {
	Data     []ProductResponse `json:"data"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

type OrderListResponse struct {
	Data     []OrderResponse `json:"data"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

type ProductVariantListResponse struct {
	Data     []ProductVariantResponse `json:"data"`
	Total    int                      `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
}
