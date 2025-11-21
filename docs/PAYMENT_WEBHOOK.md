# Payment Webhook Integration

## Overview

The payment webhook endpoint receives payment status updates from the payment processor. This implementation includes HMAC signature verification for security and transaction ID-based idempotency for reliability.

## Security Features

### HMAC Signature Verification

All webhook requests must include a valid HMAC-SHA256 signature in the `X-Webhook-Signature` header. The signature is computed as:

```
signature = HMAC-SHA256(webhook_secret, request_body)
```

**Example (bash):**
```bash
WEBHOOK_SECRET="your-webhook-secret-key"
PAYLOAD='{"order_id":"123e4567-e89b-12d3-a456-426614174000","transaction_id":"txn_12345","payment_status":"paid"}'
SIGNATURE=$(echo -n "$PAYLOAD" | openssl dgst -sha256 -hmac "$WEBHOOK_SECRET" | sed 's/^.* //')

curl -X POST http://localhost:8080/api/payment-webhook \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -d "$PAYLOAD"
```

### Configuration

Set the webhook secret in your environment:

```bash
# docker-compose.yml
WEBHOOK_SECRET=my-super-secret-webhook-key-change-in-production
```

**Important:** Change the default secret in production!

## Endpoint

**POST** `/api/payment-webhook`

### Headers

| Header | Required | Description |
|--------|----------|-------------|
| `Content-Type` | Yes | Must be `application/json` |
| `X-Webhook-Signature` | Yes | HMAC-SHA256 signature of request body |

### Request Body

```json
{
  "order_id": "123e4567-e89b-12d3-a456-426614174000",
  "transaction_id": "txn_unique_12345",
  "payment_status": "paid"
}
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `order_id` | string (UUID) | Yes | The order identifier |
| `transaction_id` | string | Yes | Unique transaction identifier for idempotency |
| `payment_status` | string | Yes | Payment status: `"paid"` or `"failed"` |

### Response

**Success (200 OK):**
```json
{
  "status": "success",
  "message": "Payment webhook processed successfully"
}
```

**Error (400 Bad Request):**
```json
{
  "error": "Invalid request body"
}
```

**Error (401 Unauthorized):**
```json
{
  "error": "Invalid webhook signature"
}
```

## Resilience Features

### 1. Idempotency

The system uses `transaction_id` to ensure idempotent processing. If the same `transaction_id` is received multiple times, subsequent requests return success without reprocessing.

### 2. Retry Support

The webhook logs track processing status and retry information:

- **Status**: `pending` → `processing` → `completed` or `failed`
- **Retry Count**: Incremented on failures
- **Next Retry**: Scheduled 5 minutes after failure

The system only returns HTTP 200 after successful database commit, allowing payment processors to retry on failure.

### 3. Audit Trail

All webhook events are logged in the `webhook_logs` table with:
- Transaction ID
- Payment status
- Processing status
- Retry count
- Raw payload
- Timestamps

## Validation Rules

1. **Signature Verification**: Request must have valid HMAC signature
2. **Transaction ID**: Must be present and unique
3. **Order ID**: Must be a valid UUID format
4. **Order Exists**: Order must exist in the database
5. **Order Status**: Order must be in `pending` status
6. **Payment Status**: Must be either `"paid"` or `"failed"`

## Behavior

### Successful Payment (`"paid"`)
- Order status: `pending` → `completed`
- Payment status: `unpaid` → `paid`
- Webhook log: Status set to `completed`

### Failed Payment (`"failed"`)
- Order status: Remains `pending` (customer can retry)
- Payment status: `unpaid` → `failed`
- Webhook log: Status set to `completed`

### Processing Error
- Webhook log: Status set to `failed`
- Retry count: Incremented
- Next retry: Scheduled for 5 minutes later
- HTTP response: 400/500 (processor will retry)

## Payment History

**GET** `/api/orders/{id}/payment-history`

Returns all webhook events for an order.

## Testing

Use the provided test script:

```bash
./test_payment_webhook.sh
```

Tests:
1. Product and order creation
2. Webhook without signature (should fail)
3. Webhook with valid signature
4. Order status updates
5. Idempotency with duplicate transactions
6. Payment history retrieval

## Error Handling

| Error | HTTP Code | Description |
|-------|-----------|-------------|
| Missing signature | 401 | `X-Webhook-Signature` header not present |
| Invalid signature | 401 | HMAC signature verification failed |
| Missing transaction_id | 400 | `transaction_id` field is required |
| Invalid request body | 400 | JSON parsing failed |
| Invalid order_id | 400 | Order ID is not a valid UUID |
| Order not found | 400 | Order does not exist |
| Invalid order status | 400 | Order is not in pending status |
| Invalid payment_status | 400 | Must be "paid" or "failed" |
| Database error | 400 | Failed to update order or create log |
