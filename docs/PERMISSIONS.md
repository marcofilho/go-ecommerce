# Role-Based Permissions Matrix

This document provides a comprehensive overview of the role-based access control (RBAC) system implemented in the Go E-Commerce API.

## Overview

The API implements a **permission-based authorization system** where each role has a specific set of permissions. The middleware validates these permissions before allowing access to protected endpoints.

## Roles

| Role | Description | Default on Registration |
|------|-------------|------------------------|
| `customer` | Standard user role | ✅ Yes |
| `admin` | Administrative role with full access | ❌ No (manual assignment) |

## Permission System

### Permission Constants

All permissions are defined in `middleware/permissions.go`:

```go
// Product permissions
PermissionCreateProduct  = "product:create"
PermissionUpdateProduct  = "product:update"
PermissionDeleteProduct  = "product:delete"
PermissionViewProduct    = "product:view"
PermissionListProducts   = "product:list"

// Order permissions
PermissionCreateOrder      = "order:create"
PermissionViewOrder        = "order:view"
PermissionListOrders       = "order:list"
PermissionUpdateOrderStatus = "order:update_status"

// Webhook permissions
PermissionViewWebhookHistory = "webhook:view_history"
```

## Complete Permission Matrix

| Permission | Customer | Admin | Description |
|------------|----------|-------|-------------|
| **Products** |
| `product:view` | ✅ | ✅ | View single product details |
| `product:list` | ✅ | ✅ | List all products with pagination |
| `product:create` | ❌ | ✅ | Create new products |
| `product:update` | ❌ | ✅ | Update existing products |
| `product:delete` | ❌ | ✅ | Delete products |
| **Orders** |
| `order:create` | ✅ | ✅ | Create new orders |
| `order:view` | ✅ | ✅ | View order details |
| `order:list` | ✅ | ✅ | List orders |
| `order:update_status` | ❌ | ✅ | Update order status (pending → completed/canceled) |
| **Webhooks** |
| `webhook:view_history` | ❌ | ✅ | View payment webhook history |

## Endpoint Authorization

### Public Endpoints (No Authentication)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user account |
| POST | `/api/auth/login` | Login and receive JWT token |
| GET | `/api/products` | List all products |
| GET | `/api/products/{id}` | Get specific product |
| POST | `/api/payment-webhook` | Payment gateway webhook |

### Customer Endpoints

Customers can perform the following actions:

#### Products (Read-Only)
```bash
# View products (public, but customers have permission too)
GET /api/products
GET /api/products/{id}
```

#### Orders
```bash
# Create order (requires: order:create)
POST /api/orders
Authorization: Bearer <customer-token>

# View orders (requires: order:list)
GET /api/orders
Authorization: Bearer <customer-token>

# View specific order (requires: order:view)
GET /api/orders/{id}
Authorization: Bearer <customer-token>
```

**Forbidden Actions for Customers:**
- ❌ Create/Update/Delete products
- ❌ Update order status
- ❌ View webhook history

### Admin Endpoints

Admins have **all permissions** and can perform any action:

#### Product Management
```bash
# Create product (requires: product:create)
POST /api/products
Authorization: Bearer <admin-token>

# Update product (requires: product:update)
PUT /api/products/{id}
Authorization: Bearer <admin-token>

# Delete product (requires: product:delete)
DELETE /api/products/{id}
Authorization: Bearer <admin-token>
```

#### Order Management
```bash
# All customer order actions PLUS:

# Update order status (requires: order:update_status)
PUT /api/orders/{id}/status
Authorization: Bearer <admin-token>
```

#### Webhook Management
```bash
# View webhook history (requires: webhook:view_history)
GET /api/orders/{id}/payment-history
Authorization: Bearer <admin-token>
```

## Authorization Flow

```
┌─────────┐           ┌──────────────┐           ┌─────────────┐
│ Client  │           │  Middleware  │           │  Handler    │
└────┬────┘           └──────┬───────┘           └──────┬──────┘
     │                       │                          │
     │  1. Request + Token   │                          │
     ├──────────────────────>│                          │
     │                       │                          │
     │                       │ 2. Validate Token        │
     │                       │    (JWT signature)       │
     │                       │                          │
     │                       │ 3. Extract User Claims   │
     │                       │    (user_id, role)       │
     │                       │                          │
     │                       │ 4. Check Permission      │
     │                       │    HasPermission(role,   │
     │                       │    permission)?          │
     │                       │                          │
     │                       ├─ ❌ No Permission        │
     │  403 Forbidden        │                          │
     │<──────────────────────┤                          │
     │                       │                          │
     │                       ├─ ✅ Has Permission       │
     │                       │ 5. Inject User Context   │
     │                       │                          │
     │                       │ 6. Forward Request       │
     │                       ├─────────────────────────>│
     │                       │                          │
     │                       │                          │ 7. Process
     │                       │ 8. Response              │    Request
     │<──────────────────────┴──────────────────────────┤
     │                                                  │
```

## Testing Permissions

### Test Scenario 1: Customer tries to create product (should fail)

```bash
# Register as customer
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"pass123","name":"Customer"}' \
  | jq -r '.token')

# Try to create product (will fail with 403)
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Unauthorized Product",
    "price": 99.99,
    "quantity": 10
  }'

# Response:
# {
#   "error": "forbidden: insufficient permissions for this action"
# }
```

### Test Scenario 2: Customer creates order (should succeed)

```bash
# Customer can create orders
curl -X POST http://localhost:8080/api/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "products": [{"product_id":"uuid-here","quantity":2}]
  }'

# Response: 201 Created with order details
```

### Test Scenario 3: Customer tries to update order status (should fail)

```bash
# Try to update order status (will fail with 403)
curl -X PUT http://localhost:8080/api/orders/{id}/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"completed"}'

# Response:
# {
#   "error": "forbidden: insufficient permissions for this action"
# }
```

### Test Scenario 4: Admin performs all actions (should succeed)

```bash
# Register admin (in production, use database seeder)
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123","name":"Admin"}' \
  | jq -r '.token')

# Note: You'll need to manually update the role in database:
# UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';

# Admin can create products ✅
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Admin Product","price":199.99,"quantity":50}'

# Admin can update order status ✅
curl -X PUT http://localhost:8080/api/orders/{id}/status \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"completed"}'

# Admin can view webhook history ✅
curl -X GET http://localhost:8080/api/orders/{id}/payment-history \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## Error Responses

### 401 Unauthorized
Returned when:
- No `Authorization` header provided
- Invalid token format
- Expired token
- Invalid token signature

```json
{
  "error": "missing authorization header"
}
```

```json
{
  "error": "invalid or expired token"
}
```

### 403 Forbidden
Returned when:
- Valid authentication but insufficient permissions
- User role doesn't have required permission

```json
{
  "error": "forbidden: insufficient permissions for this action"
}
```

## Implementation Details

### Permission Check Flow

```go
// 1. Define permission in middleware/permissions.go
const PermissionCreateProduct Permission = "product:create"

// 2. Map permission to roles
var RolePermissions = map[entity.Role][]Permission{
    entity.RoleAdmin: {
        PermissionCreateProduct,
        // ... other permissions
    },
    entity.RoleCustomer: {
        // Customers don't have product:create
    },
}

// 3. Use in route configuration
mux.Handle("POST /api/products", 
    c.AuthMiddleware.Authenticate(
        c.AuthMiddleware.RequirePermission(middleware.PermissionCreateProduct)(
            http.HandlerFunc(c.ProductHandler.CreateProduct),
        ),
    ),
)

// 4. Middleware validates
func (m *AuthMiddleware) RequirePermission(permission Permission) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims := getUserFromContext(r)
            
            if !HasPermission(claims.Role, permission) {
                // 403 Forbidden
                m.writeError(w, "forbidden: insufficient permissions for this action", 
                    http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### Adding New Permissions

1. **Define the permission:**
```go
// middleware/permissions.go
const PermissionRefundOrder Permission = "order:refund"
```

2. **Assign to roles:**
```go
var RolePermissions = map[entity.Role][]Permission{
    entity.RoleAdmin: {
        // ... existing permissions
        PermissionRefundOrder,
    },
    entity.RoleCustomer: {
        // Customers cannot refund
    },
}
```

3. **Apply to route:**
```go
// routes.go
mux.Handle("POST /api/orders/{id}/refund", 
    c.AuthMiddleware.Authenticate(
        c.AuthMiddleware.RequirePermission(middleware.PermissionRefundOrder)(
            http.HandlerFunc(c.OrderHandler.RefundOrder),
        ),
    ),
)
```

## Security Best Practices

### 1. Principle of Least Privilege
- Users get minimum permissions needed
- Default role is `customer` (limited access)
- Admin role assigned manually

### 2. Permission Granularity
- Separate permissions for different actions
- Fine-grained control (create ≠ update ≠ delete)
- Easy to audit and modify

### 3. Explicit Denials
- No permission = denied by default
- Clear error messages (403 Forbidden)
- Logged for security auditing

### 4. Token-Based Auth
- Stateless authentication
- JWT tokens with expiration
- No server-side session storage

## Future Enhancements

- [ ] **Custom Roles**: Support for custom role creation
- [ ] **Dynamic Permissions**: Runtime permission assignment
- [ ] **Permission Groups**: Bundle related permissions
- [ ] **Hierarchical Roles**: Role inheritance (super-admin > admin > customer)
- [ ] **Resource-Level Permissions**: User can only view/edit their own orders
- [ ] **Audit Logging**: Track all permission checks
- [ ] **Permission Caching**: Cache role-permission mappings
- [ ] **API Keys**: Service-to-service authentication
- [ ] **Scoped Tokens**: Tokens with subset of permissions

## Summary

The permission system provides:

✅ **Clear separation** between customer and admin capabilities  
✅ **Explicit permission checks** at the middleware level  
✅ **Proper HTTP status codes** (401 for auth, 403 for permissions)  
✅ **Granular control** over each endpoint  
✅ **Easy to extend** with new permissions  
✅ **Production-ready** security model  

**Key Rules:**
- 📋 **Customers**: View products, manage their orders
- 🔑 **Admins**: Full control over products, orders, and webhooks
- 🚫 **Explicit denials**: 403 Forbidden with clear error message
- 🔒 **Secure by default**: No permission = no access
