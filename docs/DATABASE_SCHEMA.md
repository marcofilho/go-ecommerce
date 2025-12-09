# Database Schema

## Overview

This document describes the complete database schema for the Go E-Commerce API. The system uses PostgreSQL with GORM for ORM operations and automatic migrations.

## Entity Relationship Diagram

```
┌─────────────┐
│    users    │
└──────┬──────┘
       │
       │ 1:N (customer_id)
       │
       ▼
┌─────────────┐         ┌──────────────────┐
│   orders    │────────▶│   order_items    │
└─────────────┘   1:N   └────────┬─────────┘
       │                         │
       │ 1:N                     │ N:1
       │                         │
       ▼                         ▼
┌─────────────┐         ┌──────────────────┐
│webhook_logs │         │    products      │
└─────────────┘         └────────┬─────────┘
                                 │
                    ┌────────────┼────────────┐
                    │            │            │
                    │ 1:N        │ N:M        │
                    ▼            │            ▼
            ┌───────────────┐    │    ┌─────────────┐
            │product_variants│   │    │ categories  │
            └───────────────┘    │    └─────────────┘
                                 │            ▲
                                 │            │
                                 │    N:M     │
                                 │  (junction)│
                                 ▼            │
                         ┌────────────────────┘
                         │ product_categories │
                         └────────────────────┘
```

## Tables

### 1. users

Stores user accounts with authentication and role-based access control.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | SERIAL | PRIMARY KEY | Auto-incrementing user ID |
| email | VARCHAR(255) | UNIQUE, NOT NULL | User email (login) |
| password_hash | VARCHAR(255) | NOT NULL | Bcrypt hashed password |
| name | VARCHAR(255) | NOT NULL | User's full name |
| role | VARCHAR(50) | NOT NULL, DEFAULT 'customer' | Role: 'admin' or 'customer' |
| created_at | TIMESTAMP | NOT NULL | Account creation timestamp |
| updated_at | TIMESTAMP | NOT NULL | Last update timestamp |

**Indexes:**
- PRIMARY KEY on `id`
- UNIQUE INDEX on `email`

**Example:**
```sql
id | email                | password_hash      | name       | role     | created_at          | updated_at
---+----------------------+--------------------+------------+----------+---------------------+---------------------
1  | admin@example.com    | $2a$10$...       | Admin User | admin    | 2025-12-08 10:00:00 | 2025-12-08 10:00:00
2  | customer@example.com | $2a$10$...       | John Doe   | customer | 2025-12-08 10:05:00 | 2025-12-08 10:05:00
```

---

### 2. categories

Stores product categories for organizing products.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Category unique identifier |
| name | VARCHAR(255) | UNIQUE, NOT NULL | Category name |
| created_at | TIMESTAMP | NOT NULL | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL | Last update timestamp |

**Indexes:**
- PRIMARY KEY on `id`
- UNIQUE INDEX on `name`

**Example:**
```sql
id                                   | name        | created_at          | updated_at
-------------------------------------+-------------+---------------------+---------------------
550e8400-e29b-41d4-a716-446655440000 | Electronics | 2025-12-08 10:00:00 | 2025-12-08 10:00:00
660e8400-e29b-41d4-a716-446655440001 | Clothing    | 2025-12-08 10:01:00 | 2025-12-08 10:01:00
770e8400-e29b-41d4-a716-446655440002 | Books       | 2025-12-08 10:02:00 | 2025-12-08 10:02:00
```

---

### 3. products

Stores base product information with stock management.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Product unique identifier |
| name | VARCHAR(255) | NOT NULL | Product name |
| description | TEXT | | Product description |
| price | DECIMAL(10,2) | NOT NULL, CHECK (price >= 0) | Base product price |
| quantity | INTEGER | NOT NULL, CHECK (quantity >= 0) | Stock quantity |
| created_at | TIMESTAMP | NOT NULL | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL | Last update timestamp |

**Indexes:**
- PRIMARY KEY on `id`
- INDEX on `name` for search optimization

**Business Rules:**
- Price must be non-negative
- Quantity must be non-negative
- Stock is automatically deducted when orders are created

**Example:**
```sql
id                                   | name   | description          | price  | quantity | created_at          | updated_at
-------------------------------------+--------+----------------------+--------+----------+---------------------+---------------------
550e8400-e29b-41d4-a716-446655440000 | Laptop | High-performance     | 999.99 | 50       | 2025-12-08 10:00:00 | 2025-12-08 10:00:00
660e8400-e29b-41d4-a716-446655440001 | Mouse  | Wireless gaming      | 49.99  | 100      | 2025-12-08 10:05:00 | 2025-12-08 10:05:00
```

---

### 4. product_variants

Stores product variants (e.g., colors, sizes) with optional price overrides.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Variant unique identifier |
| product_id | UUID | NOT NULL, FOREIGN KEY → products(id) | Parent product reference |
| name | VARCHAR(255) | NOT NULL | Variant name (e.g., "Color") |
| value | VARCHAR(255) | NOT NULL | Variant value (e.g., "Red") |
| quantity | INTEGER | NOT NULL, CHECK (quantity >= 0) | Variant stock quantity |
| price_override | DECIMAL(10,2) | CHECK (price_override >= 0) | Optional variant price (NULL = use product price) |
| created_at | TIMESTAMP | NOT NULL | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL | Last update timestamp |

**Indexes:**
- PRIMARY KEY on `id`
- FOREIGN KEY INDEX on `product_id`
- COMPOSITE INDEX on `(product_id, name, value)` for uniqueness

**Business Rules:**
- If `price_override` is NULL, the variant uses the parent product's price
- If `price_override` is set, it overrides the product price
- Stock is managed independently for each variant

**Example:**
```sql
id                                   | product_id                           | name  | value | quantity | price_override | created_at          | updated_at
-------------------------------------+--------------------------------------+-------+-------+----------+----------------+---------------------+---------------------
770e8400-e29b-41d4-a716-446655440000 | 550e8400-e29b-41d4-a716-446655440000 | Color | Red   | 20       | 1049.99        | 2025-12-08 10:10:00 | 2025-12-08 10:10:00
880e8400-e29b-41d4-a716-446655440001 | 550e8400-e29b-41d4-a716-446655440000 | Color | Blue  | 30       | NULL           | 2025-12-08 10:11:00 | 2025-12-08 10:11:00
```

---

### 5. product_categories (Junction Table)

N:N relationship between products and categories using composite primary key.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| product_id | UUID | PRIMARY KEY, FOREIGN KEY → products(id) ON DELETE CASCADE | Product reference |
| category_id | UUID | PRIMARY KEY, FOREIGN KEY → categories(id) ON DELETE CASCADE | Category reference |

**Indexes:**
- COMPOSITE PRIMARY KEY on `(product_id, category_id)`
- FOREIGN KEY INDEX on `product_id`
- FOREIGN KEY INDEX on `category_id`

**Business Rules:**
- Composite primary key prevents duplicate assignments
- CASCADE DELETE removes associations when product or category is deleted
- A product can have multiple categories
- A category can be assigned to multiple products

**Example:**
```sql
product_id                           | category_id
-------------------------------------+-------------------------------------
550e8400-e29b-41d4-a716-446655440000 | 550e8400-e29b-41d4-a716-446655440000  (Laptop → Electronics)
550e8400-e29b-41d4-a716-446655440000 | 770e8400-e29b-41d4-a716-446655440002  (Laptop → Computers)
660e8400-e29b-41d4-a716-446655440001 | 550e8400-e29b-41d4-a716-446655440000  (Mouse → Electronics)
```

---

### 6. orders

Stores order information with payment and status tracking.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Order unique identifier |
| customer_id | INTEGER | NOT NULL, FOREIGN KEY → users(id) | Customer reference |
| status | VARCHAR(50) | NOT NULL, DEFAULT 'pending' | Order status |
| payment_status | VARCHAR(50) | NOT NULL, DEFAULT 'unpaid' | Payment status |
| total | DECIMAL(10,2) | NOT NULL, CHECK (total >= 0) | Order total amount |
| created_at | TIMESTAMP | NOT NULL | Order creation timestamp |
| updated_at | TIMESTAMP | NOT NULL | Last update timestamp |

**Allowed Values:**
- `status`: 'pending', 'completed', 'canceled'
- `payment_status`: 'unpaid', 'paid', 'failed'

**Indexes:**
- PRIMARY KEY on `id`
- FOREIGN KEY INDEX on `customer_id`
- INDEX on `status` for filtering
- INDEX on `payment_status` for filtering

**Business Rules:**
- Total is calculated from order items (quantity × price)
- Status can transition: pending → completed/canceled
- Payment status can transition: unpaid → paid/failed
- Cannot modify order after completion

**Example:**
```sql
id                                   | customer_id | status  | payment_status | total   | created_at          | updated_at
-------------------------------------+-------------+---------+----------------+---------+---------------------+---------------------
990e8400-e29b-41d4-a716-446655440000 | 2           | pending | unpaid         | 1999.98 | 2025-12-08 11:00:00 | 2025-12-08 11:00:00
aa0e8400-e29b-41d4-a716-446655440001 | 2           | completed| paid          | 49.99   | 2025-12-08 11:05:00 | 2025-12-08 11:10:00
```

---

### 7. order_items

Stores individual items within an order with optional variant support.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Order item unique identifier |
| order_id | UUID | NOT NULL, FOREIGN KEY → orders(id) ON DELETE CASCADE | Order reference |
| product_id | UUID | NOT NULL, FOREIGN KEY → products(id) | Product reference |
| variant_id | UUID | FOREIGN KEY → product_variants(id) | Optional variant reference |
| quantity | INTEGER | NOT NULL, CHECK (quantity > 0) | Item quantity |
| price | DECIMAL(10,2) | NOT NULL, CHECK (price >= 0) | Price at purchase time |
| created_at | TIMESTAMP | NOT NULL | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL | Last update timestamp |

**Indexes:**
- PRIMARY KEY on `id`
- FOREIGN KEY INDEX on `order_id`
- FOREIGN KEY INDEX on `product_id`
- FOREIGN KEY INDEX on `variant_id`

**Business Rules:**
- Price is captured at order time (historical price)
- If `variant_id` is set, uses variant's price (or product price if no override)
- If `variant_id` is NULL, uses product's base price
- Stock is deducted from product or variant when order is created

**Example:**
```sql
id                                   | order_id                             | product_id                           | variant_id                           | quantity | price   | created_at          | updated_at
-------------------------------------+--------------------------------------+--------------------------------------+--------------------------------------+----------+---------+---------------------+---------------------
bb0e8400-e29b-41d4-a716-446655440000 | 990e8400-e29b-41d4-a716-446655440000 | 550e8400-e29b-41d4-a716-446655440000 | 770e8400-e29b-41d4-a716-446655440000 | 2        | 1049.99 | 2025-12-08 11:00:00 | 2025-12-08 11:00:00
cc0e8400-e29b-41d4-a716-446655440001 | aa0e8400-e29b-41d4-a716-446655440001 | 660e8400-e29b-41d4-a716-446655440001 | NULL                                 | 1        | 49.99   | 2025-12-08 11:05:00 | 2025-12-08 11:05:00
```

---

### 8. webhook_logs

Audit trail for payment webhook events.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Log entry unique identifier |
| order_id | VARCHAR(255) | NOT NULL, FOREIGN KEY → orders(id) | Order reference |
| transaction_id | VARCHAR(255) | UNIQUE, NOT NULL | Payment processor transaction ID |
| status | VARCHAR(50) | NOT NULL | Payment status received |
| payload | JSONB | NOT NULL | Complete webhook payload |
| created_at | TIMESTAMP | NOT NULL | Webhook receipt timestamp |

**Indexes:**
- PRIMARY KEY on `id`
- UNIQUE INDEX on `transaction_id` (idempotency key)
- FOREIGN KEY INDEX on `order_id`
- INDEX on `created_at` for chronological queries

**Business Rules:**
- `transaction_id` ensures idempotent webhook processing
- JSONB payload allows flexible querying of webhook data
- Used for audit trail and compliance
- Enables replay and debugging of payment events

**Example:**
```sql
id                                   | order_id                             | transaction_id | status | payload                      | created_at
-------------------------------------+--------------------------------------+----------------+--------+------------------------------+---------------------
dd0e8400-e29b-41d4-a716-446655440000 | aa0e8400-e29b-41d4-a716-446655440001 | txn_12345      | paid   | {"amount":49.99,"method":"card"} | 2025-12-08 11:10:00
```

---

## Migration Order

Tables must be created in this order to satisfy foreign key constraints:

1. `users` - No dependencies
2. `categories` - No dependencies
3. `products` - No dependencies
4. `product_variants` - Depends on `products`
5. `product_categories` - Depends on `products` and `categories`
6. `orders` - Depends on `users`
7. `order_items` - Depends on `orders`, `products`, and `product_variants`
8. `webhook_logs` - Depends on `orders`

## Automatic Migrations

The application uses GORM AutoMigrate to automatically create and update tables:

```go
// src/internal/infrastructure/database/database.go
func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &entity.User{},
        &entity.Category{},
        &entity.Product{},
        &entity.ProductVariant{},
        &entity.ProductCategory{},
        &entity.Order{},
        &entity.OrderItem{},
        &entity.WebhookLog{},
    )
}
```

**To run migrations manually:**
```bash
go run src/cmd/migrate/main.go
```

## Manual Database Operations

### Create Admin User

```sql
-- After registering a regular user, promote to admin:
UPDATE users 
SET role = 'admin' 
WHERE email = 'admin@example.com';
```

### View Product with Categories

```sql
SELECT 
    p.id,
    p.name,
    p.price,
    STRING_AGG(c.name, ', ') as categories
FROM products p
LEFT JOIN product_categories pc ON p.id = pc.product_id
LEFT JOIN categories c ON pc.category_id = c.id
GROUP BY p.id, p.name, p.price;
```

### View Order Details with Items

```sql
SELECT 
    o.id as order_id,
    o.status,
    o.payment_status,
    o.total,
    p.name as product_name,
    pv.value as variant,
    oi.quantity,
    oi.price
FROM orders o
JOIN order_items oi ON o.id = oi.order_id
JOIN products p ON oi.product_id = p.id
LEFT JOIN product_variants pv ON oi.variant_id = pv.id
WHERE o.id = 'order-uuid-here';
```

### Reset Database

```bash
# Using Docker Compose
docker-compose down -v
docker-compose up -d

# The application will automatically run migrations on startup
```

## Database Connection

Configuration via environment variables:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ecommerce
```

Connection string format:
```
host=localhost port=5432 user=postgres password=postgres dbname=ecommerce sslmode=disable
```

## Performance Considerations

### Indexes

The schema includes indexes on:
- All primary keys (automatic)
- All foreign keys (automatic with GORM)
- Unique constraints (`users.email`, `categories.name`, `webhook_logs.transaction_id`)
- Filter columns (`orders.status`, `orders.payment_status`)
- Search columns (`products.name`)

### Query Optimization

- Use `Preload()` for eager loading relationships
- Apply pagination to all list endpoints
- Use composite indexes for common query patterns

### Scaling Considerations

- UUID primary keys enable distributed ID generation
- JSONB in `webhook_logs` allows flexible querying without schema changes
- Soft deletes not implemented (use status flags if needed)
- Consider read replicas for high-traffic deployments

## Data Integrity

### Constraints

- **CHECK constraints**: Enforce business rules (non-negative prices, positive quantities)
- **FOREIGN KEY constraints**: Maintain referential integrity
- **UNIQUE constraints**: Prevent duplicates (emails, transaction IDs, category names)
- **CASCADE DELETE**: Automatic cleanup of dependent records

### Validation

- Application-level validation in entity `Validate()` methods
- Database-level constraints as safety net
- Transactions ensure atomic operations

## Backup and Recovery

### Recommended Backup Strategy

```bash
# Daily automated backups
pg_dump -U postgres ecommerce > backup_$(date +%Y%m%d).sql

# Restore from backup
psql -U postgres ecommerce < backup_20251208.sql
```

### Point-in-Time Recovery

Enable WAL archiving in PostgreSQL for continuous backups and point-in-time recovery capabilities.

## Security

- Passwords stored as bcrypt hashes (never plaintext)
- No sensitive data in logs
- JWT secrets stored in environment variables
- Webhook signatures verified before processing
- SQL injection prevented by GORM parameterized queries

---

**Last Updated:** December 8, 2025
