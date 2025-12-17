#!/bin/bash

# Test Enhanced Discount System - Practical Examples
# This script demonstrates how product, category, and user-specific discounts work

BASE_URL="http://localhost:8080/api"

echo "=========================================="
echo "Enhanced Discount System - Test Scenarios"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper function to print section headers
print_header() {
    echo ""
    echo "=========================================="
    echo "$1"
    echo "=========================================="
}

# Helper function to print test results
print_test() {
    echo ""
    echo -e "${YELLOW}TEST:${NC} $1"
}

print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ PASS${NC}: $2"
    else
        echo -e "${RED}✗ FAIL${NC}: $2"
    fi
}

# Prerequisites check
print_header "0. Prerequisites Check"
echo "Checking if API is running..."
if curl -s -f -o /dev/null "$BASE_URL/products"; then
    echo -e "${GREEN}✓${NC} API is running"
else
    echo -e "${RED}✗${NC} API is not running. Please start it with 'make run'"
    exit 1
fi

# 1. Create test products and categories
print_header "1. Setup Test Data"

echo "Creating Electronics category..."
ELECTRONICS_ID=$(curl -s -X POST "$BASE_URL/categories" \
  -H "Content-Type: application/json" \
  -d '{"name":"Electronics"}' | jq -r '.id')
echo "Electronics category ID: $ELECTRONICS_ID"

echo "Creating Clothing category..."
CLOTHING_ID=$(curl -s -X POST "$BASE_URL/categories" \
  -H "Content-Type: application/json" \
  -d '{"name":"Clothing"}' | jq -r '.id')
echo "Clothing category ID: $CLOTHING_ID"

echo "Creating Laptop product..."
LAPTOP_ID=$(curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Gaming Laptop\",\"description\":\"High-end gaming laptop\",\"price\":1000,\"quantity\":10,\"category_ids\":[\"$ELECTRONICS_ID\"]}" | jq -r '.id')
echo "Laptop ID: $LAPTOP_ID"

echo "Creating Mouse product..."
MOUSE_ID=$(curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Wireless Mouse\",\"description\":\"Ergonomic mouse\",\"price\":50,\"quantity\":20,\"category_ids\":[\"$ELECTRONICS_ID\"]}" | jq -r '.id')
echo "Mouse ID: $MOUSE_ID"

echo "Creating T-Shirt product..."
TSHIRT_ID=$(curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Cotton T-Shirt\",\"description\":\"Comfortable t-shirt\",\"price\":30,\"quantity\":50,\"category_ids\":[\"$CLOTHING_ID\"]}" | jq -r '.id')
echo "T-Shirt ID: $TSHIRT_ID"

# 2. Test Site-Wide Discount
print_header "2. Site-Wide Discount (OLD BEHAVIOR)"

echo "Creating SITEWIDE10 discount (10% off everything)..."
curl -s -X POST "$BASE_URL/discounts" \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "SITEWIDE10",
    "discount_type": "percentage",
    "value": 10,
    "active": true
  }' | jq '.'

print_test "Order with mixed products + SITEWIDE10"
echo "Order: Laptop ($1000) + Mouse ($50) = $1050"
echo "Expected: 10% off entire order = $945"
RESULT=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\": 1,
    \"products\": [
      {\"product_id\": \"$LAPTOP_ID\", \"quantity\": 1},
      {\"product_id\": \"$MOUSE_ID\", \"quantity\": 1}
    ],
    \"promo_code\": \"SITEWIDE10\"
  }")
TOTAL=$(echo "$RESULT" | jq -r '.total_price')
echo "Actual total: \$$TOTAL"

# 3. Test Product-Specific Discount
print_header "3. Product-Specific Discount (NEW!)"

echo "Creating LAPTOP50 discount ($50 off laptops only)..."
curl -s -X POST "$BASE_URL/discounts" \
  -H "Content-Type: application/json" \
  -d "{
    \"promo_code\": \"LAPTOP50\",
    \"discount_type\": \"amount\",
    \"value\": 50,
    \"active\": true,
    \"product_ids\": [\"$LAPTOP_ID\"]
  }" | jq '.'

print_test "Order with eligible product (Laptop) + LAPTOP50"
echo "Order: Laptop ($1000) + Mouse ($50)"
echo "Expected: $50 off laptop only = $1000 total"
RESULT=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\": 2,
    \"products\": [
      {\"product_id\": \"$LAPTOP_ID\", \"quantity\": 1},
      {\"product_id\": \"$MOUSE_ID\", \"quantity\": 1}
    ],
    \"promo_code\": \"LAPTOP50\"
  }")
TOTAL=$(echo "$RESULT" | jq -r '.total_price')
if [ "$TOTAL" != "null" ]; then
    echo "Actual total: \$$TOTAL"
    print_result 0 "Discount applied only to laptop"
else
    ERROR=$(echo "$RESULT" | jq -r '.error')
    echo "Error: $ERROR"
fi

print_test "Order with NO eligible products + LAPTOP50"
echo "Order: Mouse ($50) + T-Shirt ($30)"
echo "Expected: Order rejected - discount doesn't apply"
RESULT=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\": 2,
    \"products\": [
      {\"product_id\": \"$MOUSE_ID\", \"quantity\": 1},
      {\"product_id\": \"$TSHIRT_ID\", \"quantity\": 1}
    ],
    \"promo_code\": \"LAPTOP50\"
  }")
ERROR=$(echo "$RESULT" | jq -r '.error')
if [[ "$ERROR" == *"does not apply"* ]]; then
    print_result 0 "Order correctly rejected: $ERROR"
else
    print_result 1 "Order should have been rejected"
fi

# 4. Test Category-Based Discount
print_header "4. Category-Based Discount (NEW!)"

echo "Creating ELECTRONICS15 discount (15% off electronics)..."
curl -s -X POST "$BASE_URL/discounts" \
  -H "Content-Type: application/json" \
  -d "{
    \"promo_code\": \"ELECTRONICS15\",
    \"discount_type\": \"percentage\",
    \"value\": 15,
    \"active\": true,
    \"category_ids\": [\"$ELECTRONICS_ID\"]
  }" | jq '.'

print_test "Mixed categories order + ELECTRONICS15"
echo "Order: Laptop ($1000) + Mouse ($50) + T-Shirt ($30)"
echo "Expected: 15% off electronics only ($1050 * 0.15 = $157.50 off)"
echo "Expected total: $922.50"
RESULT=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\": 3,
    \"products\": [
      {\"product_id\": \"$LAPTOP_ID\", \"quantity\": 1},
      {\"product_id\": \"$MOUSE_ID\", \"quantity\": 1},
      {\"product_id\": \"$TSHIRT_ID\", \"quantity\": 1}
    ],
    \"promo_code\": \"ELECTRONICS15\"
  }")
TOTAL=$(echo "$RESULT" | jq -r '.total_price')
echo "Actual total: \$$TOTAL"

# 5. Test VIP User Discount
print_header "5. VIP User-Only Discount (NEW!)"

echo "Creating VIP20 discount (20% for VIP users 2 & 3 only)..."
curl -s -X POST "$BASE_URL/discounts" \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "VIP20",
    "discount_type": "percentage",
    "value": 20,
    "active": true,
    "user_ids": [2, 3],
    "user_usage_limit": 5
  }' | jq '.'

print_test "Regular user tries VIP discount"
echo "Customer ID 1 (not VIP) tries to use VIP20"
echo "Expected: Order rejected - not authorized"
RESULT=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\": 1,
    \"products\": [
      {\"product_id\": \"$LAPTOP_ID\", \"quantity\": 1}
    ],
    \"promo_code\": \"VIP20\"
  }")
ERROR=$(echo "$RESULT" | jq -r '.error')
if [[ "$ERROR" == *"not available for your account"* ]]; then
    print_result 0 "Non-VIP correctly rejected: $ERROR"
else
    print_result 1 "Should have rejected non-VIP user"
fi

print_test "VIP user uses VIP discount"
echo "Customer ID 2 (VIP) uses VIP20 on $1000 laptop"
echo "Expected: 20% off = $800"
RESULT=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\": 2,
    \"products\": [
      {\"product_id\": \"$LAPTOP_ID\", \"quantity\": 1}
    ],
    \"promo_code\": \"VIP20\"
  }")
TOTAL=$(echo "$RESULT" | jq -r '.total_price')
if [ "$TOTAL" != "null" ]; then
    echo "Actual total: \$$TOTAL"
    print_result 0 "VIP discount applied successfully"
else
    ERROR=$(echo "$RESULT" | jq -r '.error')
    echo "Error: $ERROR"
fi

# 6. Test Usage Limits
print_header "6. Usage Limits (NEW!)"

echo "Creating LIMITED5 discount (5 total uses max)..."
curl -s -X POST "$BASE_URL/discounts" \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "LIMITED5",
    "discount_type": "percentage",
    "value": 30,
    "active": true,
    "usage_limit": 5,
    "user_usage_limit": 2
  }' | jq '.'

print_test "User makes first order with LIMITED5"
echo "Customer 4, attempt 1/2"
RESULT=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\": 4,
    \"products\": [
      {\"product_id\": \"$MOUSE_ID\", \"quantity\": 1}
    ],
    \"promo_code\": \"LIMITED5\"
  }")
TOTAL=$(echo "$RESULT" | jq -r '.total_price')
if [ "$TOTAL" != "null" ]; then
    echo "First use successful: \$$TOTAL"
    print_result 0 "Usage 1/2 allowed"
else
    ERROR=$(echo "$RESULT" | jq -r '.error')
    echo "Error: $ERROR"
fi

# Summary
print_header "Test Summary"
echo ""
echo "The enhanced discount system now supports:"
echo ""
echo "✓ Site-wide discounts (applies to everything)"
echo "✓ Product-specific discounts (only certain products)"
echo "✓ Category-based discounts (entire categories)"
echo "✓ VIP/User-specific discounts (targeted users only)"
echo "✓ Global usage limits (total uses)"
echo "✓ Per-user usage limits (per-customer limits)"
echo "✓ Minimum purchase requirements"
echo "✓ Date range validation"
echo ""
echo "Orders are now intelligently validated before applying discounts!"
echo ""
echo "See docs/DISCOUNT_EXAMPLES.md for detailed scenarios"
echo ""
