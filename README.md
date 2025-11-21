# Go E-Commerce API

A RESTful API for managing products and orders in an e-commerce system, built with Go using clean architecture principles and PostgreSQL.

## Features

- Product Management (CRUD with stock tracking)
- Order Management (create orders with automatic stock deduction)
- Status workflow (pending → completed/canceled)
- Pagination & filtering
- PostgreSQL with GORM ORM
- Automatic migrations
- **Swagger/OpenAPI documentation** - Interactive API testing at `/swagger/`

## Quick Start

### Prerequisites

- Docker and Docker Compose

### Setup

**Option 1: Full Stack (Recommended)**

```bash
# Build and start everything (PostgreSQL + API)
make docker-up

# View logs
make docker-logs

# Stop everything
make docker-down
```

**Option 2: Local Development**

```bash
# Start only PostgreSQL
make db-up

# Run API locally (requires Go 1.24.1+)
make run
```

Server starts at `http://localhost:8080`

**Swagger UI:** `http://localhost:8080/swagger/`

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

**Docker (Full Stack):**

```bash
make docker-up       # Start all services (PostgreSQL + API)
make docker-down     # Stop all services
make docker-logs     # View logs
make docker-restart  # Restart all services
make docker-build    # Rebuild Docker image
```

**Local Development:**

```bash
make db-up      # Start PostgreSQL only
make db-down    # Stop PostgreSQL only
make db-reset   # Reset database
make run        # Run API locally
make build      # Build binary
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
