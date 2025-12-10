# Go E-Commerce API

A RESTful API for managing products and orders in an e-commerce system, built with Go using clean architecture principles and PostgreSQL.

## Features

- **Authentication & Authorization** (JWT-based with RBAC)
- **Role-Based Permissions** (admin vs customer access control)
- Product Management (CRUD with stock tracking)
- **Product Categories** (N:N relationship - products can have multiple categories)
- **Product Variants** (support multiple variants per product with optional price overrides)
- Order Management (create orders with automatic stock deduction)
- **Advanced Payment Webhook Security**:
  - HMAC-SHA256 signature validation
  - Timestamp-based replay attack prevention (±5 minute tolerance)
  - Transaction ID-based idempotency
  - Complete audit trail for compliance
- Payment status tracking (unpaid → paid/failed)
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

1. Unit tests (276 tests)
2. Service startup (PostgreSQL + API)
3. Database auto-seeding (if empty) with sample data
4. Integration tests (webhook + auth scenarios)
5. Opens Swagger UI in browser

**Default Admin Account:** `admin@ecommerce.com` / `password123`

### Test API

**Via Swagger UI (Recommended):**
Visit `http://localhost:8080/swagger/index.html` for interactive API testing

**Via curl:**

```bash
# Register a new customer (public, no auth required)
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"pass123","name":"John Doe"}'

# Login and get JWT token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"pass123"}' \
  | jq -r '.token')

# Create admin account (requires existing admin authentication)
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@ecommerce.com","password":"password123"}' \
  | jq -r '.token')

curl -X POST http://localhost:8080/api/auth/register \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"newadmin@example.com","password":"secure123","name":"New Admin","role":"admin"}'

# Create a product (admin only)
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"High-performance","price":999.99,"quantity":50}'

# List products (public access - includes variants and categories)
curl http://localhost:8080/api/products

# Get specific product (includes variants and categories)
curl http://localhost:8080/api/products/YOUR_PRODUCT_ID

# Create an order (authenticated users)
curl -X POST http://localhost:8080/api/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"customer_id":123,"products":[{"product_id":"YOUR_PRODUCT_ID","quantity":2}]}'
```

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register new user (public: customer role, admin creation requires admin auth)
- `POST /api/auth/login` - Login and receive JWT token

**📖 See [Authentication Documentation](docs/AUTHENTICATION.md) for complete guide including admin account creation**

**📖 See [Permissions Matrix](docs/PERMISSIONS.md) for role-based access control details**

**📖 See [Database Schema Documentation](docs/DATABASE_SCHEMA.md) for complete database structure**

### Products

- `POST /api/products` - Create product (**Admin only** 🔒)
- `GET /api/products` - List products with categories and variants (supports `?page=1&page_size=10&in_stock_only=true`) (Public)
- `GET /api/products/{id}` - Get product with categories and variants (Public)
- `PUT /api/products/{id}` - Update product (**Admin only** 🔒)
- `DELETE /api/products/{id}` - Delete product (**Admin only** 🔒)

### Categories

- `POST /api/categories` - Create category (**Admin only** 🔒)
- `GET /api/categories` - List categories (supports `?page=1&page_size=10`) (Public)
- `POST /api/products/{id}/categories` - Assign category to product (**Admin only** 🔒)
- `DELETE /api/products/{id}/categories/{category_id}` - Remove category from product (**Admin only** 🔒)
- `GET /api/products/{id}/categories` - Get product categories (Public)

### Product Variants

- `POST /api/products/{id}/variants` - Create variant for a product (**Admin only** 🔒)
- `GET /api/products/{id}/variants` - List variants for a product (supports `?page=1&page_size=10`) (Public)
- `PUT /api/variants/{variant_id}` - Update variant (**Admin only** 🔒)
- `DELETE /api/variants/{variant_id}` - Delete variant (**Admin only** 🔒)

### Orders

- `POST /api/orders` - Create order (Authenticated 🔒)
- `GET /api/orders` - List orders (supports `?page=1&page_size=10&status=pending`) (Authenticated 🔒)
- `GET /api/orders/{id}` - Get order (Authenticated 🔒)
- `PUT /api/orders/{id}/status` - Update order status (**Admin only** 🔒)

### Payment Webhooks

- `POST /api/payment-webhook` - Receive payment status updates (Public with HMAC signature & timestamp verification)
- `GET /api/orders/{id}/payment-history` - Get payment webhook history (**Admin only** 🔒)

**📖 See [Payment Webhook Documentation](docs/PAYMENT_WEBHOOK.md) for complete integration guide including:**
- HMAC-SHA256 signature generation
- Timestamp-based replay attack prevention
- Security best practices
- Code examples and test scenarios

## Testing

### Unit Tests

Run all unit tests:

```bash
# Run unit tests in Docker
make test
```

### Integration Tests

Run authentication and authorization tests:

```bash
# Requires running API
make test-auth
```

### Load Tests

Test database integrity under concurrent operations:

```bash
# Requires running API  
./test_category_load.sh
```

This comprehensive load test verifies:
- ✅ **Concurrent category creation** (20 simultaneous requests)
- ✅ **Concurrent product-category assignments** (50 simultaneous requests)
- ✅ **Concurrent read operations** (100 simultaneous requests)
- ✅ **Concurrent deletion operations** (20 simultaneous requests)
- ✅ **Database integrity** (no orphaned records, constraint violations)
- ✅ **Composite primary key** (prevents duplicate assignments)
- ✅ **CASCADE DELETE** (proper cleanup on deletion)
- ✅ **N:N relationship queries** (eager loading performance)

**Test Coverage:**

- **Domain entities: 99.0% coverage** ✅ (Product, ProductVariant, Order, Category, User validation & business logic)
- **DTO mappers: 100.0% coverage** ✅
- **HTTP handlers: 100.0% coverage** ✅
- **Product use cases: 100.0% coverage** ✅
- **Product variant use cases: 100.0% coverage** ✅
- **Category use cases: 100.0% coverage** ✅
- **Order use cases: 95.1% coverage** ✅
- **JWT Provider: 100.0% coverage** ✅
- **Total: 276 passing tests across 18 test packages**

**Test Suites:**

- Entity layer: Product, ProductVariant, Order, Category & User business logic validation, password hashing, GORM hooks
  - Product: 13 tests (validation, stock management, variants & categories relationship)
  - ProductVariant: 15 tests (price override logic, validation, UUID generation)
  - Category: 4 tests (validation, UUID generation with BeforeCreate hook)
  - Order: 12 tests (comprehensive order workflow with variant support)
  - User: Authentication and validation tests
- DTO layer: Request/Response mapping and pagination
- Handler layer: HTTP request/response handling, validation, error responses
  - Category: 15 tests (CRUD, product-category assignment, error handling)
- Use case layer: 
  - Product: 22 tests (CRUD operations with comprehensive error handling)
  - ProductVariant: 27 tests (full variant lifecycle with price override logic)
  - Category: 20 tests (CRUD, pagination, product-category relationships)
  - Order: 16 tests (order creation with variant support, stock management)
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
- **Replay attack prevention** (timestamps outside ±5 minute window rejected)

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

The project follows **Clean Architecture** principles with **Dependency Inversion** - all layers depend on interfaces, not concrete implementations.

**📖 See [Architecture & Design Principles](docs/ARCHITECTURE.md) for detailed explanation of interfaces, SOLID principles, and testing strategies**

```
src/
├── cmd/api/              # Entry point (main, container, routes)
├── internal/
│   ├── domain/           # Entities & repository interfaces
│   │   ├── entity/       # User, Product, Order, WebhookLog
│   │   └── repository/   # Repository interfaces
│   ├── infrastructure/   # Repository implementations (PostgreSQL)
│   │   ├── auth/         # JWT provider (implements TokenProvider interface)
│   │   ├── database/     # Database connection & migrations
│   │   └── repository/   # PostgreSQL implementations
│   ├── adapter/http/
│   │   ├── handler/      # HTTP handlers (auth, product, order, payment)
│   │   ├── middleware/   # Authentication & authorization
│   │   └── dto/          # Data Transfer Objects
│   └── config/           # Configuration
└── usecase/              # Business logic (auth, product, order, payment)
                          # Each use case defines service interfaces
```

## Make Commands

```bash
# Services
make start         # Start services + run all tests (unit + integration)
make stop          # Stop all services
make logs          # View service logs

# Testing
make test          # Run unit tests in Docker
make test-webhook  # Run webhook integration tests
make test-auth     # Run authentication integration tests

# Database
make seed          # Manually seed database with sample data
make clean-db      # Clean database (with confirmation prompt)
make reset-db      # Reset database (clean + seed)

# Other
make help          # Show available commands
```

**Note:** The database is automatically seeded with sample data on first startup if empty. Sample credentials:
- Admin: `admin@ecommerce.com` / `password123`
- Customer: `john.doe@example.com` / `password123`

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
🏛️ **SOLID Principles** - Interface-based design following Dependency Inversion Principle ([Architecture Guide](docs/ARCHITECTURE.md))  
🔐 **JWT Authentication** - Secure token-based authentication with bcrypt password hashing  
🛡️ **Role-Based Access Control** - Fine-grained permission system with admin privilege restrictions  
👥 **Secure Admin Creation** - Admin accounts require authenticated admin authorization  
🏷️ **Product Categories** - N:N relationship supporting multiple categories per product  
🎨 **Product Variants** - Support for multiple product variants with optional price overrides  
🧪 **Comprehensive Testing** - 282 unit tests + 17 auth tests + 12 webhook tests with 95%+ coverage  
🔒 **Advanced Webhook Security**:
  - HMAC-SHA256 signature verification with `X-Payment-Signature` header
  - Timestamp-based replay attack prevention (±5 minute tolerance window)
  - Proper HTTP status codes (401 for auth failures, 200 for success)
🔄 **Idempotency** - Transaction ID-based duplicate prevention  
📊 **Audit Trail** - Complete webhook event logging with status tracking  
⚡ **Retry Logic** - Webhook status tracking for payment processor retries  
🗄️ **Auto-Seeding** - Automatic database population with sample data on first startup  
📖 **API Documentation** - Interactive Swagger UI with complete endpoint documentation  
🐳 **Containerized** - Docker & Docker Compose for easy deployment  
🚀 **CI/CD Ready** - Automated testing on startup

## License

MIT
