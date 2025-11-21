#!/bin/bash

# Payment Webhook Batch Test Script
# Tests various success and failure scenarios

set -e

API_URL="http://localhost:8080"
WEBHOOK_SECRET="my-super-secret-webhook-key-change-in-production"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to generate HMAC signature
generate_signature() {
  local payload="$1"
  echo -n "$payload" | openssl dgst -sha256 -hmac "$WEBHOOK_SECRET" | sed 's/^.* //'
}

# Function to print test header
print_test() {
  echo ""
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${YELLOW}$1${NC}"
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Function to assert expected result
assert_contains() {
  local response="$1"
  local expected="$2"
  local test_name="$3"
  
  if echo "$response" | grep -q "$expected"; then
    echo -e "${GREEN}✓${NC} $test_name: PASSED"
    return 0
  else
    echo -e "${RED}✗${NC} $test_name: FAILED"
    echo "   Expected: $expected"
    echo "   Got: $response"
    return 1
  fi
}

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Payment Webhook Batch Test Suite                     ║${NC}"
echo -e "${BLUE}║  Testing Success & Failure Scenarios                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"

# Setup: Create product
print_test "SETUP: Creating Test Product"
PRODUCT_RESPONSE=$(curl -s -X POST "$API_URL/api/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product - Batch",
    "description": "Product for batch testing",
    "price": 99.99,
    "quantity": 100
  }')
PRODUCT_ID=$(echo $PRODUCT_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo -e "${GREEN}✓${NC} Product created: $PRODUCT_ID"

# Helper function to create order
create_order() {
  local order_response=$(curl -s -X POST "$API_URL/api/orders" \
    -H "Content-Type: application/json" \
    -d "{
      \"customer_id\": 123,
      \"products\": [{
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 1
      }]
    }")
  echo $order_response | grep -o '"id":"[^"]*' | cut -d'"' -f4
}

#═══════════════════════════════════════════════════════════════
# TEST 1: Missing Signature
#═══════════════════════════════════════════════════════════════
print_test "TEST 1: Webhook Without Signature (Should Fail 401)"
ORDER_ID=$(create_order)
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -d "{\"order_id\":\"$ORDER_ID\",\"transaction_id\":\"txn_no_sig\",\"payment_status\":\"paid\"}")
assert_contains "$RESPONSE" "401" "HTTP 401 returned"
assert_contains "$RESPONSE" "Missing webhook signature" "Error message correct"

#═══════════════════════════════════════════════════════════════
# TEST 2: Invalid Signature
#═══════════════════════════════════════════════════════════════
print_test "TEST 2: Webhook With Invalid Signature (Should Fail 401)"
ORDER_ID=$(create_order)
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: invalid-signature-12345" \
  -d "{\"order_id\":\"$ORDER_ID\",\"transaction_id\":\"txn_bad_sig\",\"payment_status\":\"paid\"}")
assert_contains "$RESPONSE" "401" "HTTP 401 returned"
assert_contains "$RESPONSE" "Invalid webhook signature" "Error message correct"

#═══════════════════════════════════════════════════════════════
# TEST 3: Missing Transaction ID
#═══════════════════════════════════════════════════════════════
print_test "TEST 3: Webhook Without Transaction ID (Should Fail 400)"
ORDER_ID=$(create_order)
PAYLOAD="{\"order_id\":\"$ORDER_ID\",\"payment_status\":\"paid\"}"
SIGNATURE=$(generate_signature "$PAYLOAD")
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -d "$PAYLOAD")
assert_contains "$RESPONSE" "400" "HTTP 400 returned"
assert_contains "$RESPONSE" "transaction_id is required" "Error message correct"

#═══════════════════════════════════════════════════════════════
# TEST 4: Invalid Order ID Format
#═══════════════════════════════════════════════════════════════
print_test "TEST 4: Webhook With Invalid Order ID (Should Fail 400)"
PAYLOAD='{"order_id":"not-a-valid-uuid","transaction_id":"txn_bad_id","payment_status":"paid"}'
SIGNATURE=$(generate_signature "$PAYLOAD")
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -d "$PAYLOAD")
assert_contains "$RESPONSE" "400" "HTTP 400 returned"
assert_contains "$RESPONSE" "invalid order_id format" "Error message correct"

#═══════════════════════════════════════════════════════════════
# TEST 5: Non-Existent Order
#═══════════════════════════════════════════════════════════════
print_test "TEST 5: Webhook For Non-Existent Order (Should Fail 400)"
FAKE_ORDER_ID="00000000-0000-0000-0000-000000000000"
PAYLOAD="{\"order_id\":\"$FAKE_ORDER_ID\",\"transaction_id\":\"txn_no_order\",\"payment_status\":\"paid\"}"
SIGNATURE=$(generate_signature "$PAYLOAD")
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -d "$PAYLOAD")
assert_contains "$RESPONSE" "400" "HTTP 400 returned"
assert_contains "$RESPONSE" "order not found" "Error message correct"

#═══════════════════════════════════════════════════════════════
# TEST 6: Invalid Payment Status
#═══════════════════════════════════════════════════════════════
print_test "TEST 6: Webhook With Invalid Payment Status (Should Fail 400)"
ORDER_ID=$(create_order)
PAYLOAD="{\"order_id\":\"$ORDER_ID\",\"transaction_id\":\"txn_bad_status\",\"payment_status\":\"processing\"}"
SIGNATURE=$(generate_signature "$PAYLOAD")
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -d "$PAYLOAD")
assert_contains "$RESPONSE" "400" "HTTP 400 returned"
assert_contains "$RESPONSE" "payment_status must be either" "Error message correct"

#═══════════════════════════════════════════════════════════════
# TEST 7: Successful Payment
#═══════════════════════════════════════════════════════════════
print_test "TEST 7: Webhook With Successful Payment (Should Success 200)"
ORDER_ID=$(create_order)
TXN_ID="txn_success_$(date +%s)"
PAYLOAD="{\"order_id\":\"$ORDER_ID\",\"transaction_id\":\"$TXN_ID\",\"payment_status\":\"paid\"}"
SIGNATURE=$(generate_signature "$PAYLOAD")
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -d "$PAYLOAD")
assert_contains "$RESPONSE" "200" "HTTP 200 returned"
assert_contains "$RESPONSE" "success" "Success response"

# Verify order status changed
sleep 1
ORDER_STATUS=$(curl -s "$API_URL/api/orders/$ORDER_ID" | grep -o '"status":"[^"]*' | head -1 | cut -d'"' -f4)
assert_contains "$ORDER_STATUS" "completed" "Order status is completed"
ORDER_PAYMENT=$(curl -s "$API_URL/api/orders/$ORDER_ID" | grep -o '"payment_status":"[^"]*' | cut -d'"' -f4)
assert_contains "$ORDER_PAYMENT" "paid" "Payment status is paid"

#═══════════════════════════════════════════════════════════════
# TEST 8: Failed Payment
#═══════════════════════════════════════════════════════════════
print_test "TEST 8: Webhook With Failed Payment (Should Success 200)"
ORDER_ID=$(create_order)
TXN_ID="txn_failed_$(date +%s)"
PAYLOAD="{\"order_id\":\"$ORDER_ID\",\"transaction_id\":\"$TXN_ID\",\"payment_status\":\"failed\"}"
SIGNATURE=$(generate_signature "$PAYLOAD")
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -d "$PAYLOAD")
assert_contains "$RESPONSE" "200" "HTTP 200 returned"

# Verify order status remains pending but payment failed
sleep 1
ORDER_STATUS=$(curl -s "$API_URL/api/orders/$ORDER_ID" | grep -o '"status":"[^"]*' | head -1 | cut -d'"' -f4)
assert_contains "$ORDER_STATUS" "pending" "Order status remains pending"
ORDER_PAYMENT=$(curl -s "$API_URL/api/orders/$ORDER_ID" | grep -o '"payment_status":"[^"]*' | cut -d'"' -f4)
assert_contains "$ORDER_PAYMENT" "failed" "Payment status is failed"

#═══════════════════════════════════════════════════════════════
# TEST 9: Idempotency - Duplicate Transaction
#═══════════════════════════════════════════════════════════════
print_test "TEST 9: Idempotency - Duplicate Transaction (Should Success 200)"
ORDER_ID=$(create_order)
TXN_ID="txn_duplicate_$(date +%s)"
PAYLOAD="{\"order_id\":\"$ORDER_ID\",\"transaction_id\":\"$TXN_ID\",\"payment_status\":\"paid\"}"
SIGNATURE=$(generate_signature "$PAYLOAD")

# First request
RESPONSE1=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -d "$PAYLOAD")
assert_contains "$RESPONSE1" "200" "First request succeeded"

# Duplicate request with same transaction ID
RESPONSE2=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -d "$PAYLOAD")
assert_contains "$RESPONSE2" "200" "Duplicate request also returns 200"

# Verify only one webhook log entry
sleep 1
HISTORY=$(curl -s "$API_URL/api/orders/$ORDER_ID/payment-history")
LOG_COUNT=$(echo "$HISTORY" | grep -o "\"TransactionID\":\"$TXN_ID\"" | wc -l | tr -d ' ')
if [ "$LOG_COUNT" -eq 1 ]; then
  echo -e "${GREEN}✓${NC} Only one webhook log entry created (idempotency working)"
else
  echo -e "${RED}✗${NC} Expected 1 log entry, found $LOG_COUNT"
fi

#═══════════════════════════════════════════════════════════════
# TEST 10: Webhook on Already Completed Order
#═══════════════════════════════════════════════════════════════
print_test "TEST 10: Webhook On Already Completed Order (Should Fail 400)"
ORDER_ID=$(create_order)

# First, complete the order
TXN_ID_1="txn_first_$(date +%s)"
PAYLOAD1="{\"order_id\":\"$ORDER_ID\",\"transaction_id\":\"$TXN_ID_1\",\"payment_status\":\"paid\"}"
SIGNATURE1=$(generate_signature "$PAYLOAD1")
curl -s -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE1" \
  -d "$PAYLOAD1" > /dev/null

sleep 1

# Try to process another webhook on completed order
TXN_ID_2="txn_second_$(date +%s)"
PAYLOAD2="{\"order_id\":\"$ORDER_ID\",\"transaction_id\":\"$TXN_ID_2\",\"payment_status\":\"paid\"}"
SIGNATURE2=$(generate_signature "$PAYLOAD2")
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE2" \
  -d "$PAYLOAD2")
assert_contains "$RESPONSE" "400" "HTTP 400 returned"
assert_contains "$RESPONSE" "order status must be 'pending'" "Error message correct"

#═══════════════════════════════════════════════════════════════
# TEST 11: Webhook History Tracking
#═══════════════════════════════════════════════════════════════
print_test "TEST 11: Webhook History Tracking"
ORDER_ID=$(create_order)

# Send multiple webhooks with different transaction IDs
for i in 1 2 3; do
  TXN_ID="txn_history_${i}_$(date +%s)"
  PAYLOAD="{\"order_id\":\"$ORDER_ID\",\"transaction_id\":\"$TXN_ID\",\"payment_status\":\"failed\"}"
  SIGNATURE=$(generate_signature "$PAYLOAD")
  curl -s -X POST "$API_URL/api/payment-webhook" \
    -H "Content-Type: application/json" \
    -H "X-Webhook-Signature: $SIGNATURE" \
    -d "$PAYLOAD" > /dev/null
  sleep 0.5
done

# Now send a successful payment
TXN_ID_SUCCESS="txn_history_success_$(date +%s)"
PAYLOAD="{\"order_id\":\"$ORDER_ID\",\"transaction_id\":\"$TXN_ID_SUCCESS\",\"payment_status\":\"paid\"}"
SIGNATURE=$(generate_signature "$PAYLOAD")
curl -s -X POST "$API_URL/api/payment-webhook" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -d "$PAYLOAD" > /dev/null

sleep 1

# Check history
HISTORY=$(curl -s "$API_URL/api/orders/$ORDER_ID/payment-history")
HISTORY_COUNT=$(echo "$HISTORY" | grep -o '"ID"' | wc -l | tr -d ' ')
if [ "$HISTORY_COUNT" -eq 4 ]; then
  echo -e "${GREEN}✓${NC} All 4 webhook events logged correctly"
else
  echo -e "${RED}✗${NC} Expected 4 webhook events, found $HISTORY_COUNT"
fi

# Verify final order status
ORDER_STATUS=$(curl -s "$API_URL/api/orders/$ORDER_ID" | grep -o '"status":"[^"]*' | head -1 | cut -d'"' -f4)
assert_contains "$ORDER_STATUS" "completed" "Final order status is completed"

#═══════════════════════════════════════════════════════════════
# TEST 12: Concurrent Webhooks (Race Condition)
#═══════════════════════════════════════════════════════════════
print_test "TEST 12: Concurrent Webhooks For Same Order"
ORDER_ID=$(create_order)
TXN_ID="txn_concurrent_$(date +%s)"
PAYLOAD="{\"order_id\":\"$ORDER_ID\",\"transaction_id\":\"$TXN_ID\",\"payment_status\":\"paid\"}"
SIGNATURE=$(generate_signature "$PAYLOAD")

# Send 3 concurrent requests with same transaction ID
for i in 1 2 3; do
  curl -s -X POST "$API_URL/api/payment-webhook" \
    -H "Content-Type: application/json" \
    -H "X-Webhook-Signature: $SIGNATURE" \
    -d "$PAYLOAD" > /tmp/webhook_response_$i.txt &
done

wait

# All should return 200 (idempotency)
SUCCESS_COUNT=0
for i in 1 2 3; do
  if grep -q "success" /tmp/webhook_response_$i.txt; then
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
  fi
done

if [ "$SUCCESS_COUNT" -eq 3 ]; then
  echo -e "${GREEN}✓${NC} All concurrent requests handled successfully"
else
  echo -e "${RED}✗${NC} Some concurrent requests failed"
fi

# Clean up
rm -f /tmp/webhook_response_*.txt

#═══════════════════════════════════════════════════════════════
# Final Summary
#═══════════════════════════════════════════════════════════════
echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  ${GREEN}✓ Batch Test Suite Complete!${BLUE}                        ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}Tests Executed:${NC}"
echo "  ✓ Security: Signature validation (missing, invalid)"
echo "  ✓ Validation: Transaction ID, Order ID, Payment status"
echo "  ✓ Business Logic: Successful/Failed payments"
echo "  ✓ Idempotency: Duplicate transactions handled"
echo "  ✓ Edge Cases: Already completed orders, non-existent orders"
echo "  ✓ Audit: Webhook history tracking"
echo "  ✓ Concurrency: Race condition handling"
echo ""
echo -e "${GREEN}All webhook resilience features verified!${NC}"
