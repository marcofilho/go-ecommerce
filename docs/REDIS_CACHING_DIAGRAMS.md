# Redis Caching Flow Diagram

## Request Flow with Cache

```
┌──────────────┐
│   Client     │
└──────┬───────┘
       │ GET /api/products
       ▼
┌──────────────────────────────┐
│   ProductHandler             │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│ CachedProductRepository      │
│ (Decorator Pattern)          │
└──────┬───────────────────────┘
       │
       │ 1. Check Redis Cache
       ▼
┌──────────────────────────────┐      ┌─────────────┐
│      Redis Client            │─────▶│   Redis     │
└──────┬───────────────────────┘      │  (Cache)    │
       │                               └─────────────┘
       │
       ├─[Cache HIT]──────────────────┐
       │                               │
       │                          ┌────▼──────┐
       │                          │  Return   │
       │                          │  Cached   │
       │                          │   Data    │
       │                          └───────────┘
       │
       └─[Cache MISS]────┐
                         │
                         ▼
              ┌──────────────────────┐
              │ ProductRepository    │
              │ (PostgreSQL)         │
              └──────┬───────────────┘
                     │
                     │ 2. Query Database
                     ▼
              ┌──────────────────────┐
              │    PostgreSQL        │
              └──────┬───────────────┘
                     │
                     │ 3. Return Data
                     ▼
              ┌──────────────────────┐
              │  Store in Redis      │
              │  (TTL: 10 min)       │
              └──────┬───────────────┘
                     │
                     │ 4. Return to Client
                     ▼
              ┌──────────────────────┐
              │      Client          │
              └──────────────────────┘
```

## Cache Invalidation Flow

```
┌──────────────┐
│   Client     │
└──────┬───────┘
       │ PUT/DELETE /api/products/{id}
       ▼
┌──────────────────────────────┐
│   ProductHandler             │
│   (Admin Only)               │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│ CachedProductRepository      │
└──────┬───────────────────────┘
       │
       │ 1. Update Database
       ▼
┌──────────────────────────────┐
│  ProductRepositoryPostgres   │
└──────┬───────────────────────┘
       │
       │ 2. Database Updated
       ▼
┌──────────────────────────────┐
│ Cache Invalidation Logic     │
│                              │
│ • Delete: product:id:{uuid}  │
│ • Delete: products:list:*    │
└──────┬───────────────────────┘
       │
       │ 3. Clear Redis Keys
       ▼
┌──────────────────────────────┐      ┌─────────────┐
│      Redis Client            │─────▶│   Redis     │
│  DeleteByPattern("...")      │      │   KEYS      │
└──────┬───────────────────────┘      │  DELETED    │
       │                               └─────────────┘
       │
       │ 4. Success Response
       ▼
┌──────────────────────────────┐
│         Client               │
└──────────────────────────────┘
```

## Performance Comparison

### Without Cache (Cold Request)
```
Client ──[16ms]──▶ API ──[15ms]──▶ PostgreSQL ──[14ms]──▶ Data ──[1ms]──▶ Client
Total: ~16-20ms
```

### With Cache (Warm Request)
```
Client ──[3ms]──▶ API ──[2ms]──▶ Redis ──[1ms]──▶ Data ──▶ Client
Total: ~2-5ms (70-80% faster!)
```

## Cache Key Structure

```
product:id:{uuid}
├─ Example: product:id:d4444444-4444-4444-4444-444444444444
└─ TTL: 10 minutes

products:list:page:{page}:size:{size}:instock:{bool}
├─ Example: products:list:page:1:size:10:instock:false
└─ TTL: 10 minutes
```

## Architecture Layers

```
┌─────────────────────────────────────────────────────┐
│                  HTTP Layer                         │
│  (handlers, middleware, routes)                     │
└───────────────────────┬─────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│              Use Case Layer                         │
│  (business logic, orchestration)                    │
└───────────────────────┬─────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│         Repository Layer (Interface)                │
└───────────────────────┬─────────────────────────────┘
                        │
            ┌───────────┴───────────┐
            ▼                       ▼
┌────────────────────────┐  ┌──────────────────┐
│ CachedProductRepository│  │  RedisClient     │
│    (Decorator)         │  │  (Cache)         │
└───────────┬────────────┘  └──────────────────┘
            │
            ▼
┌────────────────────────┐
│ ProductRepositoryPG    │
│    (Database)          │
└────────────────────────┘
```
