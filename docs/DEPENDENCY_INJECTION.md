# Dependency Injection Pattern

## Overview

This project uses a **Services Pattern** for dependency injection, which provides a scalable and maintainable approach to managing dependencies across use cases.

## Architecture

### Services Struct

Located in `src/cmd/api/container.go`, the `Services` struct groups common infrastructure services:

```go
type Services struct {
    audit audit.AuditService
}

func (s *Services) GetAuditService() audit.AuditService {
    return s.audit
}
```

### Use Case Pattern

Use cases that require infrastructure services receive a `Services` interface instead of individual dependencies:

```go
type UseCase struct {
    repo     repository.ProductRepository
    services Services  // Interface, not concrete type
}

type Services interface {
    GetAuditService() audit.AuditService
}

func NewUseCase(repo repository.ProductRepository, services Services) *UseCase {
    return &UseCase{
        repo:     repo,
        services: services,
    }
}
```

### Using Services in Use Cases

Services are accessed via getter methods:

```go
func (uc *UseCase) CreateProduct(ctx context.Context, name, description string, price float64, quantity int) (*entity.Product, error) {
    // Business logic...
    
    // Access audit service
    uc.services.GetAuditService().LogChange(
        ctx,
        userID,
        "CREATE",
        "product",
        product.ID,
        nil,
        product,
    )
    
    return product, nil
}
```

## Benefits

### 1. **Scalability**
Adding a new service requires updates in only 2 places:
- Add field to `Services` struct
- Add getter method

**Before (Old Pattern):**
```go
// Had to update:
// 1. Use case struct (add field)
// 2. Use case constructor (add parameter)
// 3. Container initialization (add argument)
// 4. Every test file (add mock parameter)
```

**After (New Pattern):**
```go
// Only update:
// 1. Services struct (add field + getter)
// 2. Tests continue working with existing MockServices
```

### 2. **Testability**
Centralized mock services in `src/internal/testing/mocks.go`:

```go
type MockServices struct {
    AuditService audit.AuditService
}

func (m *MockServices) GetAuditService() audit.AuditService {
    if m.AuditService != nil {
        return m.AuditService
    }
    return &MockAuditService{}
}
```

Tests can use the default mock or inject custom behavior:

```go
// Default behavior
uc := NewUseCase(mockRepo, &mockServices.MockServices{})

// Custom behavior
customMock := &mockServices.MockServices{
    AuditService: &CustomAuditMock{},
}
uc := NewUseCase(mockRepo, customMock)
```

### 3. **Separation of Concerns**
- Use cases depend on **interfaces**, not concrete implementations
- Services can be swapped without changing use case code
- Clear boundary between business logic and infrastructure

### 4. **Maintainability**
- Single source of truth for service dependencies
- Consistent pattern across all use cases
- Easy to understand and follow for new developers

## When to Use Services Pattern

### ✅ Use Services for:
- **Infrastructure services** (audit, logging, metrics, caching)
- **Cross-cutting concerns** needed by multiple use cases
- **Stateless services** that provide utility functions

### ❌ Don't use Services for:
- **Domain repositories** - these should remain direct dependencies
- **Use case-specific dependencies** - pass these directly
- **Configuration values** - use config structs or environment variables

## Migration Guide

### Adding a New Service

1. **Add to Services struct:**
```go
type Services struct {
    audit  audit.AuditService
    logger logger.LoggerService  // New service
}
```

2. **Add getter method:**
```go
func (s *Services) GetLoggerService() logger.LoggerService {
    return s.logger
}
```

3. **Initialize in container:**
```go
services := &Services{
    audit:  auditService,
    logger: loggerService,  // New initialization
}
```

4. **Add to Services interface in use cases:**
```go
type Services interface {
    GetAuditService() audit.AuditService
    GetLoggerService() logger.LoggerService  // Add this
}
```

5. **Update MockServices (if needed):**
```go
type MockServices struct {
    AuditService  audit.AuditService
    LoggerService logger.LoggerService  // New field
}

func (m *MockServices) GetLoggerService() logger.LoggerService {
    if m.LoggerService != nil {
        return m.LoggerService
    }
    return &MockLoggerService{}  // Default mock
}
```

### Converting a Use Case to Use Services

If a use case doesn't currently use the Services pattern:

1. **Add Services field to use case:**
```go
type UseCase struct {
    repo     repository.SomeRepository
    services Services  // Add this
}
```

2. **Define Services interface:**
```go
type Services interface {
    GetAuditService() audit.AuditService
    // Add other getters as needed
}
```

3. **Update constructor:**
```go
func NewUseCase(repo repository.SomeRepository, services Services) *UseCase {
    return &UseCase{
        repo:     repo,
        services: services,
    }
}
```

4. **Update container initialization:**
```go
someUseCase := somepackage.NewUseCase(someRepo, services)
```

5. **Update tests:**
```go
import mockServices "github.com/marcofilho/go-ecommerce/src/internal/testing"

uc := NewUseCase(mockRepo, &mockServices.MockServices{})
```

## Current Implementation Status

### Use Cases Using Services Pattern

| Use Case | Services Used |
|----------|--------------|
| **Product** | ✅ Audit Service |
| **Order** | ✅ Audit Service |
| **Payment** | ✅ Audit Service |
| **Category** | ❌ No services needed |
| **Product Variant** | ❌ No services needed |
| **Auth** | ❌ No services needed |

### Test Files

All test files use centralized `MockServices` from `src/internal/testing/mocks.go`:

- ✅ `src/usecase/product/product_usecase_test.go`
- ✅ `src/usecase/order/order_usecase_test.go`
- ✅ `src/internal/adapter/http/handler/product_handler_test.go`
- ✅ `src/internal/adapter/http/handler/order_handler_test.go`

## Best Practices

### 1. Keep Interfaces Minimal
Only include getters for services actually used by the use case:

```go
// Good - only what's needed
type Services interface {
    GetAuditService() audit.AuditService
}

// Bad - exposing unnecessary services
type Services interface {
    GetAuditService() audit.AuditService
    GetLoggerService() logger.LoggerService  // Not used
    GetCacheService() cache.CacheService     // Not used
}
```

### 2. Use Getter Methods
Always access services through getter methods, never directly:

```go
// Good
uc.services.GetAuditService().LogChange(...)

// Bad
uc.services.audit.LogChange(...)  // Breaks encapsulation
```

### 3. Nil-Safety in Mocks
Always provide default behavior in mock getters:

```go
func (m *MockServices) GetAuditService() audit.AuditService {
    if m.AuditService != nil {
        return m.AuditService  // Custom mock
    }
    return &MockAuditService{}  // Default mock
}
```

### 4. Document Service Responsibilities
Keep a clear understanding of what each service does:

```go
// AuditService logs all entity changes for compliance and debugging
func (s *Services) GetAuditService() audit.AuditService {
    return s.audit
}
```

## Related Documentation

- [Architecture](./ARCHITECTURE.md) - Overall system architecture
- [Testing](./TESTING.md) - Testing guidelines and patterns
- [Audit Logging](./SOFT_DELETES_AUDIT.md) - Audit service implementation

## References

- Dependency Injection: https://martinfowler.com/articles/injection.html
- Interface Segregation Principle: https://en.wikipedia.org/wiki/Interface_segregation_principle
