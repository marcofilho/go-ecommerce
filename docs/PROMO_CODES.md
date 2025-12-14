# Promo Codes & Discounts

## Overview

The e-commerce API supports promotional discount codes that can be applied to orders at checkout. Discounts can be either **percentage-based** (e.g., 20% off) or **fixed amount** (e.g., $15 off).

## Features

- ✅ **Two discount types**: Percentage (0-100%) or fixed amount
- ✅ **Active/Inactive control**: Enable or disable promo codes without deletion
- ✅ **Validation endpoint**: Public endpoint to check if a promo code is valid
- ✅ **Applied at checkout**: Discount is calculated when creating an order
- ✅ **No storage**: Discount details are not stored with the order (only final price)
- ✅ **Case-sensitive codes**: Promo codes must match exactly
- ✅ **Admin-only management**: Only admins can create, update, or view discounts

## Database Schema

### Discounts Table

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Unique identifier |
| promo_code | VARCHAR(255) | The promotional code (e.g., "SAVE20") |
| discount_type | VARCHAR(50) | Either "percentage" or "amount" |
| value | DECIMAL(10,2) | The discount value (0-100 for percentage, dollar amount for fixed) |
| active | BOOLEAN | Whether the discount is currently active |
| created_at | TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | Last update timestamp |

**Business Rules:**
- Promo codes must be unique
- Percentage discounts: value must be 0-100
- Amount discounts: value is the dollar amount to subtract
- Only active discounts can be applied to orders
- Discounts are queried with `WHERE active = true` automatically

## API Endpoints

### Admin Endpoints (Require Admin Role)

#### Create Discount
```http
POST /api/discounts
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "promo_code": "SAVE20",
  "discount_type": "percentage",
  "value": 20.0,
  "active": true
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "promo_code": "SAVE20",
  "discount_type": "percentage",
  "value": 20.0,
  "active": true,
  "created_at": "2025-12-14T10:00:00Z",
  "updated_at": "2025-12-14T10:00:00Z"
}
```

#### Get Discount by ID
```http
GET /api/discounts/{id}
Authorization: Bearer {admin_token}
```

#### Update Discount
```http
PUT /api/discounts/{id}
Authorization: Bearer {admin_token}
Content-Type: application/json

{
  "promo_code": "SAVE20",
  "discount_type": "percentage",
  "value": 25.0,
  "active": true
}
```

### Public Endpoints

#### Validate Promo Code
```http
POST /api/discounts/validate
Content-Type: application/json

{
  "promo_code": "SAVE20"
}
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "promo_code": "SAVE20",
  "discount_type": "percentage",
  "value": 20.0,
  "active": true,
  "created_at": "2025-12-14T10:00:00Z",
  "updated_at": "2025-12-14T10:00:00Z"
}
```

**Response (404 Not Found):**
```json
{
  "error": "Promo code not found"
}
```

## Usage in Orders

When creating an order, include the optional `promo_code` field:

```http
POST /api/orders
Authorization: Bearer {token}
Content-Type: application/json

{
  "customer_id": 123,
  "products": [
    {
      "product_id": "660e8400-e29b-41d4-a716-446655440001",
      "quantity": 2
    }
  ],
  "promo_code": "SAVE20"
}
```

### Discount Calculation

**Percentage Discount:**
```
Original Total: $100.00
Discount: 20%
Final Price: $100.00 * (1 - 0.20) = $80.00
```

**Fixed Amount Discount:**
```
Original Total: $100.00
Discount: $15.00
Final Price: $100.00 - $15.00 = $85.00
```

**Note:** If the discount amount exceeds the order total, the final price will be $0.00 (not negative).

## Testing

### Unit Tests

All discount functionality is covered by unit tests:

```bash
# Run all tests
go test ./... -short

# Run discount-specific tests
go test ./src/internal/domain/entity -run="Discount" -v
go test ./src/usecase/discount -v
go test ./src/internal/adapter/http/handler -run="Discount" -v
```

**Test Coverage:**
- **Entity tests**: 5 test functions (validation, discount calculation)
- **Use case tests**: 11 test functions (CRUD operations, validation)
- **Handler tests**: 8 test functions (HTTP endpoints, error handling)
- **Order integration**: 3 test functions (orders with promo codes)

### Integration Tests

Comprehensive end-to-end testing of the promo code feature:

```bash
# Run promo code integration tests
make test-promo
```

**Test Scenarios (10 total):**

1. ✅ Create percentage discount (20% off)
2. ✅ Create amount discount ($15 off)
3. ✅ Order without promo code (baseline $100)
4. ✅ Order with percentage promo ($100 → $80)
5. ✅ Order with amount promo ($100 → $85)
6. ✅ Invalid promo code rejection
7. ✅ Inactive promo code rejection
8. ✅ Validate active promo endpoint
9. ✅ Validate inactive promo endpoint
10. ✅ Multiple items with promo ($300 → $285)

## Error Handling

### Common Errors

**400 Bad Request:**
```json
{
  "error": "Invalid discount type. Must be 'percentage' or 'amount'"
}
```

**400 Bad Request:**
```json
{
  "error": "Percentage discount must be between 0 and 100"
}
```

**404 Not Found:**
```json
{
  "error": "Promo code not found"
}
```

**403 Forbidden:**
```json
{
  "error": "Admin access required"
}
```

## Examples

### Example 1: Create and Use Percentage Discount

```bash
# 1. Login as admin
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@ecommerce.com","password":"password123"}' \
  | jq -r '.token')

# 2. Create 20% off promo code
curl -X POST http://localhost:8080/api/discounts \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "SAVE20",
    "discount_type": "percentage",
    "value": 20.0,
    "active": true
  }'

# 3. Login as customer
CUSTOMER_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john.doe@example.com","password":"password123"}' \
  | jq -r '.token')

# 4. Create order with promo code
curl -X POST http://localhost:8080/api/orders \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 2,
    "products": [{"product_id": "PRODUCT_UUID", "quantity": 2}],
    "promo_code": "SAVE20"
  }'
```

### Example 2: Create and Use Fixed Amount Discount

```bash
# Create $15 off promo code
curl -X POST http://localhost:8080/api/discounts \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "SAVE15",
    "discount_type": "amount",
    "value": 15.0,
    "active": true
  }'

# Use in order
curl -X POST http://localhost:8080/api/orders \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 2,
    "products": [{"product_id": "PRODUCT_UUID", "quantity": 1}],
    "promo_code": "SAVE15"
  }'
```

### Example 3: Validate Promo Code Before Order

```bash
# Validate promo code (public endpoint, no auth required)
curl -X POST http://localhost:8080/api/discounts/validate \
  -H "Content-Type: application/json" \
  -d '{"promo_code": "SAVE20"}'
```

### Example 4: Deactivate a Promo Code

```bash
# Get discount ID first
DISCOUNT_ID="550e8400-e29b-41d4-a716-446655440000"

# Update to inactive
curl -X PUT http://localhost:8080/api/discounts/$DISCOUNT_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "SAVE20",
    "discount_type": "percentage",
    "value": 20.0,
    "active": false
  }'
```

## Architecture

The promo code feature follows clean architecture principles:

### Domain Layer
- **Entity**: `src/internal/domain/entity/discount.go`
  - Validation logic
  - Discount calculation (ApplyDiscount method)
  - UUID generation hook

### Repository Layer
- **Interface**: `src/internal/domain/repository/discount_repository.go`
- **Implementation**: `src/infrastructure/repository/discount_repository_postgres.go`
  - GetByPromoCode filters for `active = true` automatically

### Use Case Layer
- **Interface**: `src/usecase/discount/discount_usecase.go`
- Methods: CreateDiscount, GetDiscountByID, UpdateDiscount, GetDiscountByPromoCode

### Handler Layer
- **Handler**: `src/internal/adapter/http/handler/discount_handler.go`
- Endpoints: Create, Get, Update, Validate
- Validation: Checks discount_type and value ranges

### Order Integration
- **Order Use Case**: `src/usecase/order/order_usecase.go`
- Optional promo_code parameter in CreateOrder
- Applies discount before saving final order price

## Security

- ✅ **Admin-only management**: Only admins can create, view, or modify discounts
- ✅ **Public validation**: Anyone can validate a promo code (for UI feedback)
- ✅ **Active-only queries**: Database queries automatically filter for active discounts
- ✅ **No injection risks**: UUID-based IDs, parameterized queries via GORM

## Future Enhancements

Potential improvements for the promo code system:

- [ ] Usage limits (max uses per code)
- [ ] Expiration dates
- [ ] User-specific promo codes
- [ ] Minimum order amount requirements
- [ ] Category-specific discounts
- [ ] First-time customer discounts
- [ ] Stackable discounts
- [ ] Discount usage analytics

## References

- [Main README](../README.md)
- [API Documentation](http://localhost:8080/swagger/index.html)
- [Testing Guide](TESTING.md)
- [Database Schema](DATABASE_SCHEMA.md)
