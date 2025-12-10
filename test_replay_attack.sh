#!/bin/bash

API_URL="http://localhost:8080"
WEBHOOK_SECRET="webhook_secret_key"

# Login as customer to create an order
CUSTOMER_TOKEN=$(curl -s -X POST "$API_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"password123"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Get product ID (from seeded data)
PRODUCT_ID=$(curl -s "$API_URL/api/products" | grep -o '"ID":"[^"]*' | head -1 | cut -d'"' -f4)

# Create an order
ORDER_ID=$(curl -s -X POST "$API_URL/api/orders" \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"items\":[{\"product_id\":\"$PRODUCT_ID\",\"quantity\":1}]}" | grep -o '"ID":"[^"]*' | head -1 | cut -d'"' -f4)

echo "Created order: $ORDER_ID"

# Test 1: Timestamp 10 minutes in the past (should be rejected)
echo -e "\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "TEST 1: Webhook with old timestamp (10 minutes ago)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
OLD_TIMESTAMP=$(($(date +%s) - 600)) # 10 minutes ago
PAYLOAD="{\"order_id\":\"$ORDER_ID\",\"timestamp\":$OLD_TIMESTAMP,\"transaction_id\":\"txn_old\",\"payment_status\":\"paid\"}"
SIGNATURE=$(echo -n "$PAYLOAD" | openssl dgst -sha256 -hmac "$WEBHOOK_SECRET" | sed 's/^.* //')

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Payment-Signature: $SIGNATURE" \
  -d "$PAYLOAD")

if echo "$RESPONSE" | grep -q "401"; then
  echo "✓ OLD timestamp rejected with 401 (replay attack prevented)"
else
  echo "✗ Expected 401, got: $RESPONSE"
fi

# Test 2: Timestamp 10 minutes in the future (should be rejected)
echo -e "\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "TEST 2: Webhook with future timestamp (10 minutes ahead)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
FUTURE_TIMESTAMP=$(($(date +%s) + 600)) # 10 minutes in future
PAYLOAD="{\"order_id\":\"$ORDER_ID\",\"timestamp\":$FUTURE_TIMESTAMP,\"transaction_id\":\"txn_future\",\"payment_status\":\"paid\"}"
SIGNATURE=$(echo -n "$PAYLOAD" | openssl dgst -sha256 -hmac "$WEBHOOK_SECRET" | sed 's/^.* //')

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Payment-Signature: $SIGNATURE" \
  -d "$PAYLOAD")

if echo "$RESPONSE" | grep -q "401"; then
  echo "✓ FUTURE timestamp rejected with 401 (clock skew protection)"
else
  echo "✗ Expected 401, got: $RESPONSE"
fi

# Test 3: Current timestamp (should be accepted)
echo -e "\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "TEST 3: Webhook with current timestamp (should succeed)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
CURRENT_TIMESTAMP=$(date +%s)
PAYLOAD="{\"order_id\":\"$ORDER_ID\",\"timestamp\":$CURRENT_TIMESTAMP,\"transaction_id\":\"txn_current\",\"payment_status\":\"paid\"}"
SIGNATURE=$(echo -n "$PAYLOAD" | openssl dgst -sha256 -hmac "$WEBHOOK_SECRET" | sed 's/^.* //')

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Payment-Signature: $SIGNATURE" \
  -d "$PAYLOAD")

if echo "$RESPONSE" | grep -q "200"; then
  echo "✓ CURRENT timestamp accepted with 200 (within tolerance)"
else
  echo "✗ Expected 200, got: $RESPONSE"
fi

echo -e "\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✓ Replay attack prevention tests complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
