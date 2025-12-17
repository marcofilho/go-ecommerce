# Enhanced Discount System - Order Flow

## Visual Flow Diagram

```
                          CREATE ORDER REQUEST
                                  │
                                  ▼
                    ┌─────────────────────────┐
                    │  Parse Products & Items  │
                    └─────────────────────────┘
                                  │
                                  ▼
                    ┌─────────────────────────┐
                    │   Validate Stock         │
                    │   Calculate Item Totals  │
                    └─────────────────────────┘
                                  │
                                  ▼
                    ┌─────────────────────────┐
                    │   Calculate Order Total  │
                    └─────────────────────────┘
                                  │
                                  ▼
                         Has Promo Code?
                         /           \
                        NO           YES
                        │             │
                        │             ▼
                        │    ┌──────────────────────────────┐
                        │    │ ENHANCED DISCOUNT VALIDATION │
                        │    └──────────────────────────────┘
                        │             │
                        │             ▼
                        │    ┌──────────────────────────────┐
                        │    │ 1. Get Discount with         │
                        │    │    Relations (products,       │
                        │    │    categories, users)         │
                        │    └──────────────────────────────┘
                        │             │
                        │             ▼
                        │    ┌──────────────────────────────┐
                        │    │ 2. Validate Discount         │
                        │    │    ✓ Is active?              │
                        │    │    ✓ Within date range?      │
                        │    │    ✓ Not expired?            │
                        │    │    ✓ Has remaining uses?     │
                        │    └──────────────────────────────┘
                        │             │
                        │             ▼
                        │         Valid?
                        │         /    \
                        │       NO      YES
                        │        │       │
                        │        │       ▼
                        │        │  ┌──────────────────────────────┐
                        │        │  │ 3. Check User Authorization  │
                        │        │  │    Has user restrictions?     │
                        │        │  └──────────────────────────────┘
                        │        │       │
                        │        │       ▼
                        │        │  User in List?
                        │        │  /          \
                        │        │ NO          YES
                        │        │  │           │
                        │        │  │           ▼
                        │        │  │  ┌──────────────────────────┐
                        │        │  │  │ Check per-user limits    │
                        │        │  │  └──────────────────────────┘
                        │        │  │           │
                        │        │  │           ▼
                        │        │  │      Within limit?
                        │        │  │      /          \
                        │        │  │    NO            YES
                        │        │  │     │             │
                        │        │  │     │             ▼
                        │        │  │     │  ┌──────────────────────────────┐
                        │        │  │     │  │ 4. Match Products/Categories │
                        │        │  │     │  │    Check each order item     │
                        │        │  │     │  └──────────────────────────────┘
                        │        │  │     │             │
                        │        │  │     │             ▼
                        │        │  │     │    Any items match?
                        │        │  │     │    /              \
                        │        │  │     │  NO               YES
                        │        │  │     │   │                │
                        │        └──┼─────┼───┼────────────────┤
                        │           │     │   │                │
                        │           ▼     ▼   ▼                ▼
                        │    ┌─────────────────────┐  ┌──────────────────┐
                        │    │   REJECT ORDER      │  │ 5. Calculate     │
                        │    │   Return Error:     │  │    Discount      │
                        │    │   - Invalid code    │  └──────────────────┘
                        │    │   - Expired         │           │
                        │    │   - No access       │           ▼
                        │    │   - No items match  │  ┌──────────────────┐
                        │    │   - Usage exceeded  │  │ Site-wide or     │
                        │    └─────────────────────┘  │ Specific items?  │
                        │                             └──────────────────┘
                        │                                      │
                        │                            ┌─────────┴─────────┐
                        │                            │                   │
                        │                       Site-wide           Specific
                        │                            │                   │
                        │                            ▼                   ▼
                        │              ┌──────────────────┐  ┌──────────────────┐
                        │              │ Apply to entire  │  │ Apply only to    │
                        │              │ order total      │  │ matching items   │
                        │              └──────────────────┘  └──────────────────┘
                        │                            │                   │
                        │                            └─────────┬─────────┘
                        │                                      │
                        │                                      ▼
                        │                         ┌──────────────────────────┐
                        │                         │ 6. Update Usage Counters │
                        │                         │    ✓ Global counter++    │
                        │                         │    ✓ User counter++      │
                        │                         └──────────────────────────┘
                        │                                      │
                        └──────────────────────────────────────┘
                                                  │
                                                  ▼
                                    ┌─────────────────────────┐
                                    │   Save Order to DB      │
                                    └─────────────────────────┘
                                                  │
                                                  ▼
                                    ┌─────────────────────────┐
                                    │   Return Order Response │
                                    │   with applied discount │
                                    └─────────────────────────┘
```

## Key Decision Points

### 1. **Discount Validation**
```
IsValid() checks:
- active == true
- deleted_at == NULL
- NOW() >= valid_from (if set)
- NOW() <= valid_until (if set)
- usage_count < usage_limit (if set)
```

### 2. **User Authorization** (if discount has user restrictions)
```
For each discount.Users:
  If user.ID == order.customer_id:
    - User is authorized
    - Check user_usage_count < user_usage_limit
```

### 3. **Product/Category Matching**
```
For each order item:
  Get product with categories
  
  AppliesTo(product_id, category_ids):
    - If discount.Products is empty AND discount.Categories is empty:
        Return TRUE (site-wide)
    
    - If product_id in discount.Products:
        Return TRUE
    
    - If any category_id in discount.Categories:
        Return TRUE
    
    Return FALSE
```

### 4. **Discount Calculation**
```
If discount applies to specific items:
  applicable_total = sum of matching items
  discount_amount = Calculate(applicable_total)
  final_total = order_total - discount_amount

Else (site-wide):
  discount_amount = Calculate(order_total)
  final_total = discount_amount
```

## Database Queries Involved

```sql
-- 1. Get discount with relationships
SELECT d.*, 
       array_agg(DISTINCT dp.product_id) as product_ids,
       array_agg(DISTINCT dc.category_id) as category_ids,
       array_agg(DISTINCT du.user_id) as user_ids
FROM discounts d
LEFT JOIN discount_products dp ON d.id = dp.discount_id
LEFT JOIN discount_categories dc ON d.id = dc.discount_id
LEFT JOIN discount_users du ON d.id = du.discount_id
WHERE d.promo_code = $1
  AND d.active = true
  AND d.deleted_at IS NULL
GROUP BY d.id;

-- 2. Check user-specific usage
SELECT usage_count, usage_limit
FROM discount_users
WHERE discount_id = $1 AND user_id = $2;

-- 3. Get product with categories
SELECT p.*, array_agg(c.id) as category_ids
FROM products p
LEFT JOIN product_categories pc ON p.id = pc.product_id
LEFT JOIN categories c ON pc.category_id = c.id
WHERE p.id = $1
GROUP BY p.id;

-- 4. Increment usage counters
UPDATE discounts 
SET usage_count = usage_count + 1
WHERE id = $1;

UPDATE discount_users 
SET usage_count = usage_count + 1
WHERE discount_id = $1 AND user_id = $2;
```

## Performance Considerations

### Before (Simple System):
- 1 query: Get discount by promo code
- Total: **1 DB query**

### After (Enhanced System):
- 1 query: Get discount with relationships (3 JOINs)
- 1 query: Get user usage (if user-specific)
- N queries: Get product categories (1 per order item)
- 2 queries: Update usage counters
- Total: **4 + N DB queries**

### Optimization Opportunities:
1. **Cache product categories** (Redis) - reduces N queries
2. **Cache discount relationships** - reduces JOINs
3. **Batch category lookups** - single query for all products
4. **Background usage updates** - async counter increments

## Code Location

The enhanced logic is in:
- [src/usecase/order/order_usecase.go](../src/usecase/order/order_usecase.go) - Line 148-235
- [src/internal/domain/entity/discount.go](../src/internal/domain/entity/discount.go) - Methods: `IsValid()`, `CanBeUsedBy()`, `AppliesTo()`, `ApplyDiscount()`
- [src/internal/infrastructure/repository/discount_repository_postgres.go](../src/internal/infrastructure/repository/discount_repository_postgres.go) - `GetByPromoCodeWithRelations()`
