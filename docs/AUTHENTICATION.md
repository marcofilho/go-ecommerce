# Authentication & Authorization

This document describes the authentication and authorization system implemented in the Go E-Commerce API.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [User Roles](#user-roles)
- [Endpoints](#endpoints)
- [Authentication Flow](#authentication-flow)
- [Protected Routes](#protected-routes)
- [Token Structure](#token-structure)
- [Usage Examples](#usage-examples)
- [Testing](#testing)
- [Security Considerations](#security-considerations)

## Overview

The API implements **JWT (JSON Web Token)** based authentication with **role-based access control (RBAC)**. The system includes:

- User registration and login
- JWT token generation and validation
- Authentication middleware
- Role-based authorization middleware
- Secure password hashing (bcrypt)

## Architecture

The authentication system follows Clean Architecture principles:

```
src/
├── domain/entity/
│   └── user.go                  # User entity with validation & password hashing
├── domain/repository/
│   └── user_repository.go       # User repository interface
├── infrastructure/
│   ├── auth/
│   │   └── jwt_provider.go      # JWT token generation & validation
│   └── repository/
│       └── user_repository_postgres.go  # PostgreSQL user repository
├── usecase/auth/
│   └── auth_usecase.go          # Authentication business logic
└── adapter/http/
    ├── handler/
    │   └── auth_handler.go      # HTTP handlers for auth endpoints
    └── middleware/
        ├── auth_middleware.go   # Authentication & authorization middleware
        └── context.go           # Context helper functions
```

### Key Components

1. **User Entity** (`entity/user.go`)
   - Validates user data
   - Hashes passwords using bcrypt
   - Checks password validity
   - Role management

2. **JWT Provider** (`infrastructure/auth/jwt_provider.go`)
   - Generates JWT tokens
   - Validates and parses tokens
   - Manages token expiration

3. **Auth Use Case** (`usecase/auth/auth_usecase.go`)
   - Handles registration logic
   - Authenticates users
   - Generates auth responses

4. **Auth Middleware** (`middleware/auth_middleware.go`)
   - Validates JWT tokens
   - Injects user context into requests
   - Enforces role-based access control

## User Roles

The system supports two roles:

| Role | Description | Permissions |
|------|-------------|-------------|
| `customer` | Default role for registered users | Can view products, create orders, view own orders |
| `admin` | Administrative role | Full access - can manage products, view all orders, update order status |

**Note:** Currently, all registered users get the `customer` role by default. In production, you should implement a separate endpoint or database seeder to create admin users.

## Endpoints

### Public Endpoints (No Authentication Required)

#### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "role": "customer",
  "expires_at": "2025-12-06T12:00:00Z"
}
```

#### Login User
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "role": "customer",
  "expires_at": "2025-12-06T12:00:00Z"
}
```

### Protected Endpoints

All protected endpoints require the `Authorization` header with a valid JWT token:

```http
Authorization: Bearer <your-jwt-token>
```

## Authentication Flow

```
┌─────────┐                 ┌─────────┐                ┌──────────┐
│ Client  │                 │   API   │                │ Database │
└────┬────┘                 └────┬────┘                └────┬─────┘
     │                           │                          │
     │  1. POST /auth/register   │                          │
     ├──────────────────────────>│                          │
     │                           │  2. Hash password        │
     │                           │  3. Create user          │
     │                           ├─────────────────────────>│
     │                           │                          │
     │                           │  4. User created         │
     │                           │<─────────────────────────┤
     │                           │  5. Generate JWT         │
     │                           │                          │
     │  6. Return token & user   │                          │
     │<──────────────────────────┤                          │
     │                           │                          │
     │  7. Request with token    │                          │
     │    Authorization: Bearer  │                          │
     ├──────────────────────────>│                          │
     │                           │  8. Validate token       │
     │                           │  9. Extract claims       │
     │                           │ 10. Inject user context  │
     │                           │ 11. Check permissions    │
     │                           │ 12. Process request      │
     │                           │                          │
     │  13. Return response      │                          │
     │<──────────────────────────┤                          │
     │                           │                          │
```

## Protected Routes

### Authentication Required (Any Role)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/orders` | Create a new order |
| GET | `/api/orders` | List user's orders |
| GET | `/api/orders/{id}` | Get specific order |

### Admin Only

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/products` | Create product |
| PUT | `/api/products/{id}` | Update product |
| DELETE | `/api/products/{id}` | Delete product |
| PUT | `/api/orders/{id}/status` | Update order status |
| GET | `/api/orders/{id}/payment-history` | View webhook history |

### Public (No Authentication)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/products` | List all products |
| GET | `/api/products/{id}` | Get specific product |
| POST | `/api/payment-webhook` | Payment webhook (external) |

## Token Structure

JWT tokens contain the following claims:

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "role": "customer",
  "iss": "go-ecommerce",
  "exp": 1733489123,
  "iat": 1733402723
}
```

**Claims:**
- `user_id`: UUID of the user
- `email`: User's email address
- `role`: User's role (`customer` or `admin`)
- `iss`: Issuer (always "go-ecommerce")
- `exp`: Expiration time (Unix timestamp)
- `iat`: Issued at (Unix timestamp)

**Token Expiration:** 24 hours (configurable via `JWT_EXPIRATION_HOURS`)

## Usage Examples

### 1. Register and Get Token

```bash
# Register a new user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword",
    "name": "John Doe"
  }'

# Save the token from the response
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 2. Login with Existing User

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### 3. Access Protected Endpoint

```bash
# Create an order (requires authentication)
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "customer_id": 1,
    "products": [
      {
        "product_id": "550e8400-e29b-41d4-a716-446655440000",
        "quantity": 2
      }
    ]
  }'
```

### 4. Admin Operations

```bash
# Create a product (requires admin role)
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 1299.99,
    "quantity": 50
  }'
```

### 5. Error Responses

**Missing Token (401):**
```json
{
  "error": "missing authorization header"
}
```

**Invalid Token (401):**
```json
{
  "error": "invalid or expired token"
}
```

**Insufficient Permissions (403):**
```json
{
  "error": "insufficient permissions"
}
```

**Invalid Credentials (401):**
```json
{
  "error": "invalid credentials"
}
```

## Testing

Run the comprehensive authentication test suite:

```bash
# Start the application
make start

# In another terminal, run the tests
./test_authentication.sh
```

The test script covers:
- ✅ User registration
- ✅ User login
- ✅ Token validation
- ✅ Protected endpoint access control
- ✅ Role-based authorization
- ✅ Proper error responses (401, 403)

## Security Considerations

### Implemented Security Features

1. **Password Hashing**
   - Uses bcrypt with default cost factor (10)
   - Passwords never stored in plain text
   - Constant-time comparison for password verification

2. **JWT Token Security**
   - Signed with HMAC-SHA256
   - Includes expiration time
   - Contains minimal user data
   - Secret key configurable via environment variable

3. **Input Validation**
   - Email format validation
   - Password minimum length (6 characters)
   - Name minimum length (2 characters)
   - All inputs sanitized

4. **HTTP Security**
   - Tokens transmitted in Authorization header (not URL)
   - Content-Type validation
   - Proper error messages (no information leakage)

### Best Practices for Production

1. **JWT Secret**
   ```bash
   # Use a strong, random secret (32+ characters)
   JWT_SECRET=$(openssl rand -base64 32)
   ```

2. **Token Storage (Client-Side)**
   - Store tokens securely (e.g., httpOnly cookies)
   - Never store in localStorage for sensitive apps
   - Clear tokens on logout

3. **Token Refresh**
   - Implement refresh token mechanism
   - Short-lived access tokens (15-30 minutes)
   - Long-lived refresh tokens (7-30 days)

4. **Rate Limiting**
   - Implement rate limiting on auth endpoints
   - Prevent brute force attacks
   - Use tools like `golang.org/x/time/rate`

5. **HTTPS**
   - Always use HTTPS in production
   - Tokens sent over encrypted connection only

6. **Account Security**
   - Implement password reset flow
   - Email verification
   - Multi-factor authentication (MFA)
   - Account lockout after failed attempts

7. **Token Blacklisting**
   - Implement token revocation for logout
   - Use Redis to store blacklisted tokens
   - Check blacklist in middleware

### Environment Variables

Configure authentication in your environment:

```bash
# JWT Configuration
JWT_SECRET=your-jwt-secret-key-change-in-production-use-strong-secret
JWT_EXPIRATION_HOURS=24

# Database (users stored here)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ecommerce
```

## Implementation Details

### Middleware Chaining

The system supports flexible middleware chaining:

```go
// Require authentication only
mux.Handle("GET /api/orders", 
  c.AuthMiddleware.Authenticate(
    http.HandlerFunc(c.OrderHandler.ListOrders),
  ),
)

// Require authentication + admin role
mux.Handle("POST /api/products", 
  c.AuthMiddleware.Authenticate(
    c.AuthMiddleware.RequireRole(entity.RoleAdmin)(
      http.HandlerFunc(c.ProductHandler.CreateProduct),
    ),
  ),
)

// Optional authentication (token validation if present)
mux.Handle("GET /api/products", 
  c.AuthMiddleware.OptionalAuth(
    http.HandlerFunc(c.ProductHandler.ListProducts),
  ),
)
```

### Accessing User Context in Handlers

```go
import "github.com/marcofilho/go-ecommerce/src/internal/adapter/http/middleware"

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    // Get authenticated user from context
    claims, err := middleware.GetUserFromContext(r)
    if err != nil {
        respondError(w, http.StatusUnauthorized, "unauthorized")
        return
    }
    
    // Access user data
    userID := claims.UserID
    email := claims.Email
    role := claims.Role
    
    // Use in business logic...
}
```

## Future Enhancements

- [ ] Refresh token mechanism
- [ ] Password reset flow
- [ ] Email verification
- [ ] OAuth2/Social login (Google, GitHub)
- [ ] Multi-factor authentication (MFA)
- [ ] Account lockout after failed attempts
- [ ] Password strength requirements
- [ ] Token blacklisting for logout
- [ ] Audit log for authentication events
- [ ] IP-based rate limiting
