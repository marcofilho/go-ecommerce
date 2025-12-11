# Soft Deletes & Audit Logging

This document describes the implementation of soft deletes and audit logging features in the e-commerce API.

## Soft Deletes

### Overview

Soft deletes allow records to be marked as deleted without actually removing them from the database. This enables:
- Data recovery if needed
- Historical analysis
- Compliance with data retention policies
- Maintaining referential integrity

### Implementation

Soft deletes are implemented using GORM's `DeletedAt` field with automatic query scoping.

#### Entities with Soft Deletes

The following entities support soft deletes:

1. **Products** (`products` table)
2. **Product Variants** (`product_variants` table)
3. **Categories** (`categories` table)

#### Database Schema

Each table includes:
```sql
deleted_at timestamp with time zone
```

With an index for query performance:
```sql
CREATE INDEX idx_<table>_deleted_at ON <table> (deleted_at);
```

#### Entity Definition

```go
type Product struct {
    ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
    Name        string         `gorm:"size:255;not null"`
    Description string         `gorm:"type:text"`
    Price       float64        `gorm:"type:decimal(10,2);not null"`
    Quantity    int            `gorm:"not null"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"` // Soft delete field
    
    // Relations
    Variants   []ProductVariant `gorm:"foreignKey:ProductID"`
    Categories []Category       `gorm:"many2many:product_categories;"`
}
```

#### Behavior

**Automatic Filtering:**
GORM automatically excludes soft-deleted records from queries:

```go
// Only returns non-deleted products
products, err := repo.GetAll(ctx, page, pageSize, false)
```

**Delete Operation:**
```go
// Soft delete - sets deleted_at to current timestamp
err := repo.Delete(ctx, productID)
// Product still exists in DB with deleted_at set
```

**Include Deleted Records:**
```go
// Explicitly include soft-deleted records
db.Unscoped().Where("id = ?", id).Find(&product)
```

**Permanent Delete:**
```go
// Permanently delete from database
db.Unscoped().Delete(&product)
```

**Restore:**
```go
// Restore soft-deleted record
db.Model(&product).Update("deleted_at", nil)
```

### API Behavior

#### List Endpoints
Soft-deleted items are automatically excluded:

```bash
GET /api/products
# Returns only active products (deleted_at IS NULL)
```

#### Get by ID
Attempting to retrieve a soft-deleted item returns 404:

```bash
GET /api/products/{id}
# Returns 404 if product is soft-deleted
```

#### Delete
Delete operations perform soft delete:

```bash
DELETE /api/products/{id}
# Sets deleted_at timestamp, doesn't remove from DB
```

---

## Audit Logging

### Overview

Audit logging tracks all significant changes to critical resources, providing:
- Complete change history
- Before/after state comparison
- User attribution (when available)
- Compliance and security monitoring

### Implementation

#### Audit Logs Table

```sql
CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY,
    user_id         UUID,           -- Nullable for system actions
    action          VARCHAR(100) NOT NULL,
    resource_type   VARCHAR(100) NOT NULL,
    resource_id     UUID NOT NULL,
    payload_before  JSONB,          -- State before change
    payload_after   JSONB,          -- State after change
    timestamp       TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Indexes for efficient queries
CREATE INDEX idx_audit_logs_user_id ON audit_logs (user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs (action);
CREATE INDEX idx_audit_logs_resource_type ON audit_logs (resource_type);
CREATE INDEX idx_audit_logs_resource_id ON audit_logs (resource_id);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs (timestamp);
```

#### Entity Definition

```go
type AuditLog struct {
    ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
    UserID        *uuid.UUID     `gorm:"type:uuid;index"` // Nullable
    Action        string         `gorm:"size:100;not null;index"`
    ResourceType  string         `gorm:"size:100;not null;index"`
    ResourceID    uuid.UUID      `gorm:"type:uuid;not null;index"`
    PayloadBefore datatypes.JSON `gorm:"type:jsonb"`
    PayloadAfter  datatypes.JSON `gorm:"type:jsonb"`
    Timestamp     time.Time      `gorm:"not null;index"`
}
```

### Logged Actions

#### 1. Product Changes

**Actions Logged:**
- `CREATE` - New product created
- `UPDATE` - Product details modified
- `DELETE` - Product soft-deleted

**Example:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": null,
  "action": "UPDATE",
  "resource_type": "Product",
  "resource_id": "d4444444-4444-4444-4444-444444444444",
  "payload_before": {
    "id": "d4444444-4444-4444-4444-444444444444",
    "name": "MacBook Pro",
    "price": 2999.00,
    "quantity": 50
  },
  "payload_after": {
    "id": "d4444444-4444-4444-4444-444444444444",
    "name": "MacBook Pro 16\"",
    "price": 3499.00,
    "quantity": 25
  },
  "timestamp": "2025-12-11T03:30:00Z"
}
```

#### 2. Order Status Updates

**Actions Logged:**
- `UPDATE_STATUS` - Order status changed

**Example:**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "user_id": null,
  "action": "UPDATE_STATUS",
  "resource_type": "Order",
  "resource_id": "10101010-0101-0101-0101-010101010101",
  "payload_before": {
    "status": "pending"
  },
  "payload_after": {
    "status": "completed"
  },
  "timestamp": "2025-12-11T03:35:00Z"
}
```

#### 3. Payment Webhook Updates

**Actions Logged:**
- `PAYMENT_WEBHOOK` - Payment status updated via webhook

**Example:**
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "user_id": null,
  "action": "PAYMENT_WEBHOOK",
  "resource_type": "Order",
  "resource_id": "10101010-0101-0101-0101-010101010101",
  "payload_before": {
    "payment_status": "unpaid",
    "status": "pending"
  },
  "payload_after": {
    "payment_status": "paid",
    "status": "completed",
    "transaction_id": "TXN123456789"
  },
  "timestamp": "2025-12-11T03:40:00Z"
}
```

### Audit Service

#### Interface

```go
type AuditService interface {
    LogChange(
        ctx context.Context,
        userID *uuid.UUID,      // nil for system actions
        action string,          // e.g., "CREATE", "UPDATE", "DELETE"
        resourceType string,    // e.g., "Product", "Order"
        resourceID uuid.UUID,   // ID of the affected resource
        before interface{},     // State before change (nil for CREATE)
        after interface{}       // State after change (nil for DELETE)
    ) error
}
```

#### Usage Example

```go
// In product use case - Update operation
func (uc *UseCase) UpdateProduct(ctx context.Context, id uuid.UUID, ...) (*entity.Product, error) {
    product, err := uc.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Store original state
    original := *product
    
    // Apply updates
    product.Name = name
    product.Price = price
    // ... other updates
    
    if err := uc.repo.Update(ctx, product); err != nil {
        return nil, err
    }
    
    // Log the change
    uc.auditService.LogChange(ctx, nil, "UPDATE", "Product", product.ID, &original, product)
    
    return product, nil
}
```

### Query Audit Logs

#### Repository Methods

```go
// Get all audit logs with filters
logs, total, err := repo.List(ctx, filters, page, pageSize)

// Get audit logs for a specific resource
logs, err := repo.GetByResourceID(ctx, "Product", productID)
```

#### Filter Options

```go
type AuditLogFilters struct {
    UserID       *uuid.UUID
    Action       *string
    ResourceType *string
    ResourceID   *uuid.UUID
    StartDate    *string  // RFC3339 format
    EndDate      *string  // RFC3339 format
}
```

### Best Practices

#### 1. User Attribution
- Pass `userID` from authenticated context when available
- Use `nil` for system/automated actions (webhooks, background jobs)

#### 2. Payload Content
- Include only relevant fields in before/after payloads
- Avoid sensitive data (passwords, tokens)
- Keep payloads concise but informative

#### 3. Action Naming
- Use consistent, descriptive action names
- Follow convention: `CREATE`, `UPDATE`, `DELETE`, `UPDATE_STATUS`, etc.
- Use uppercase for consistency

#### 4. Performance
- Audit logging is non-blocking
- Failures don't prevent the main operation
- Indexes optimize audit log queries

### Compliance & Security

#### Data Retention
- Audit logs are retained indefinitely by default
- Implement retention policies based on compliance requirements
- Consider archiving old logs to separate storage

#### Access Control
- Audit log access should be restricted to administrators
- Implement separate endpoints with proper authorization
- Log access to audit logs themselves for security

#### Immutability
- Audit log entries should never be modified
- Only INSERT operations allowed
- DELETE only with proper authorization and logging

### Monitoring & Alerts

#### Key Metrics to Monitor
- Audit log creation rate
- Failed audit log writes
- Unusual patterns in logged actions
- Suspicious user activities

#### Example Queries

**Most active users:**
```sql
SELECT user_id, COUNT(*) as action_count
FROM audit_logs
WHERE timestamp > NOW() - INTERVAL '24 hours'
GROUP BY user_id
ORDER BY action_count DESC
LIMIT 10;
```

**Recent deletions:**
```sql
SELECT resource_type, resource_id, payload_before, timestamp
FROM audit_logs
WHERE action = 'DELETE'
ORDER BY timestamp DESC
LIMIT 50;
```

**Changes to specific product:**
```sql
SELECT action, payload_before, payload_after, timestamp
FROM audit_logs
WHERE resource_type = 'Product'
  AND resource_id = 'd4444444-4444-4444-4444-444444444444'
ORDER BY timestamp DESC;
```

---

## Testing

Run the test script to verify implementation:

```bash
./test_soft_deletes_audit.sh
```

This validates:
- ✅ Soft delete columns exist on all entities
- ✅ Indexes created for performance
- ✅ Audit logs table properly structured
- ✅ GORM automatically filters deleted records
- ✅ Audit service integrated in use cases
