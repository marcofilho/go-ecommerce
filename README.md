# Go E-Commerce API

A RESTful API for managing products and orders in an e-commerce system, built with Go using clean architecture principles and PostgreSQL.

## Features

- **Authentication & Authorization** (JWT-based with RBAC)
- **Role-Based Permissions** (admin vs customer access control)
- Product Management (CRUD with stock tracking)
- Order Management (create orders with automatic stock deduction)
- **Payment Webhook Integration** (simulated payment gateway)
- Payment status tracking (unpaid → paid/failed)
- Status workflow (pending → completed/canceled)
- Webhook audit trail for compliance
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
# Start everything (PostgreSQL + API) with automated tests
make start

# Or manually start services
docker-compose up -d

# View logs
make logs

# Stop everything
make stop
```

Server starts at `http://localhost:8080`

**Swagger UI:** `http://localhost:8080/swagger/index.html`

**Note:** `make start` automatically runs:

1. Unit tests (105 tests)
2. Service startup (PostgreSQL + API)
3. Integration tests (12 webhook scenarios)
4. Opens Swagger UI in browser

### Test API

**Via Swagger UI (Recommended):**
Visit `http://localhost:8080/swagger/index.html` for interactive API testing

**Via curl:**

```bash
# Register a new user (customer role by default)
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"pass123","name":"John Doe"}'

# Login and get JWT token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"pass123"}' \
  | jq -r '.token')

# Create a product (admin only)
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"High-performance","price":999.99,"quantity":50}'

# List products (public access)
curl http://localhost:8080/api/products

# Create an order (authenticated users)
curl -X POST http://localhost:8080/api/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"customer_id":123,"products":[{"product_id":"YOUR_PRODUCT_ID","quantity":2}]}'
```

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register new user account
- `POST /api/auth/login` - Login and receive JWT token

**📖 See [Authentication Documentation](docs/AUTHENTICATION.md) for complete guide**

**📖 See [Permissions Matrix](docs/PERMISSIONS.md) for role-based access control details**

### Products

- `POST /api/products` - Create product (**Admin only** 🔒)
- `GET /api/products` - List products (supports `?page=1&page_size=10&in_stock_only=true`) (Public)
- `GET /api/products/{id}` - Get product (Public)
- `PUT /api/products/{id}` - Update product (**Admin only** 🔒)
- `DELETE /api/products/{id}` - Delete product (**Admin only** 🔒)

### Orders

- `POST /api/orders` - Create order (Authenticated 🔒)
- `GET /api/orders` - List orders (supports `?page=1&page_size=10&status=pending`) (Authenticated 🔒)
- `GET /api/orders/{id}` - Get order (Authenticated 🔒)
- `PUT /api/orders/{id}/status` - Update order status (**Admin only** 🔒)

### Payment Webhooks

- `POST /api/payment-webhook` - Receive payment status updates (Public with signature verification)
- `GET /api/orders/{id}/payment-history` - Get payment webhook history (**Admin only** 🔒)

**📖 See [Payment Webhook Documentation](docs/PAYMENT_WEBHOOK.md) for complete integration guide**

## Testing

### Unit Tests

Run all unit tests:

```bash
# Run unit tests in Docker
make test
```

**Test Coverage:**

- **Domain entities: 100.0% coverage** ✅ (Product, Order, User validation & business logic)
- **DTO mappers: 100.0% coverage** ✅
- **HTTP handlers: 100.0% coverage** ✅
- **Product use cases: 100.0% coverage** ✅
- **Order use cases: 95.1% coverage** ✅
- **JWT Provider: 100.0% coverage** ✅
- **Total: 95 passing tests across 8 test suites**

**Test Suites:**

- Entity layer: Product, Order & User business logic validation, password hashing, GORM hooks
- DTO layer: Request/Response mapping and pagination
- Handler layer: HTTP request/response handling, validation, error responses
- Use case layer: Product & Order CRUD operations with comprehensive error handling
- Infrastructure layer: JWT token generation, validation, expiration, and security
- All edge cases: Invalid inputs, repository errors, validation failures, pagination defaults

### Integration Tests

**Authentication Tests:**

```bash
# Run authentication integration tests
make test-auth
```

**Auth Test Coverage (11 scenarios):**

✅ **Registration & Login:**
- Successful user registration (customer role)
- Login with valid credentials
- Login with invalid credentials (401)

✅ **Token Validation:**
- Access protected endpoint with valid token
- Access protected endpoint without token (401)
- Access protected endpoint with invalid token (401)

✅ **Permission Tests:**
- Customer can view products (public)
- Customer can create orders (authenticated)
- Customer cannot create products (403 Forbidden)
- Customer cannot update order status (403 Forbidden)
- Admin can create products and update order status

**Webhook Tests:**

```bash
# Run payment webhook integration tests
make test-webhook
```

**Webhook Test Coverage (12 scenarios):**

✅ **Security Tests:**

- Missing HMAC signature (401)
- Invalid HMAC signature (401)

✅ **Validation Tests:**

- Missing transaction ID (400)
- Invalid order ID format (400)
- Non-existent order (400)
- Invalid payment status (400)

✅ **Business Logic Tests:**

- Successful payment processing (200)
- Failed payment processing (200)

✅ **Resilience Tests:**

- Idempotency with duplicate transactions
- Webhook on already completed order
- Webhook history audit trail
- Concurrent webhook handling (race conditions)

**Note:** Integration tests run automatically with `make start`

## Architecture

```
src/
├── cmd/api/              # Entry point (main, container, routes)
├── internal/
│   ├── domain/           # Entities & repository interfaces
│   │   ├── entity/       # User, Product, Order, WebhookLog
│   │   └── repository/   # Repository interfaces
│   ├── infrastructure/   # Repository implementations (PostgreSQL)
│   │   ├── auth/         # JWT provider
│   │   ├── database/     # Database connection & migrations
│   │   └── repository/   # PostgreSQL implementations
│   ├── adapter/http/
│   │   ├── handler/      # HTTP handlers (auth, product, order, payment)
│   │   ├── middleware/   # Authentication & authorization
│   │   └── dto/          # Data Transfer Objects
│   └── config/           # Configuration
└── usecase/              # Business logic (auth, product, order, payment)
```

## Make Commands

```bash
make start         # Start services + run all tests (unit + integration)
make stop          # Stop all services
make logs          # View service logs
make test          # Run unit tests in Docker
make test-webhook  # Run webhook integration tests
make test-auth     # Run authentication integration tests
make help          # Show available commands
```

## Configuration

Environment variables (defaults):

- `DB_HOST=localhost`
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=postgres`
- `DB_NAME=ecommerce`
- `SERVER_PORT=8080`
- `JWT_SECRET=your-secret-key` (⚠️ Change in production!)
- `JWT_EXPIRATION_HOURS=24` (Token validity period)
- `WEBHOOK_SECRET=your-webhook-secret-key` (⚠️ Change in production!)

## Project Highlights

✨ **Clean Architecture** - Separation of concerns with domain, use case, and infrastructure layers  
🔐 **JWT Authentication** - Secure token-based authentication with bcrypt password hashing  
🛡️ **Role-Based Access Control** - Fine-grained permission system (admin vs customer)  
🧪 **Comprehensive Testing** - 95 unit tests + 11 auth integration tests + 12 webhook integration tests with 95%+ coverage  
🔒 **Webhook Security** - HMAC-SHA256 signature verification for payment webhooks  
🔄 **Idempotency** - Transaction ID-based duplicate prevention  
📊 **Audit Trail** - Complete webhook event logging with status tracking  
⚡ **Retry Logic** - Webhook status tracking for payment processor retries  
📖 **API Documentation** - Interactive Swagger UI with complete endpoint documentation  
🐳 **Containerized** - Docker & Docker Compose for easy deployment  
🚀 **CI/CD Ready** - Automated testing on startup

## License

MIT
