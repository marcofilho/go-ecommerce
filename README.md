# Go E-Commerce API

A RESTful API for managing products and orders in an e-commerce system, built with Go using clean architecture principles and PostgreSQL database.

## Features

- **Product Management**: CRUD operations with stock tracking
- **Order Management**: Create and manage orders with automatic stock deduction
- **Status Workflow**: Controlled order status transitions (pending → completed/canceled)
- **Pagination & Filtering**: List endpoints with pagination and filtering support
- **PostgreSQL Database**: Persistent storage with GORM ORM
- **Automatic Migrations**: Database schema managed automatically

## Quick Start

### Prerequisites

- Go 1.24.1 or higher
- Docker and Docker Compose (for PostgreSQL)

### 1. Setup Database

**Make sure Docker is running first!**

```bash
# Start PostgreSQL using Docker Compose
make db-up

# Or manually with docker-compose
docker-compose up -d
```

**Alternative: Use existing PostgreSQL**

If you have PostgreSQL installed locally or want to use a remote instance, just set the environment variables:

```bash
export DB_HOST=your-host
export DB_PORT=5432
export DB_USER=your-user
export DB_PASSWORD=your-password
export DB_NAME=ecommerce
```

### 2. Configure Environment (Optional)

Copy `.env.example` to `.env` and adjust if needed:

```bash
cp .env.example .env
```

Default configuration:

- Database: `localhost:5432`
- User/Password: `postgres/postgres`
- Database name: `ecommerce`

### 3. Run the Server

```bash
# Build and run (will auto-start database if not running)
make run

# Or run directly
go run ./src/cmd/api/main.go
```

The server will:

1. Connect to PostgreSQL
2. Run migrations automatically (creates tables)
3. Start on `http://localhost:8080`

### Test the API

```bash
# Create a product
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "quantity": 50
  }'

# List products
curl http://localhost:8080/api/products

# Create an order (use product ID from previous response)
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 123,
    "products": [
      {
        "product_id": "YOUR_PRODUCT_ID",
        "quantity": 2
      }
    ]
  }'
```

Run automated tests:

```bash
./test_api.sh
```

## API Endpoints

### Products

| Method | Endpoint             | Description                               |
| ------ | -------------------- | ----------------------------------------- |
| POST   | `/api/products`      | Create product                            |
| GET    | `/api/products`      | List products (with pagination & filters) |
| GET    | `/api/products/{id}` | Get product by ID                         |
| PUT    | `/api/products/{id}` | Update product                            |
| DELETE | `/api/products/{id}` | Delete product                            |

**Query Parameters for List:**

- `page` - Page number (default: 1)
- `page_size` - Items per page (default: 10, max: 100)
- `in_stock_only` - Filter products with quantity > 0

### Orders

| Method | Endpoint                  | Description                             |
| ------ | ------------------------- | --------------------------------------- |
| POST   | `/api/orders`             | Create order                            |
| GET    | `/api/orders`             | List orders (with pagination & filters) |
| GET    | `/api/orders/{id}`        | Get order by ID                         |
| PUT    | `/api/orders/{id}/status` | Update order status                     |

**Query Parameters for List:**

- `page` - Page number (default: 1)
- `page_size` - Items per page (default: 10, max: 100)
- `status` - Filter by status (pending, cancelled, completed)
- `payment_status` - Filter by payment status (unpaid, paid, failed)

**Valid Status Transitions:**

- `pending` → `completed`
- `pending` → `cancelled`

## Request/Response Examples

### Create Product

```json
// Request
{
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "quantity": 50
}

// Response (201 Created)
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "quantity": 50,
  "created_at": "2025-11-20T10:00:00Z",
  "updated_at": "2025-11-20T10:00:00Z"
}
```

### Create Order

```json
// Request
{
  "customer_id": 123,
  "products": [
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440000",
      "quantity": 2
    }
  ]
}

// Response (201 Created)
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "customer_id": 123,
  "items": [
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440000",
      "quantity": 2,
      "price": 999.99,
      "subtotal": 1999.98
    }
  ],
  "total_price": 1999.98,
  "status": "pending",
  "payment_status": "unpaid",
  "created_at": "2025-11-20T10:05:00Z",
  "updated_at": "2025-11-20T10:05:00Z"
}
```

## Architecture

```
src/
├── cmd/api/              # Application entry point
├── internal/
│   ├── domain/
│   │   ├── entity/       # Business entities
│   │   └── repository/   # Repository interfaces
│   ├── infrastructure/
│   │   └── repository/   # Repository implementations
│   └── adapter/http/     # HTTP handlers
└── usecase/             # Business logic
    ├── product/
    └── order/
```

## Make Commands

```bash
make db-up      # Start PostgreSQL database
make db-down    # Stop PostgreSQL database
make db-reset   # Reset database (fresh start)
make build      # Build the application
make run        # Start database and run application
make test       # Run tests
make test-api   # Test API endpoints (server must be running)
make clean      # Clean build artifacts
make fmt        # Format code
make help       # Show all commands
```

## Database Schema

The application automatically creates the following tables:

### products

- `id` (UUID, Primary Key)
- `name` (VARCHAR, NOT NULL)
- `description` (TEXT)
- `price` (DECIMAL)
- `quantity` (INT)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### orders

- `id` (UUID, Primary Key)
- `customer_id` (INT, NOT NULL)
- `total_price` (DECIMAL)
- `status` (VARCHAR: pending/cancelled/completed)
- `payment_status` (VARCHAR: unpaid/paid/failed)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### order_items

- `id` (INT, Primary Key)
- `order_id` (UUID, Foreign Key)
- `product_id` (UUID)
- `quantity` (INT)
- `price` (DECIMAL)

## Business Rules

1. **Products**: Name required, price/quantity cannot be negative
2. **Orders**: Customer ID required, minimum one item, automatic stock deduction
3. **Status Transitions**: Only pending orders can be completed or cancelled

## Error Responses

```json
{
  "error": "Error message describing what went wrong"
}
```

Common status codes: `400` (Bad Request), `404` (Not Found), `500` (Internal Server Error)

## License

MIT
