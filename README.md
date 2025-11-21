# Go E-Commerce API

A RESTful API for managing products and orders in an e-commerce system, built with Go using clean architecture principles and PostgreSQL.

## Features

- Product Management (CRUD with stock tracking)
- Order Management (create orders with automatic stock deduction)
- Status workflow (pending → completed/canceled)
- Pagination & filtering
- PostgreSQL with GORM ORM
- Automatic migrations

## Quick Start

### Prerequisites

- Go 1.24.1+
- Docker and Docker Compose

### Setup

```bash
# Start PostgreSQL
make db-up

# Run the server
make run
```

Server starts at `http://localhost:8080`

### Test API

```bash
# Create a product
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"High-performance","price":999.99,"quantity":50}'

# List products
curl http://localhost:8080/api/products

# Create an order
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_id":123,"products":[{"product_id":"YOUR_PRODUCT_ID","quantity":2}]}'
```

## API Endpoints

### Products

- `POST /api/products` - Create product
- `GET /api/products` - List products (supports `?page=1&page_size=10&in_stock_only=true`)
- `GET /api/products/{id}` - Get product
- `PUT /api/products/{id}` - Update product
- `DELETE /api/products/{id}` - Delete product

### Orders

- `POST /api/orders` - Create order
- `GET /api/orders` - List orders (supports `?page=1&page_size=10&status=pending`)
- `GET /api/orders/{id}` - Get order
- `PUT /api/orders/{id}/status` - Update order status

## Architecture

```
src/
├── cmd/api/              # Entry point (main, container, routes)
├── internal/
│   ├── domain/           # Entities & repository interfaces
│   ├── infrastructure/   # Repository implementations (PostgreSQL)
│   └── adapter/http/     # HTTP handlers & DTOs
└── usecase/              # Business logic
```

## Make Commands

```bash
make db-up      # Start PostgreSQL
make db-down    # Stop PostgreSQL
make db-reset   # Reset database
make run        # Run application
make build      # Build binary
make clean      # Clean artifacts
```

## Configuration

Environment variables (defaults):

- `DB_HOST=localhost`
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=postgres`
- `DB_NAME=ecommerce`
- `SERVER_PORT=8080`

## License

MIT
