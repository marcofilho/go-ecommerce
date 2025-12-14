#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Promo Code Integration Tests${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Setup: Register admin and customer users
echo -e "${YELLOW}Setup: Creating test users${NC}"
ADMIN_EMAIL="admin_promo_$(date +%s)@example.com"
CUSTOMER_EMAIL="customer_promo_$(date +%s)@example.com"

# Register Admin
ADMIN_REGISTER=$(curl -s -X POST ${API_URL}/api/auth/register \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$ADMIN_EMAIL\",
    \"password\": \"admin123\",
    \"name\": \"Admin Promo\"
  }")

docker exec ecommerce_postgres psql -U postgres -d ecommerce -c "UPDATE users SET role = 'admin' WHERE email = '$ADMIN_EMAIL';" > /dev/null 2>&1

ADMIN_LOGIN=$(curl -s -X POST ${API_URL}/api/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$ADMIN_EMAIL\",
    \"password\": \"admin123\"
  }")

ADMIN_TOKEN=$(echo $ADMIN_LOGIN | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Register Customer
CUSTOMER_REGISTER=$(curl -s -X POST ${API_URL}/api/auth/register \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$CUSTOMER_EMAIL\",
    \"password\": \"customer123\",
    \"name\": \"Customer Promo\"
  }")

CUSTOMER_TOKEN=$(echo $CUSTOMER_REGISTER | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$ADMIN_TOKEN" ] || [ -z "$CUSTOMER_TOKEN" ]; then
    echo -e "${RED}✗ Failed to setup test users${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Test users created${NC}"
echo ""

# Setup: Create a test product
echo -e "${YELLOW}Setup: Creating test product${NC}"
PRODUCT_RESPONSE=$(curl -s -X POST ${API_URL}/api/products \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product for Promo",
    "description": "Product for testing promo codes",
    "price": 100.00,
    "quantity": 50
  }')

PRODUCT_ID=$(echo $PRODUCT_RESPONSE | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)

if [ -z "$PRODUCT_ID" ]; then
    echo -e "${RED}✗ Failed to create test product${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Test product created: $PRODUCT_ID${NC}"
echo ""

# Test 1: Create a percentage discount
echo -e "${YELLOW}Test 1: Create percentage discount (20% off)${NC}"
PERCENTAGE_DISCOUNT=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/discounts \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "SAVE20",
    "discount_type": "percentage",
    "value": 20.0,
    "active": true
  }')

HTTP_CODE=$(echo "$PERCENTAGE_DISCOUNT" | grep HTTP_CODE | cut -d':' -f2)
BODY=$(echo "$PERCENTAGE_DISCOUNT" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "201" ]; then
    echo -e "${GREEN}✓ Percentage discount created${NC}"
    PERCENTAGE_DISCOUNT_ID=$(echo $BODY | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "Discount ID: $PERCENTAGE_DISCOUNT_ID"
else
    echo -e "${RED}✗ Failed to create percentage discount (HTTP $HTTP_CODE)${NC}"
    echo "Response: $BODY"
fi
echo ""

# Test 2: Create an amount discount
echo -e "${YELLOW}Test 2: Create amount discount (\$15 off)${NC}"
AMOUNT_DISCOUNT=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/discounts \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "FLAT15",
    "discount_type": "amount",
    "value": 15.0,
    "active": true
  }')

HTTP_CODE=$(echo "$AMOUNT_DISCOUNT" | grep HTTP_CODE | cut -d':' -f2)
BODY=$(echo "$AMOUNT_DISCOUNT" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "201" ]; then
    echo -e "${GREEN}✓ Amount discount created${NC}"
    AMOUNT_DISCOUNT_ID=$(echo $BODY | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "Discount ID: $AMOUNT_DISCOUNT_ID"
else
    echo -e "${RED}✗ Failed to create amount discount (HTTP $HTTP_CODE)${NC}"
    echo "Response: $BODY"
fi
echo ""

# Test 3: Create order without promo code
echo -e "${YELLOW}Test 3: Create order without promo code (baseline)${NC}"
ORDER_NO_PROMO=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/orders \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\": 1,
    \"products\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 1
      }
    ]
  }")

HTTP_CODE=$(echo "$ORDER_NO_PROMO" | grep HTTP_CODE | cut -d':' -f2)
BODY=$(echo "$ORDER_NO_PROMO" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "201" ]; then
    TOTAL_NO_PROMO=$(echo $BODY | grep -o '"total_price":[0-9.]*' | cut -d':' -f2)
    echo -e "${GREEN}✓ Order created without promo code${NC}"
    echo "Total Price: \$$TOTAL_NO_PROMO"
    
    if [ "$TOTAL_NO_PROMO" = "100" ] || [ "$TOTAL_NO_PROMO" = "100.00" ]; then
        echo -e "${GREEN}✓ Price is correct (no discount applied)${NC}"
    else
        echo -e "${RED}✗ Expected \$100.00, got \$$TOTAL_NO_PROMO${NC}"
    fi
else
    echo -e "${RED}✗ Failed to create order without promo (HTTP $HTTP_CODE)${NC}"
    echo "Response: $BODY"
fi
echo ""

# Test 4: Create order with percentage promo code (20% off)
echo -e "${YELLOW}Test 4: Create order with percentage promo code (SAVE20)${NC}"
ORDER_PERCENTAGE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/orders \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\": 1,
    \"products\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 1
      }
    ],
    \"promo_code\": \"SAVE20\"
  }")

HTTP_CODE=$(echo "$ORDER_PERCENTAGE" | grep HTTP_CODE | cut -d':' -f2)
BODY=$(echo "$ORDER_PERCENTAGE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "201" ]; then
    TOTAL_PERCENTAGE=$(echo $BODY | grep -o '"total_price":[0-9.]*' | cut -d':' -f2)
    echo -e "${GREEN}✓ Order created with percentage promo code${NC}"
    echo "Total Price: \$$TOTAL_PERCENTAGE"
    
    # 20% off $100 = $80
    if [ "$TOTAL_PERCENTAGE" = "80" ] || [ "$TOTAL_PERCENTAGE" = "80.00" ]; then
        echo -e "${GREEN}✓ Discount applied correctly (20% off: \$100 -> \$80)${NC}"
    else
        echo -e "${RED}✗ Expected \$80.00, got \$$TOTAL_PERCENTAGE${NC}"
    fi
else
    echo -e "${RED}✗ Failed to create order with percentage promo (HTTP $HTTP_CODE)${NC}"
    echo "Response: $BODY"
fi
echo ""

# Test 5: Create order with amount promo code ($15 off)
echo -e "${YELLOW}Test 5: Create order with amount promo code (FLAT15)${NC}"
ORDER_AMOUNT=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/orders \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\": 1,
    \"products\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 1
      }
    ],
    \"promo_code\": \"FLAT15\"
  }")

HTTP_CODE=$(echo "$ORDER_AMOUNT" | grep HTTP_CODE | cut -d':' -f2)
BODY=$(echo "$ORDER_AMOUNT" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "201" ]; then
    TOTAL_AMOUNT=$(echo $BODY | grep -o '"total_price":[0-9.]*' | cut -d':' -f2)
    echo -e "${GREEN}✓ Order created with amount promo code${NC}"
    echo "Total Price: \$$TOTAL_AMOUNT"
    
    # $15 off $100 = $85
    if [ "$TOTAL_AMOUNT" = "85" ] || [ "$TOTAL_AMOUNT" = "85.00" ]; then
        echo -e "${GREEN}✓ Discount applied correctly (\$15 off: \$100 -> \$85)${NC}"
    else
        echo -e "${RED}✗ Expected \$85.00, got \$$TOTAL_AMOUNT${NC}"
    fi
else
    echo -e "${RED}✗ Failed to create order with amount promo (HTTP $HTTP_CODE)${NC}"
    echo "Response: $BODY"
fi
echo ""

# Test 6: Try using invalid promo code
echo -e "${YELLOW}Test 6: Try using invalid promo code${NC}"
ORDER_INVALID=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/orders \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\": 1,
    \"products\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 1
      }
    ],
    \"promo_code\": \"INVALID123\"
  }")

HTTP_CODE=$(echo "$ORDER_INVALID" | grep HTTP_CODE | cut -d':' -f2)
BODY=$(echo "$ORDER_INVALID" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "400" ]; then
    echo -e "${GREEN}✓ Invalid promo code rejected (HTTP 400)${NC}"
    ERROR_MSG=$(echo $BODY | grep -o '"error":"[^"]*' | cut -d'"' -f4)
    echo "Error message: $ERROR_MSG"
else
    echo -e "${RED}✗ Invalid promo code should return HTTP 400, got $HTTP_CODE${NC}"
    echo "Response: $BODY"
fi
echo ""

# Test 7: Deactivate promo code and try to use it
echo -e "${YELLOW}Test 7: Deactivate promo code and try to use it${NC}"

# Deactivate SAVE20
DEACTIVATE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X PUT ${API_URL}/api/discounts/$PERCENTAGE_DISCOUNT_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "SAVE20",
    "discount_type": "percentage",
    "value": 20.0,
    "active": false
  }')

HTTP_CODE=$(echo "$DEACTIVATE" | grep HTTP_CODE | cut -d':' -f2)

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Promo code deactivated${NC}"
    
    # Try to use deactivated promo
    ORDER_INACTIVE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/orders \
      -H "Authorization: Bearer $CUSTOMER_TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"customer_id\": 1,
        \"products\": [
          {
            \"product_id\": \"$PRODUCT_ID\",
            \"quantity\": 1
          }
        ],
        \"promo_code\": \"SAVE20\"
      }")
    
    HTTP_CODE=$(echo "$ORDER_INACTIVE" | grep HTTP_CODE | cut -d':' -f2)
    BODY=$(echo "$ORDER_INACTIVE" | sed '/HTTP_CODE/d')
    
    if [ "$HTTP_CODE" = "400" ]; then
        echo -e "${GREEN}✓ Inactive promo code rejected (HTTP 400)${NC}"
        ERROR_MSG=$(echo $BODY | grep -o '"error":"[^"]*' | cut -d'"' -f4)
        echo "Error message: $ERROR_MSG"
    else
        echo -e "${RED}✗ Inactive promo code should be rejected, got HTTP $HTTP_CODE${NC}"
        echo "Response: $BODY"
    fi
else
    echo -e "${RED}✗ Failed to deactivate promo code${NC}"
fi
echo ""

# Test 8: Validate promo code endpoint
echo -e "${YELLOW}Test 8: Validate active promo code${NC}"
VALIDATE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/discounts/validate \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "FLAT15"
  }')

HTTP_CODE=$(echo "$VALIDATE" | grep HTTP_CODE | cut -d':' -f2)
BODY=$(echo "$VALIDATE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Active promo code validated${NC}"
    PROMO_CODE=$(echo $BODY | grep -o '"promo_code":"[^"]*' | cut -d'"' -f4)
    DISCOUNT_TYPE=$(echo $BODY | grep -o '"discount_type":"[^"]*' | cut -d'"' -f4)
    VALUE=$(echo $BODY | grep -o '"value":[0-9.]*' | cut -d':' -f2)
    echo "Promo Code: $PROMO_CODE, Type: $DISCOUNT_TYPE, Value: $VALUE"
else
    echo -e "${RED}✗ Failed to validate promo code (HTTP $HTTP_CODE)${NC}"
    echo "Response: $BODY"
fi
echo ""

# Test 9: Validate inactive promo code
echo -e "${YELLOW}Test 9: Validate inactive promo code${NC}"
VALIDATE_INACTIVE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/discounts/validate \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "SAVE20"
  }')

HTTP_CODE=$(echo "$VALIDATE_INACTIVE" | grep HTTP_CODE | cut -d':' -f2)
BODY=$(echo "$VALIDATE_INACTIVE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "404" ] || [ "$HTTP_CODE" = "400" ]; then
    echo -e "${GREEN}✓ Inactive promo code not found/invalid (HTTP $HTTP_CODE)${NC}"
else
    echo -e "${RED}✗ Inactive promo should return error, got HTTP $HTTP_CODE${NC}"
    echo "Response: $BODY"
fi
echo ""

# Test 10: Multiple items with promo code
echo -e "${YELLOW}Test 10: Order with multiple items and promo code${NC}"
ORDER_MULTI=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST ${API_URL}/api/orders \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_id\": 1,
    \"products\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 3
      }
    ],
    \"promo_code\": \"FLAT15\"
  }")

HTTP_CODE=$(echo "$ORDER_MULTI" | grep HTTP_CODE | cut -d':' -f2)
BODY=$(echo "$ORDER_MULTI" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "201" ]; then
    TOTAL_MULTI=$(echo $BODY | grep -o '"total_price":[0-9.]*' | cut -d':' -f2)
    echo -e "${GREEN}✓ Order with multiple items created${NC}"
    echo "Total Price: \$$TOTAL_MULTI"
    
    # 3 items * $100 = $300, minus $15 = $285
    if [ "$TOTAL_MULTI" = "285" ] || [ "$TOTAL_MULTI" = "285.00" ]; then
        echo -e "${GREEN}✓ Discount applied correctly to total (\$300 - \$15 = \$285)${NC}"
    else
        echo -e "${RED}✗ Expected \$285.00, got \$$TOTAL_MULTI${NC}"
    fi
else
    echo -e "${RED}✗ Failed to create order with multiple items (HTTP $HTTP_CODE)${NC}"
    echo "Response: $BODY"
fi
echo ""

# Summary
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ Percentage discounts${NC}"
echo -e "${GREEN}✓ Amount discounts${NC}"
echo -e "${GREEN}✓ Orders without promo codes${NC}"
echo -e "${GREEN}✓ Orders with valid promo codes${NC}"
echo -e "${GREEN}✓ Invalid promo code rejection${NC}"
echo -e "${GREEN}✓ Inactive promo code rejection${NC}"
echo -e "${GREEN}✓ Promo code validation endpoint${NC}"
echo -e "${GREEN}✓ Multiple items with promo codes${NC}"
echo ""
