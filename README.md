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
- (Optional) Go 1.24.1+ for local development

### Setup

```bash
# Start everything (PostgreSQL + API)
make up

# View logs
make logs

# Stop everything
make down
```

Server starts at `http://localhost:8080`

**Swagger UI:** `http://localhost:8080/swagger/index.html`

### Test API

**Via Swagger UI (Recommended):**
Visit `http://localhost:8080/swagger/index.html` for interactive API testing

**Via curl:**

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

## Testing

Run all unit tests:

```bash
# Run tests in Docker
make test
```

**Test Coverage:**

- **Domain entities: 100.0% coverage** ✅
- **DTO mappers: 100.0% coverage** ✅
- **HTTP handlers: 100.0% coverage** ✅
- **Product use cases: 100.0% coverage** ✅
- **Order use cases: 95.1% coverage** ✅
- **Total: 77 passing tests across 7 test files**

**Test Suites:**

- Entity layer: Product & Order business logic validation, GORM hooks
- DTO layer: Request/Response mapping and pagination
- Handler layer: HTTP request/response handling, validation, error responses
- Use case layer: Product & Order CRUD operations with comprehensive error handling
- All edge cases: Invalid inputs, repository errors, validation failures, pagination defaults

## Architecture

```
src/
├── cmd/api/              # Entry point (main, container, routes)
├── internal/
│   ├── domain/           # Entities & repository interfaces
│   ├── infrastructure/   # Repository implementations (PostgreSQL)
│   ├── adapter/http/
│   │   ├── handler/      # HTTP handlers
│   │   └── dto/          # Data Transfer Objects
│   └── config/           # Configuration
└── usecase/              # Business logic
```

## Make Commands

```bash
make up     # Start all services (PostgreSQL + API)
make down   # Stop all services
make logs   # View service logs
make test   # Run tests in Docker
make help   # Show available commands
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
