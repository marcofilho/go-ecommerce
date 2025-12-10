#!/bin/bash

# Payment Webhook Batch Test Script
# Tests various success and failure scenarios with enhanced security

set -e

API_URL="http://localhost:8080"
WEBHOOK_SECRET="my-super-secret-webhook-key-change-in-production"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

generate_signature() {
  local payload="$1"
  echo -n "$payload" | openssl dgst -sha256 -hmac "$WEBHOOK_SECRET" | sed 's/^.* //'
}

get_timestamp() {
  date +%s
}

print_test() {
  echo ""
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${YELLOW}$1${NC}"
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

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
echo -e "${BLUE}║  Payment Webhook Enhanced Security Test Suite         ║${NC}"
echo -e "${BLUE}║  Testing Signature & Timestamp Validation             ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"

# Setup admin and customer users (same as original script)
# ... (keeping the setup code from original)
