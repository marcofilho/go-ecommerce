# API Pagination Standardization

## Overview

All list endpoints in the API follow a consistent pagination format with standardized query parameters and response structure.

## Query Parameters

All list endpoints accept the following query parameters:

| Parameter    | Type   | Default        | Description                                      |
|-------------|--------|----------------|--------------------------------------------------|
| `page`      | int    | 1              | Page number (1-indexed)                          |
| `page_size` | int    | 10             | Number of items per page                         |
| `sort_by`   | string | varies         | Field name to sort by (endpoint-specific)        |
| `sort_order`| string | varies         | Sort direction: `asc` or `desc`                  |

### Default Sort Settings by Endpoint

- **Products**: `sort_by=created_at`, `sort_order=desc`
- **Categories**: `sort_by=name`, `sort_order=asc`
- **Orders**: `sort_by=created_at`, `sort_order=desc`
- **Product Variants**: `sort_by=created_at`, `sort_order=asc`

## Response Format

All list endpoints return data in the following standardized format:

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 120,
    "total_pages": 6
  }
}
```

### Response Fields

| Field                 | Type  | Description                                      |
|----------------------|-------|--------------------------------------------------|
| `data`               | array | Array of items for the current page              |
| `pagination.page`    | int   | Current page number                              |
| `pagination.page_size`| int  | Number of items per page                         |
| `pagination.total`   | int   | Total number of items across all pages           |
| `pagination.total_pages`| int| Total number of pages                            |

## Examples

### Categories List

**Request:**
```bash
GET /api/categories?page=1&page_size=2
```

**Response:**
```json
{
  "data": [
    {
      "id": "b2222222-2222-2222-2222-222222222222",
      "name": "Computers"
    },
    {
      "id": "a1111111-1111-1111-1111-111111111111",
      "name": "Electronics"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 2,
    "total": 3,
    "total_pages": 2
  }
}
```

### Products List

**Request:**
```bash
GET /api/products?page=1&page_size=2&sort_by=price&sort_order=asc
```

**Response:**
```json
{
  "data": [
    {
      "id": "e5555555-5555-5555-5555-555555555555",
      "name": "Logitech MX Master 3S Mouse",
      "description": "Advanced wireless mouse with 8K DPI sensor",
      "price": 99.99,
      "quantity": 150,
      "categories": [...],
      "variants": [...],
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "0001-01-01T00:00:00Z"
    },
    {
      "id": "f6666666-6666-6666-6666-666666666666",
      "name": "Sony WH-1000XM5 Headphones",
      "description": "Industry-leading noise canceling headphones",
      "price": 449.99,
      "quantity": 75,
      "categories": [...],
      "variants": [],
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "0001-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 2,
    "total": 3,
    "total_pages": 2
  }
}
```

### Orders List

**Request:**
```bash
GET /api/orders?page=1&page_size=5
```

**Response:**
```json
{
  "data": [
    {
      "id": "10101010-0101-0101-0101-010101010101",
      "customer_id": 100001,
      "products": [...],
      "status": "pending",
      "total": 3598.99,
      "created_at": "2024-12-01T10:00:00Z",
      "updated_at": "2024-12-01T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 5,
    "total": 3,
    "total_pages": 1
  }
}
```

### Product Variants List

**Request:**
```bash
GET /api/products/{id}/variants?page=1&page_size=10
```

**Response:**
```json
{
  "data": [
    {
      "id": "77777777-7777-7777-7777-777777777777",
      "product_id": "d4444444-4444-4444-4444-444444444444",
      "variant_name": "Color",
      "variant_value": "Space Black",
      "price": 0,
      "has_override": false,
      "quantity": 15,
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "0001-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 2,
    "total_pages": 1
  }
}
```

## Implementation Details

### Total Pages Calculation

The `total_pages` field is calculated using the following formula:

```go
totalPages := (total + pageSize - 1) / pageSize
```

This ensures proper rounding up for edge cases:
- If `total = 0`, then `total_pages = 0`
- If `total = 1` and `page_size = 10`, then `total_pages = 1`
- If `total = 11` and `page_size = 10`, then `total_pages = 2`

### Type Definitions

```go
// Pagination metadata structure
type Pagination struct {
    Page       int `json:"page"`
    PageSize   int `json:"page_size"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}

// Generic paginated response
type PaginatedResponse[T any] struct {
    Data       []T        `json:"data"`
    Pagination Pagination `json:"pagination"`
}

// Type aliases for backward compatibility
type ProductListResponse = PaginatedResponse[ProductResponse]
type OrderListResponse = PaginatedResponse[OrderResponse]
type CategoryListResponse = PaginatedResponse[CategoryResponse]
type ProductVariantListResponse = PaginatedResponse[ProductVariantResponse]
```

## Benefits

1. **Consistency**: All list endpoints follow the same pattern
2. **Client-friendly**: `total_pages` allows easy pagination UI implementation
3. **Flexible**: Supports sorting with `sort_by` and `sort_order` parameters
4. **Type-safe**: Generic `PaginatedResponse[T]` ensures type safety
5. **Backward compatible**: Type aliases maintain existing endpoint signatures
