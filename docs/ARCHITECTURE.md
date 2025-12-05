# Architecture & Design Principles

## Overview

This project follows **Clean Architecture** principles with a strong emphasis on the **Dependency Inversion Principle** (SOLID). All dependencies point inward toward the domain, and outer layers depend on interfaces rather than concrete implementations.

## Dependency Inversion Implementation

### Core Principle

> **"Depend upon abstractions, not concretions"** - Robert C. Martin

All handlers and use cases depend on **interfaces**, not concrete implementations. This enables:

- **Testability**: Easy to mock dependencies in unit tests
- **Flexibility**: Swap implementations without changing business logic
- **Maintainability**: Changes to implementations don't affect dependents
- **Decoupling**: Clear boundaries between layers

### Interface Hierarchy

```
┌─────────────────────────────────────────────────┐
│           Adapter Layer (HTTP Handlers)         │
│  Depends on: Service Interfaces                 │
│  - AuthHandler → AuthService                    │
│  - ProductHandler → ProductService              │
│  - OrderHandler → OrderService                  │
│  - PaymentHandler → PaymentService              │
└─────────────────────────────────────────────────┘
                      ↓ depends on
┌─────────────────────────────────────────────────┐
│            Use Case Layer (Services)            │
│  Defines: Service Interfaces                    │
│  Depends on: Repository & Provider Interfaces   │
│  - AuthService (interface)                      │
│    - UseCase (implementation)                   │
│      → UserRepository                           │
│      → TokenProvider                            │
└─────────────────────────────────────────────────┘
                      ↓ depends on
┌─────────────────────────────────────────────────┐
│         Infrastructure Layer                    │
│  Implements: TokenProvider interface            │
│  - JWTProvider implements TokenProvider         │
└─────────────────────────────────────────────────┘
                      ↓ depends on
┌─────────────────────────────────────────────────┐
│            Domain Layer (Core)                  │
│  Defines: Repository Interfaces                 │
│  - UserRepository                               │
│  - ProductRepository                            │
│  - OrderRepository                              │
│  - WebhookRepository                            │
└─────────────────────────────────────────────────┘
```

## Interface Definitions

### 1. Authentication Layer

**TokenProvider Interface** (`src/internal/infrastructure/auth/jwt_provider.go`)
```go
type TokenProvider interface {
    GenerateToken(user *entity.User) (string, error)
    ValidateToken(tokenString string) (*Claims, error)
}
```

**AuthService Interface** (`src/usecase/auth/auth_usecase.go`)
```go
type AuthService interface {
    Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
    Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
    ValidateToken(tokenString string) (*Claims, error)
}
```

**Implementation**: `UseCase` implements `AuthService`

**Benefits**:
- AuthHandler can be tested with mock AuthService
- JWT implementation can be swapped (e.g., to OAuth, SAML) without changing handler code

### 2. Product Layer

**ProductService Interface** (`src/usecase/product/product_usecase.go`)
```go
type ProductService interface {
    CreateProduct(ctx context.Context, name, description string, price float64, quantity int) (*entity.Product, error)
    GetProduct(ctx context.Context, id uuid.UUID) (*entity.Product, error)
    ListProducts(ctx context.Context, page, pageSize int, inStockOnly bool) ([]*entity.Product, int, error)
    UpdateProduct(ctx context.Context, id uuid.UUID, name, description string, price float64, quantity int) (*entity.Product, error)
    DeleteProduct(ctx context.Context, id uuid.UUID) error
}
```

**Implementation**: `UseCase` implements `ProductService`

### 3. Order Layer

**OrderService Interface** (`src/usecase/order/order_usecase.go`)
```go
type OrderService interface {
    CreateOrder(ctx context.Context, customerID int, items []CreateOrderItem) (*entity.Order, error)
    GetOrder(ctx context.Context, id uuid.UUID) (*entity.Order, error)
    ListOrders(ctx context.Context, page, pageSize int, status *entity.OrderStatus, paymentStatus *entity.PaymentStatus) ([]*entity.Order, int, error)
    UpdateOrderStatus(ctx context.Context, id uuid.UUID, newStatus entity.OrderStatus) (*entity.Order, error)
}
```

**Implementation**: `UseCase` implements `OrderService`

### 4. Payment Layer

**PaymentService Interface** (`src/usecase/payment/payment_usecase.go`)
```go
type PaymentService interface {
    ProcessWebhook(ctx context.Context, req *entity.PaymentWebhookRequest) error
    GetWebhookHistory(ctx context.Context, orderID string) ([]entity.WebhookLog, error)
}
```

**Implementation**: `PaymentUseCase` implements `PaymentService`

## Dependency Injection

The `Container` (`src/cmd/api/container.go`) wires up all dependencies:

```go
// Concrete implementations
jwtProvider := auth.NewJWTProvider(secret, expiration)
authUseCase := authUseCase.NewUseCase(userRepo, jwtProvider)

// Injected as interfaces
authHandler := handler.NewAuthHandler(authUseCase) // accepts AuthService
```

**Key Points**:
- Container creates concrete implementations
- Handlers and middleware receive interfaces
- Production code uses real implementations
- Test code uses mocks

## Testing Strategy

### Unit Tests with Mocks

Example: `auth_handler_test.go`

```go
type mockAuthService struct {
    registerFunc func(ctx context.Context, req authUseCase.RegisterRequest) (*authUseCase.AuthResponse, error)
    // ... other methods
}

func (m *mockAuthService) Register(ctx context.Context, req authUseCase.RegisterRequest) (*authUseCase.AuthResponse, error) {
    return m.registerFunc(ctx, req)
}

// Test with mock
mockService := &mockAuthService{
    registerFunc: func(ctx context.Context, req authUseCase.RegisterRequest) (*authUseCase.AuthResponse, error) {
        return &authUseCase.AuthResponse{Token: "test-token"}, nil
    },
}
handler := NewAuthHandler(mockService)
```

### Test Coverage

| Package | Coverage | Details |
|---------|----------|---------|
| `infrastructure/auth` | 92.9% | JWT provider with edge cases |
| `domain/entity` | 98.7% | User, Product, Order entities |
| `adapter/http/handler` | 78.3% | HTTP handlers with mocks |
| `usecase/product` | 95%+ | Product business logic |
| `usecase/order` | 95%+ | Order business logic |

**Total**: 105 unit tests

## Benefits Achieved

### 1. **Testability**
- ✅ Each layer can be tested in isolation
- ✅ No need for database or external services in unit tests
- ✅ Mock implementations are simple and maintainable

### 2. **Maintainability**
- ✅ Clear separation of concerns
- ✅ Changes to implementations don't break dependents
- ✅ Easy to understand dependencies through interfaces

### 3. **Flexibility**
- ✅ Swap JWT for OAuth without touching handlers
- ✅ Switch from PostgreSQL to MongoDB by implementing repository interfaces
- ✅ Add caching layer transparently

### 4. **SOLID Compliance**
- ✅ **S**ingle Responsibility: Each interface has one reason to change
- ✅ **O**pen/Closed: Open for extension, closed for modification
- ✅ **L**iskov Substitution: Implementations can be substituted
- ✅ **I**nterface Segregation: Focused, minimal interfaces
- ✅ **D**ependency Inversion: Depend on abstractions ✨

## Example: Adding a New Feature

### Scenario: Add Redis Caching to Product Service

**Step 1**: Create decorator implementing `ProductService`

```go
type CachedProductService struct {
    base  product.ProductService
    cache *redis.Client
}

func (c *CachedProductService) GetProduct(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
    // Check cache first
    if cached, err := c.cache.Get(ctx, id.String()).Result(); err == nil {
        return unmarshal(cached), nil
    }
    
    // Delegate to base service
    p, err := c.base.GetProduct(ctx, id)
    if err == nil {
        c.cache.Set(ctx, id.String(), marshal(p), time.Hour)
    }
    return p, err
}
```

**Step 2**: Update container wiring

```go
baseService := productUseCase.NewUseCase(productRepo)
cachedService := NewCachedProductService(baseService, redisClient)
productHandler := handler.NewProductHandler(cachedService)
```

**Benefits**:
- ✅ No changes to `ProductHandler`
- ✅ No changes to `ProductUseCase`
- ✅ Tests still work with mocks
- ✅ Can enable/disable caching via configuration

## Clean Architecture Layers

### 1. **Domain Layer** (Innermost)
- Entities: `User`, `Product`, `Order`
- Repository interfaces
- No external dependencies
- Pure business logic

### 2. **Use Case Layer**
- Service interfaces
- Business rules orchestration
- Depends on domain interfaces

### 3. **Adapter Layer**
- HTTP handlers
- DTOs and mappers
- Depends on use case interfaces

### 4. **Infrastructure Layer** (Outermost)
- Database implementations
- JWT provider
- External service clients
- Implements domain/use case interfaces

## Best Practices

### DO ✅
- Define interfaces close to where they're consumed
- Keep interfaces small and focused
- Use dependency injection via constructors
- Mock interfaces in unit tests
- Return interfaces from factories when needed

### DON'T ❌
- Pass concrete types to handlers/use cases
- Define interfaces in implementation packages
- Create "god interfaces" with many methods
- Mock concrete types (use interfaces instead)
- Skip interface definitions for "simplicity"

## References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)
- [Go Proverbs: "The bigger the interface, the weaker the abstraction"](https://go-proverbs.github.io/)
