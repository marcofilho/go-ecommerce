# Testing Guide

## Unit Tests

Run all unit tests:
```bash
make test
```

Current test coverage: **123 tests** across 17 packages, all passing.

Key test files:
- `src/internal/domain/entity/*_test.go` - Domain entity tests (User, Product, ProductVariant, Order)
- `src/internal/adapter/http/handler/*_test.go` - HTTP handler tests
- `src/usecase/*_test.go` - Use case business logic tests
- `src/internal/infrastructure/auth/*_test.go` - JWT authentication tests

## Integration Tests

### Prerequisites
```bash
# Start services
docker-compose up -d

# Wait for services to be ready
sleep 3
```

### Basic Integration Test
Tests all major endpoints for correct HTTP status codes and authorization:
```bash
./test_authentication.sh
```

### Full Workflow Test
Comprehensive test including product variants, orders with variants, and authorization:
```bash
./test_product_variants_full.sh
```

### Manual API Testing

#### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

#### 2. Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Save the token from the response.

#### 3. Create a Product (Admin only)
```bash
TOKEN="your_admin_token_here"

curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "T-Shirt",
    "description": "Cotton T-Shirt",
    "price": 29.99,
    "quantity": 100
  }'
```

#### 4. List Products (Public)
```bash
curl http://localhost:8080/api/products
```

#### 5. Create Product Variant (Admin only)
```bash
PRODUCT_ID="product_uuid_here"
TOKEN="your_admin_token_here"

curl -X POST "http://localhost:8080/api/products/$PRODUCT_ID/variants" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Color",
    "value": "Red",
    "quantity": 50,
    "price_override": 34.99
  }'
```

#### 6. List Product Variants (Public)
```bash
PRODUCT_ID="product_uuid_here"

curl "http://localhost:8080/api/products/$PRODUCT_ID/variants"
```

#### 7. Create Order with Variant
```bash
TOKEN="your_token_here"
PRODUCT_ID="product_uuid_here"
VARIANT_ID="variant_uuid_here"

curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"items\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"variant_id\": \"$VARIANT_ID\",
        \"quantity\": 2
      }
    ]
  }"
```

#### 8. Create Order without Variant (Base Product)
```bash
TOKEN="your_token_here"
PRODUCT_ID="product_uuid_here"

curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"items\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 1
      }
    ]
  }"
```

## Swagger UI

Interactive API documentation is available at:
```
http://localhost:8080/swagger/index.html
```

To update Swagger documentation after API changes:
```bash
make swagger
```

## Test Results Summary

### Integration Test Results (All Endpoints)

| Test Category | Endpoint | Status | Notes |
|--------------|----------|--------|-------|
| **Authentication** |
| | POST /api/auth/register | ✅ | Returns 201 with JWT token |
| | POST /api/auth/login | ✅ | Returns 401 for invalid credentials |
| **Products** |
| | POST /api/products | ✅ | Returns 401 without auth, 403 for non-admin |
| | GET /api/products | ✅ | Public access, returns 200 |
| | GET /api/products/{id} | ✅ | Public access |
| | PUT /api/products/{id} | ✅ | Admin only, returns 403 for customer |
| | DELETE /api/products/{id} | ✅ | Admin only, returns 403 for customer |
| **Product Variants** |
| | POST /api/products/{id}/variants | ✅ | Admin only, returns 403 for customer |
| | GET /api/products/{id}/variants | ✅ | Public access, returns 200 |
| | PUT /api/variants/{variant_id} | ✅ | Admin only, returns 403 for customer |
| | DELETE /api/variants/{variant_id} | ✅ | Admin only, returns 403 for customer |
| **Orders** |
| | POST /api/orders | ✅ | Authenticated users, supports variant_id |
| | GET /api/orders | ✅ | Authenticated users, returns 200 |
| | GET /api/orders/{id} | ✅ | Authenticated users |
| | PUT /api/orders/{id}/status | ✅ | Admin only, returns 403 for customer |
| **Payment Webhooks** |
| | POST /api/webhooks/payment | ✅ | Returns 401 without signature |
| **Documentation** |
| | GET /swagger/index.html | ✅ | Swagger UI accessible |

**Total: 18 endpoints, all working correctly**

## Known Limitations

1. **Admin User Creation**: The system doesn't have a public admin registration endpoint. To test admin features, you need to:
   - Manually update user role in database: `UPDATE users SET role = 'admin' WHERE email = 'user@example.com';`
   - Or use database seeding/migration scripts
   - Or implement a protected admin creation endpoint

2. **Email Validation**: The system checks for duplicate emails. Use unique emails for each test run or reset the database:
   ```bash
   docker-compose down -v
   docker-compose up -d
   ```

## Troubleshooting

### Integration Tests Fail with 404
**Issue**: Endpoints return "404 page not found"

**Solution**: Rebuild the Docker container to include latest code:
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### Tests Fail with 400 "Email already registered"
**Solution**: Reset the database:
```bash
docker-compose down -v
docker-compose up -d
```

### Can't Test Admin Features
**Solution**: Create an admin user directly in the database:
```bash
# Get user ID
docker exec -it ecommerce_postgres psql -U ecommerce -d ecommerce \
  -c "SELECT id, email, role FROM users;"

# Update role to admin
docker exec -it ecommerce_postgres psql -U ecommerce -d ecommerce \
  -c "UPDATE users SET role = 'admin' WHERE email = 'your@email.com';"
```

## CI/CD Integration

For automated testing in CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
test:
  runs-on: ubuntu-latest
  services:
    postgres:
      image: postgres:15
      env:
        POSTGRES_PASSWORD: password
  steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.24'
    - run: make test
```
