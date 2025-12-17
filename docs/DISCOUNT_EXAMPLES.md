# Enhanced Discount System - Practical Examples

## Overview

This document explains how the enhanced discount system works in practice when creating orders, with real-world scenarios showing how product, category, and user-based discounts are applied.

## How It Works

When a customer creates an order with a promo code, the system now performs **intelligent validation** instead of blindly applying the discount:

### 1. **Discount Validation**
- ✅ Checks if discount is active
- ✅ Validates date range (valid_from → valid_until)
- ✅ Checks global usage limits
- ✅ Validates minimum purchase requirements

### 2. **User Authorization**
- ✅ If discount is VIP/targeted → checks if user is in allowed list
- ✅ Validates per-user usage limits
- ✅ Tracks individual user usage

### 3. **Product/Category Matching**
- ✅ Checks which order items match discount criteria
- ✅ Applies discount only to eligible items
- ✅ Rejects if no items match (for specific discounts)

### 4. **Usage Tracking**
- ✅ Increments global usage counter
- ✅ Increments per-user usage counter
- ✅ Prevents over-use

---

## Practical Scenarios

### Scenario 1: Site-Wide Discount (OLD BEHAVIOR)

**Discount Setup:**
```json
{
  "promo_code": "WELCOME10",
  "discount_type": "percentage",
  "value": 10,
  "active": true,
  "products": [],      // Empty = applies to ALL
  "categories": [],    // Empty = applies to ALL
  "users": []          // Empty = any user can use
}
```

**Order Creation:**
```json
{
  "customer_id": 123,
  "products": [
    {"product_id": "laptop-uuid", "quantity": 1},  // $1000
    {"product_id": "mouse-uuid", "quantity": 1}    // $50
  ],
  "promo_code": "WELCOME10"
}
```

**Result:**
- ✅ Discount applies to entire order ($1050)
- 💰 Discount amount: $105 (10% of $1050)
- 💰 Final total: **$945**

---

### Scenario 2: Product-Specific Discount (NEW!)

**Discount Setup:**
```json
{
  "promo_code": "LAPTOP50",
  "discount_type": "amount",
  "value": 50,
  "active": true,
  "product_ids": ["laptop-uuid"],  // Only applies to laptops!
  "categories": [],
  "users": []
}
```

**Order Creation:**
```json
{
  "customer_id": 123,
  "products": [
    {"product_id": "laptop-uuid", "quantity": 1},  // $1000 ✅ ELIGIBLE
    {"product_id": "mouse-uuid", "quantity": 1}    // $50 ❌ NOT ELIGIBLE
  ],
  "promo_code": "LAPTOP50"
}
```

**Result:**
- ✅ Discount only applies to laptop ($1000)
- ❌ Mouse is excluded from discount
- 💰 Discount amount: $50 (only on laptop)
- 💰 Subtotal: Laptop: $950, Mouse: $50
- 💰 Final total: **$1000**

**What if no laptops in order?**
```json
{
  "customer_id": 123,
  "products": [
    {"product_id": "mouse-uuid", "quantity": 1},    // $50 ❌ NOT ELIGIBLE
    {"product_id": "keyboard-uuid", "quantity": 1}  // $80 ❌ NOT ELIGIBLE
  ],
  "promo_code": "LAPTOP50"
}
```
- ❌ **Order rejected**: "discount does not apply to any items in your order"

---

### Scenario 3: Category-Based Discount (NEW!)

**Discount Setup:**
```json
{
  "promo_code": "TECH15",
  "discount_type": "percentage",
  "value": 15,
  "active": true,
  "product_ids": [],
  "category_ids": ["electronics-uuid"],  // Only electronics!
  "users": []
}
```

**Order with Mixed Categories:**
```json
{
  "customer_id": 123,
  "products": [
    {"product_id": "laptop-uuid", "quantity": 1},     // $1000, Categories: [Electronics] ✅
    {"product_id": "headphones-uuid", "quantity": 1}, // $200, Categories: [Electronics] ✅
    {"product_id": "t-shirt-uuid", "quantity": 1}     // $30, Categories: [Clothing] ❌
  ],
  "promo_code": "TECH15"
}
```

**Result:**
- ✅ Discount applies to laptop + headphones ($1200)
- ❌ T-shirt excluded (different category)
- 💰 Discount amount: $180 (15% of $1200)
- 💰 Subtotal: Electronics: $1020, T-shirt: $30
- 💰 Final total: **$1050**

---

### Scenario 4: VIP User-Only Discount (NEW!)

**Discount Setup:**
```json
{
  "promo_code": "VIP20",
  "discount_type": "percentage",
  "value": 20,
  "active": true,
  "product_ids": [],
  "category_ids": [],
  "user_ids": [456, 789],  // Only specific VIP users!
  "user_usage_limit": 5     // Each user can use 5 times
}
```

**Regular User Tries to Use:**
```json
{
  "customer_id": 123,  // ❌ Not in VIP list
  "products": [
    {"product_id": "laptop-uuid", "quantity": 1}
  ],
  "promo_code": "VIP20"
}
```
- ❌ **Order rejected**: "this discount is not available for your account"

**VIP User Uses (First Time):**
```json
{
  "customer_id": 456,  // ✅ In VIP list
  "products": [
    {"product_id": "laptop-uuid", "quantity": 1}  // $1000
  ],
  "promo_code": "VIP20"
}
```
- ✅ Discount applies (user 456, usage 1/5)
- 💰 Discount: $200 (20% of $1000)
- 💰 Final total: **$800**

**VIP User After 5 Uses:**
```json
{
  "customer_id": 456,  // User already used 5 times
  "products": [
    {"product_id": "laptop-uuid", "quantity": 1}
  ],
  "promo_code": "VIP20"
}
```
- ❌ **Order rejected**: "you have reached the usage limit for this discount"

---

### Scenario 5: Targeted Product + VIP Users (NEW!)

**Discount Setup:**
```json
{
  "promo_code": "VIPLAPTOP",
  "discount_type": "amount",
  "value": 100,
  "active": true,
  "product_ids": ["laptop-uuid", "gaming-laptop-uuid"],
  "category_ids": [],
  "user_ids": [456, 789],  // Only VIP users
  "user_usage_limit": 1     // Once per user
}
```

**VIP User Orders Eligible Product:**
```json
{
  "customer_id": 456,  // ✅ VIP user
  "products": [
    {"product_id": "laptop-uuid", "quantity": 1},  // $1000 ✅ Eligible product
    {"product_id": "mouse-uuid", "quantity": 1}    // $50 ❌ Not eligible
  ],
  "promo_code": "VIPLAPTOP"
}
```
- ✅ User authorized ✅ Product eligible
- 💰 Discount: $100 (only on laptop)
- 💰 Final total: **$950**

**Non-VIP User with Eligible Product:**
```json
{
  "customer_id": 123,  // ❌ Not VIP
  "products": [
    {"product_id": "laptop-uuid", "quantity": 1}  // Product is eligible but...
  ],
  "promo_code": "VIPLAPTOP"
}
```
- ❌ **Order rejected**: "this discount is not available for your account"

**VIP User Orders Non-Eligible Product:**
```json
{
  "customer_id": 456,  // ✅ VIP user
  "products": [
    {"product_id": "mouse-uuid", "quantity": 1}  // ❌ Not in product list
  ],
  "promo_code": "VIPLAPTOP"
}
```
- ❌ **Order rejected**: "discount does not apply to any items in your order"

---

### Scenario 6: Global + Per-User Limits (NEW!)

**Discount Setup:**
```json
{
  "promo_code": "HOLIDAY50",
  "discount_type": "percentage",
  "value": 50,
  "active": true,
  "usage_limit": 100,      // Total 100 uses across all users
  "user_usage_limit": 1,   // Each user can only use once
  "product_ids": [],
  "category_ids": [],
  "user_ids": []
}
```

**First User's First Order:**
```json
{
  "customer_id": 123,
  "products": [
    {"product_id": "laptop-uuid", "quantity": 1}  // $1000
  ],
  "promo_code": "HOLIDAY50"
}
```
- ✅ Global usage: 1/100 ✅ User usage: 1/1
- 💰 Discount: $500
- 💰 Final total: **$500**

**Same User's Second Order:**
```json
{
  "customer_id": 123,  // Already used once
  "products": [
    {"product_id": "mouse-uuid", "quantity": 1}
  ],
  "promo_code": "HOLIDAY50"
}
```
- ❌ **Order rejected**: "you have reached the usage limit for this discount"

**After 100 Total Uses:**
- ✅ Global usage: 100/100
- ❌ **All orders rejected**: "discount is not valid or has reached usage limit"

---

### Scenario 7: Minimum Purchase + Date Range (NEW!)

**Discount Setup:**
```json
{
  "promo_code": "BIGSPENDER",
  "discount_type": "percentage",
  "value": 25,
  "active": true,
  "min_purchase_amount": 500,      // Minimum $500 order
  "max_discount_amount": 200,      // Cap at $200 discount
  "valid_from": "2025-12-01 00:00:00",
  "valid_until": "2025-12-31 23:59:59",
  "product_ids": [],
  "category_ids": [],
  "user_ids": []
}
```

**Small Order (Below Minimum):**
```json
{
  "customer_id": 123,
  "products": [
    {"product_id": "mouse-uuid", "quantity": 1}  // $50
  ],
  "promo_code": "BIGSPENDER"
}
```
- ❌ **Order rejected**: "minimum purchase amount not met"

**Medium Order (Qualifies):**
```json
{
  "customer_id": 123,
  "products": [
    {"product_id": "laptop-uuid", "quantity": 1}  // $1000
  ],
  "promo_code": "BIGSPENDER"
}
```
- ✅ Order total $1000 ≥ $500 minimum
- 💰 25% of $1000 = $250, but capped at $200
- 💰 Final total: **$800**

**Order Outside Date Range:**
```json
{
  "customer_id": 123,
  "products": [
    {"product_id": "laptop-uuid", "quantity": 1}
  ],
  "promo_code": "BIGSPENDER"
}
// If current date is 2026-01-01
```
- ❌ **Order rejected**: "discount is not valid, expired, or has reached usage limit"

---

## Key Differences: Old vs New System

### OLD System (Simple)
```
❌ Discount applies to EVERYTHING blindly
❌ No user restrictions
❌ No product/category filtering
❌ No usage limits
❌ No validation beyond "is code valid?"
```

### NEW System (Enhanced)
```
✅ Smart validation checks everything
✅ VIP/targeted user support
✅ Product-specific discounts
✅ Category-specific discounts
✅ Per-user usage tracking
✅ Global usage limits
✅ Minimum purchase requirements
✅ Maximum discount caps
✅ Date range validation
✅ Complex scenarios (VIP + Product + Date)
```

---

## Database Impact

When an order is created with a discount, the system:

1. **Reads from `discounts` table** (basic info)
2. **JOINs with `discount_products`** (if product-specific)
3. **JOINs with `discount_categories`** (if category-specific)
4. **JOINs with `discount_users`** (if user-specific)
5. **Checks `usage_count`** (global limit)
6. **Checks `discount_users.usage_count`** (per-user limit)
7. **Updates usage counters** on success

### Example Query Flow:
```sql
-- 1. Get discount with relationships
SELECT d.*, 
       dp.product_id, 
       dc.category_id, 
       du.user_id, du.usage_limit, du.usage_count
FROM discounts d
LEFT JOIN discount_products dp ON d.id = dp.discount_id
LEFT JOIN discount_categories dc ON d.id = dc.discount_id
LEFT JOIN discount_users du ON d.id = du.discount_id
WHERE d.promo_code = 'VIP20'
  AND d.active = true
  AND d.deleted_at IS NULL;

-- 2. Validate date range
-- Check if NOW() BETWEEN valid_from AND valid_until

-- 3. Check global usage
-- Check if usage_count < usage_limit

-- 4. For each order item, get product categories
SELECT pc.category_id
FROM products p
JOIN product_categories pc ON p.id = pc.product_id
WHERE p.id = 'laptop-uuid';

-- 5. Check if discount applies to product/category
-- Logic: Does product_id match discount_products?
--        Do any category_ids match discount_categories?

-- 6. Increment usage counters
UPDATE discounts SET usage_count = usage_count + 1 WHERE id = 'discount-uuid';
UPDATE discount_users SET usage_count = usage_count + 1 
WHERE discount_id = 'discount-uuid' AND user_id = 123;
```

---

## Testing the System

### Create Various Discounts:

```bash
# 1. Site-wide discount
curl -X POST http://localhost:8080/api/discounts \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "SITEWIDE10",
    "discount_type": "percentage",
    "value": 10,
    "active": true
  }'

# 2. Product-specific discount
curl -X POST http://localhost:8080/api/discounts \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "LAPTOP50",
    "discount_type": "amount",
    "value": 50,
    "active": true,
    "product_ids": ["<laptop-product-id>"]
  }'

# 3. Category discount
curl -X POST http://localhost:8080/api/discounts \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "ELECTRONICS15",
    "discount_type": "percentage",
    "value": 15,
    "active": true,
    "category_ids": ["<electronics-category-id>"]
  }'

# 4. VIP discount
curl -X POST http://localhost:8080/api/discounts \
  -H "Content-Type: application/json" \
  -d '{
    "promo_code": "VIP20",
    "discount_type": "percentage",
    "value": 20,
    "active": true,
    "user_ids": [2, 3],
    "user_usage_limit": 5
  }'
```

### Test Order Creation:

```bash
# Test with different scenarios
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "products": [
      {"product_id": "<laptop-id>", "quantity": 1},
      {"product_id": "<mouse-id>", "quantity": 1}
    ],
    "promo_code": "LAPTOP50"
  }'
```

---

## Future Enhancements

Possible additions to make it even more powerful:

1. **Stackable Discounts**: Apply multiple discounts to one order
2. **Buy X Get Y**: Quantity-based rules
3. **Bundle Discounts**: Discount when buying specific combinations
4. **Time-of-Day Discounts**: Happy hour pricing
5. **Referral Discounts**: Unique codes per user
6. **Tiered Discounts**: Bigger discount for higher spend

The current architecture supports adding these features easily!
