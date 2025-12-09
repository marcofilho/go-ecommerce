#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Authentication & Authorization Tests${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Test 1: Register a new customer
echo -e "${YELLOW}Test 1: Register new customer${NC}"
CUSTOMER_EMAIL="customer_$(date +%s)@example.com"
REGISTER_RESPONSE=$(curl -s -X POST ${API_URL}/api/auth/register \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$CUSTOMER_EMAIL\",
    \"password\": \"password123\",
    \"name\": \"John Customer\"
  }")

CUSTOMER_TOKEN=$(echo $REGISTER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$CUSTOMER_TOKEN" ]; then
    echo -e "${RED}✗ Failed to register customer${NC}"
    echo "Response: $REGISTER_RESPONSE"
else
    echo -e "${GREEN}✓ Customer registered successfully${NC}"
    echo "Token: ${CUSTOMER_TOKEN:0:50}..."
fi
echo ""

# Test 2: Register an admin (manually - in production you'd have a seeder or admin creation endpoint)
echo -e "${YELLOW}Test 2: Register admin user${NC}"
ADMIN_EMAIL="admin_$(date +%s)@example.com"
ADMIN_REGISTER_RESPONSE=$(curl -s -X POST ${API_URL}/api/auth/register \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$ADMIN_EMAIL\",
    \"password\": \"admin123\",
    \"name\": \"Admin User\"
  }")

# Promote user to admin role in database
docker exec ecommerce_postgres psql -U postgres -d ecommerce -c "UPDATE users SET role = 'admin' WHERE email = '$ADMIN_EMAIL';" > /dev/null 2>&1

# Re-login to get fresh token with admin role
ADMIN_LOGIN_RESPONSE=$(curl -s -X POST ${API_URL}/api/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$ADMIN_EMAIL\",
    \"password\": \"admin123\"
  }")

ADMIN_TOKEN=$(echo $ADMIN_LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$ADMIN_TOKEN" ]; then
    echo -e "${RED}✗ Failed to register admin${NC}"
else
    echo -e "${GREEN}✓ Admin registered successfully${NC}"
    echo "Token: ${ADMIN_TOKEN:0:50}..."
fi
echo ""

# Test 3: Login as customer
echo -e "${YELLOW}Test 3: Login as customer${NC}"
LOGIN_RESPONSE=$(curl -s -X POST ${API_URL}/api/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$CUSTOMER_EMAIL\",
    \"password\": \"password123\"
  }")

LOGIN_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$LOGIN_TOKEN" ]; then
    echo -e "${RED}✗ Failed to login${NC}"
else
    echo -e "${GREEN}✓ Login successful${NC}"
    echo "Token: ${LOGIN_TOKEN:0:50}..."
fi
echo ""

# Test 4: Access protected endpoint without token
echo -e "${YELLOW}Test 4: Access protected endpoint without token (should fail)${NC}"
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "products": [{"product_id": "550e8400-e29b-41d4-a716-446655440000", "quantity": 1}]
  }')

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

if [ "$HTTP_CODE" == "401" ]; then
    echo -e "${GREEN}✓ Correctly rejected (401 Unauthorized)${NC}"
else
    echo -e "${RED}✗ Unexpected response code: $HTTP_CODE${NC}"
fi
echo ""

# Test 5: Access protected endpoint with invalid token
echo -e "${YELLOW}Test 5: Access protected endpoint with invalid token (should fail)${NC}"
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer invalid-token" \
  -d '{
    "customer_id": 1,
    "products": [{"product_id": "550e8400-e29b-41d4-a716-446655440000", "quantity": 1}]
  }')

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

if [ "$HTTP_CODE" == "401" ]; then
    echo -e "${GREEN}✓ Correctly rejected (401 Unauthorized)${NC}"
else
    echo -e "${RED}✗ Unexpected response code: $HTTP_CODE${NC}"
fi
echo ""

# Test 6: Customer tries to create a product (admin only - should fail)
echo -e "${YELLOW}Test 6: Customer tries to create product (should fail - admin only)${NC}"
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -d '{
    "name": "Unauthorized Product",
    "description": "This should not be created",
    "price": 99.99,
    "quantity": 10
  }')

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

if [ "$HTTP_CODE" == "403" ]; then
    echo -e "${GREEN}✓ Correctly rejected (403 Forbidden)${NC}"
else
    echo -e "${RED}✗ Unexpected response code: $HTTP_CODE${NC}"
fi
echo ""

# Test 7: Admin creates a product (should succeed)
echo -e "${YELLOW}Test 7: Admin creates product (should succeed)${NC}"
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Authorized Product",
    "description": "Created by admin",
    "price": 149.99,
    "quantity": 20
  }')

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

if [ "$HTTP_CODE" == "201" ]; then
    echo -e "${GREEN}✓ Product created successfully${NC}"
    PRODUCT_RESPONSE=$(echo "$RESPONSE" | sed 's/HTTP_CODE.*//')
    PRODUCT_ID=$(echo $PRODUCT_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "Product ID: $PRODUCT_ID"
else
    echo -e "${RED}✗ Failed to create product. HTTP code: $HTTP_CODE${NC}"
    PRODUCT_ID=""
fi
echo ""

# Test 8: Customer lists products (public endpoint - should succeed)
echo -e "${YELLOW}Test 8: Customer lists products (public endpoint)${NC}"
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET ${API_URL}/api/products)

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

if [ "$HTTP_CODE" == "200" ]; then
    echo -e "${GREEN}✓ Products listed successfully${NC}"
else
    echo -e "${RED}✗ Failed to list products. HTTP code: $HTTP_CODE${NC}"
fi
echo ""

# Test 9: Customer creates an order (authenticated endpoint - should succeed)
echo -e "${YELLOW}Test 9: Customer creates an order (should succeed)${NC}"
if [ ! -z "$PRODUCT_ID" ]; then
    RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/orders \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $CUSTOMER_TOKEN" \
      -d "{
        \"customer_id\": 1,
        \"products\": [{\"product_id\": \"$PRODUCT_ID\", \"quantity\": 2}]
      }")

    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

    if [ "$HTTP_CODE" == "201" ]; then
        echo -e "${GREEN}✓ Order created successfully${NC}"
        ORDER_RESPONSE=$(echo "$RESPONSE" | sed 's/HTTP_CODE.*//')
        ORDER_ID=$(echo $ORDER_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
        echo "Order ID: $ORDER_ID"
    else
        echo -e "${RED}✗ Failed to create order. HTTP code: $HTTP_CODE${NC}"
        ORDER_ID=""
    fi
else
    echo -e "${YELLOW}⚠ Skipping (no product ID available)${NC}"
    ORDER_ID=""
fi
echo ""

# Test 10: Customer tries to update order status (admin only - should fail)
echo -e "${YELLOW}Test 10: Customer tries to update order status (should fail - admin only)${NC}"
if [ ! -z "$ORDER_ID" ]; then
    RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X PUT ${API_URL}/api/orders/${ORDER_ID}/status \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $CUSTOMER_TOKEN" \
      -d '{
        "status": "completed"
      }')

    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

    if [ "$HTTP_CODE" == "403" ]; then
        echo -e "${GREEN}✓ Correctly rejected (403 Forbidden)${NC}"
    else
        echo -e "${RED}✗ Unexpected response code: $HTTP_CODE${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Skipping (no order ID available)${NC}"
fi
echo ""

# Test 11: Admin updates order status (should succeed)
echo -e "${YELLOW}Test 11: Admin updates order status (should succeed)${NC}"
if [ ! -z "$ORDER_ID" ]; then
    RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X PUT ${API_URL}/api/orders/${ORDER_ID}/status \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $ADMIN_TOKEN" \
      -d '{
        "status": "completed"
      }')

    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

    if [ "$HTTP_CODE" == "200" ]; then
        echo -e "${GREEN}✓ Order status updated successfully${NC}"
    else
        echo -e "${RED}✗ Failed to update order status. HTTP code: $HTTP_CODE${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Skipping (no order ID available)${NC}"
fi
echo ""

# Test 12: Customer tries to create category (admin only - should fail)
echo -e "${YELLOW}Test 12: Customer tries to create category (should fail - admin only)${NC}"
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -d '{
    "name": "Unauthorized Category"
  }')

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

if [ "$HTTP_CODE" == "403" ]; then
    echo -e "${GREEN}✓ Correctly rejected (403 Forbidden)${NC}"
else
    echo -e "${RED}✗ Unexpected response code: $HTTP_CODE${NC}"
fi
echo ""

# Test 13: Admin creates a category (should succeed)
echo -e "${YELLOW}Test 13: Admin creates category (should succeed)${NC}"
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Electronics"
  }')

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

if [ "$HTTP_CODE" == "201" ]; then
    echo -e "${GREEN}✓ Category created successfully${NC}"
    CATEGORY_RESPONSE=$(echo "$RESPONSE" | sed 's/HTTP_CODE.*//')
    CATEGORY_ID=$(echo $CATEGORY_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "Category ID: $CATEGORY_ID"
else
    echo -e "${RED}✗ Failed to create category. HTTP code: $HTTP_CODE${NC}"
    CATEGORY_ID=""
fi
echo ""

# Test 14: Customer lists categories (public endpoint - should succeed)
echo -e "${YELLOW}Test 14: Customer lists categories (public endpoint)${NC}"
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET ${API_URL}/api/categories)

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

if [ "$HTTP_CODE" == "200" ]; then
    echo -e "${GREEN}✓ Categories listed successfully${NC}"
else
    echo -e "${RED}✗ Failed to list categories. HTTP code: $HTTP_CODE${NC}"
fi
echo ""

# Test 15: Admin assigns category to product (should succeed)
echo -e "${YELLOW}Test 15: Admin assigns category to product (should succeed)${NC}"
if [ ! -z "$PRODUCT_ID" ] && [ ! -z "$CATEGORY_ID" ]; then
    RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/products/${PRODUCT_ID}/categories \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $ADMIN_TOKEN" \
      -d "{
        \"category_id\": \"$CATEGORY_ID\"
      }")

    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

    if [ "$HTTP_CODE" == "200" ]; then
        echo -e "${GREEN}✓ Category assigned to product successfully${NC}"
    else
        echo -e "${RED}✗ Failed to assign category. HTTP code: $HTTP_CODE${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Skipping (no product or category ID available)${NC}"
fi
echo ""

# Test 16: Customer gets product categories (public - should succeed)
echo -e "${YELLOW}Test 16: Customer gets product categories (public endpoint)${NC}"
if [ ! -z "$PRODUCT_ID" ]; then
    RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET ${API_URL}/api/products/${PRODUCT_ID}/categories)

    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

    if [ "$HTTP_CODE" == "200" ]; then
        echo -e "${GREEN}✓ Product categories retrieved successfully${NC}"
        CATEGORIES_RESPONSE=$(echo "$RESPONSE" | sed 's/HTTP_CODE.*//')
        echo "Categories: $CATEGORIES_RESPONSE"
    else
        echo -e "${RED}✗ Failed to get product categories. HTTP code: $HTTP_CODE${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Skipping (no product ID available)${NC}"
fi
echo ""

# Test 17: Customer tries to remove category from product (admin only - should fail)
echo -e "${YELLOW}Test 17: Customer tries to remove category (should fail - admin only)${NC}"
if [ ! -z "$PRODUCT_ID" ] && [ ! -z "$CATEGORY_ID" ]; then
    RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X DELETE ${API_URL}/api/products/${PRODUCT_ID}/categories/${CATEGORY_ID} \
      -H "Authorization: Bearer $CUSTOMER_TOKEN")

    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

    if [ "$HTTP_CODE" == "403" ]; then
        echo -e "${GREEN}✓ Correctly rejected (403 Forbidden)${NC}"
    else
        echo -e "${RED}✗ Unexpected response code: $HTTP_CODE${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Skipping (no product or category ID available)${NC}"
fi
echo ""

# Summary
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "All authentication and authorization tests completed!"
echo -e "\n${GREEN}Key Features Demonstrated:${NC}"
echo -e "✓ User registration and login"
echo -e "✓ JWT token generation"
echo -e "✓ Token validation in middleware"
echo -e "✓ Role-based access control (admin vs customer)"
echo -e "✓ Public vs protected endpoints"
echo -e "✓ Product categories (N:N relationship)"
echo -e "✓ Category assignment to products"
echo -e "✓ Proper HTTP status codes (401, 403)"
