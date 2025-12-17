# Enhanced Discount System Design

## Problem Statement

We need a flexible discount system that can apply discounts based on:
- **Specific Products** - Discount applies to certain products
- **Categories** - Discount applies to all products in a category
- **Users** - Discount applies to specific users (or all users if not specified)
- **Combinations** - A discount can apply to multiple products, categories, and users simultaneously

## Database Schema

### Core Tables

```sql
-- Main discount table (existing, enhanced)
CREATE TABLE discounts (
    id UUID PRIMARY KEY,
    promo_code VARCHAR(50) UNIQUE NOT NULL,
    discount_type VARCHAR(20) NOT NULL DEFAULT 'amount', -- 'percentage' or 'amount'
    value DECIMAL(10,2) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    min_purchase_amount DECIMAL(10,2), -- Optional: minimum purchase to apply discount
    max_discount_amount DECIMAL(10,2), -- Optional: cap for percentage discounts
    usage_limit INTEGER, -- Optional: max total uses
    usage_count INTEGER DEFAULT 0, -- Track how many times used
    valid_from TIMESTAMP, -- Optional: start date
    valid_until TIMESTAMP, -- Optional: end date
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Junction table: Discount <-> Products
CREATE TABLE discount_products (
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (discount_id, product_id)
);
CREATE INDEX idx_discount_products_discount ON discount_products(discount_id);
CREATE INDEX idx_discount_products_product ON discount_products(product_id);

-- Junction table: Discount <-> Categories
CREATE TABLE discount_categories (
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (discount_id, category_id)
);
CREATE INDEX idx_discount_categories_discount ON discount_categories(discount_id);
CREATE INDEX idx_discount_categories_category ON discount_categories(category_id);

-- Junction table: Discount <-> Users
CREATE TABLE discount_users (
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    usage_count INTEGER DEFAULT 0, -- Track per-user usage
    usage_limit INTEGER, -- Optional: per-user limit
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (discount_id, user_id)
);
CREATE INDEX idx_discount_users_discount ON discount_users(discount_id);
CREATE INDEX idx_discount_users_user ON discount_users(user_id);
```

## Entity Models

### Enhanced Discount Entity

```go
package entity

import (
    "errors"
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type DiscountType string

const (
    Percentage DiscountType = "percentage"
    Amount     DiscountType = "amount"
)

type Discount struct {
    ID                 uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
    PromoCode          string         `gorm:"uniqueIndex;not null" json:"promo_code"`
    DiscountType       DiscountType   `gorm:"type:varchar(20);not null;default:'amount'" json:"discount_type"`
    Value              float64        `gorm:"type:decimal(10,2);not null" json:"value"`
    Active             bool           `gorm:"not null;default:true" json:"active"`
    MinPurchaseAmount  *float64       `gorm:"type:decimal(10,2)" json:"min_purchase_amount,omitempty"`
    MaxDiscountAmount  *float64       `gorm:"type:decimal(10,2)" json:"max_discount_amount,omitempty"`
    UsageLimit         *int           `json:"usage_limit,omitempty"`
    UsageCount         int            `gorm:"default:0" json:"usage_count"`
    ValidFrom          *time.Time     `json:"valid_from,omitempty"`
    ValidUntil         *time.Time     `json:"valid_until,omitempty"`
    CreatedAt          time.Time      `json:"created_at"`
    UpdatedAt          time.Time      `json:"updated_at"`
    DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
    
    // Relationships
    Products           []Product      `gorm:"many2many:discount_products;" json:"products,omitempty"`
    Categories         []Category     `gorm:"many2many:discount_categories;" json:"categories,omitempty"`
    Users              []User         `gorm:"many2many:discount_users;" json:"users,omitempty"`
}

func (d *Discount) BeforeCreate(tx *gorm.DB) error {
    if d.ID == uuid.Nil {
        d.ID = uuid.New()
    }
    return nil
}

func (d *Discount) Validate() error {
    if d.PromoCode == "" {
        return errors.New("promo code is required")
    }
    if d.Value <= 0 {
        return errors.New("discount value must be positive")
    }
    if d.DiscountType == Percentage && d.Value > 100 {
        return errors.New("percentage discount cannot exceed 100%")
    }
    if d.ValidFrom != nil && d.ValidUntil != nil && d.ValidFrom.After(*d.ValidUntil) {
        return errors.New("valid_from must be before valid_until")
    }
    return nil
}

func (d *Discount) IsValid() bool {
    if !d.Active {
        return false
    }
    
    now := time.Now()
    
    if d.ValidFrom != nil && now.Before(*d.ValidFrom) {
        return false
    }
    
    if d.ValidUntil != nil && now.After(*d.ValidUntil) {
        return false
    }
    
    if d.UsageLimit != nil && d.UsageCount >= *d.UsageLimit {
        return false
    }
    
    return true
}

func (d *Discount) CanBeUsedBy(userID uuid.UUID, userUsageCount int, userUsageLimit *int) bool {
    if !d.IsValid() {
        return false
    }
    
    if userUsageLimit != nil && userUsageCount >= *userUsageLimit {
        return false
    }
    
    return true
}

func (d *Discount) AppliesTo(productID uuid.UUID, categoryIDs []uuid.UUID) bool {
    // If no specific products or categories, applies to all
    if len(d.Products) == 0 && len(d.Categories) == 0 {
        return true
    }
    
    // Check if product is in discount's product list
    for _, p := range d.Products {
        if p.ID == productID {
            return true
        }
    }
    
    // Check if product's category is in discount's category list
    for _, dc := range d.Categories {
        for _, catID := range categoryIDs {
            if dc.ID == catID {
                return true
            }
        }
    }
    
    return false
}

func (d *Discount) ApplyDiscount(total float64) (float64, error) {
    if d.MinPurchaseAmount != nil && total < *d.MinPurchaseAmount {
        return total, errors.New("minimum purchase amount not met")
    }
    
    var discountAmount float64
    
    switch d.DiscountType {
    case Percentage:
        discountAmount = total * (d.Value / 100)
        if d.MaxDiscountAmount != nil && discountAmount > *d.MaxDiscountAmount {
            discountAmount = *d.MaxDiscountAmount
        }
    case Amount:
        discountAmount = d.Value
    default:
        return total, errors.New("invalid discount type")
    }
    
    newTotal := total - discountAmount
    if newTotal < 0 {
        newTotal = 0
    }
    
    return newTotal, nil
}
```

### Junction Table Entities

```go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type DiscountUser struct {
    DiscountID  uuid.UUID  `gorm:"primaryKey;type:uuid"`
    UserID      uuid.UUID  `gorm:"primaryKey;type:uuid"`
    UsageCount  int        `gorm:"default:0"`
    UsageLimit  *int       `gorm:"type:integer"`
    CreatedAt   time.Time
    
    // Relationships
    Discount    Discount   `gorm:"foreignKey:DiscountID"`
    User        User       `gorm:"foreignKey:UserID"`
}

func (DiscountUser) TableName() string {
    return "discount_users"
}
```

## Use Cases Examples

### Scenario 1: Category-wide discount for all users
```
Discount: "ELECTRONICS20"
- Type: Percentage (20%)
- Categories: [Electronics]
- Products: [] (empty)
- Users: [] (empty - applies to all users)
- Result: 20% off all electronics for everyone
```

### Scenario 2: User-specific discount on specific products
```
Discount: "VIP50"
- Type: Amount ($50)
- Categories: [] (empty)
- Products: [MacBook Pro, iPhone 15]
- Users: [user-123, user-456]
- Result: $50 off MacBook/iPhone only for VIP users
```

### Scenario 3: Mixed discount
```
Discount: "SUMMER30"
- Type: Percentage (30%)
- Categories: [Summer Collection, Beach Wear]
- Products: [Special Item A, Special Item B]
- Users: [premium-user-789]
- Result: 30% off for premium user on:
  - All products in Summer Collection
  - All products in Beach Wear
  - Plus Special Item A and B (even if not in those categories)
```

## Validation Logic

### When applying a discount to an order:

```go
func (uc *DiscountUseCase) ValidateDiscount(
    ctx context.Context,
    promoCode string,
    userID uuid.UUID,
    orderItems []OrderItem,
) (*Discount, error) {
    // 1. Get discount with all relationships
    discount, err := uc.repo.GetByPromoCodeWithRelations(ctx, promoCode)
    if err != nil {
        return nil, errors.New("invalid promo code")
    }
    
    // 2. Check if discount is valid (active, dates, usage limits)
    if !discount.IsValid() {
        return nil, errors.New("discount is not valid or has expired")
    }
    
    // 3. Check if user is eligible
    userUsage, err := uc.repo.GetUserUsage(ctx, discount.ID, userID)
    if err == nil {
        if !discount.CanBeUsedBy(userID, userUsage.UsageCount, userUsage.UsageLimit) {
            return nil, errors.New("user has reached usage limit for this discount")
        }
    }
    
    // 4. Check if discount applies to at least one item in the order
    hasApplicableItem := false
    for _, item := range orderItems {
        product, _ := uc.productRepo.GetByID(ctx, item.ProductID)
        categoryIDs := extractCategoryIDs(product.Categories)
        
        if discount.AppliesTo(item.ProductID, categoryIDs) {
            hasApplicableItem = true
            break
        }
    }
    
    if !hasApplicableItem {
        return nil, errors.New("discount does not apply to any items in your order")
    }
    
    return discount, nil
}
```

## API Endpoints

### Create Discount (Admin)
```http
POST /api/admin/discounts
{
  "promo_code": "SUMMER30",
  "discount_type": "percentage",
  "value": 30,
  "min_purchase_amount": 100,
  "valid_from": "2025-06-01T00:00:00Z",
  "valid_until": "2025-08-31T23:59:59Z",
  "product_ids": ["uuid1", "uuid2"],
  "category_ids": ["uuid3"],
  "user_ids": ["uuid4", "uuid5"]
}
```

### Validate Discount (Customer)
```http
POST /api/discounts/validate
{
  "promo_code": "SUMMER30",
  "order_items": [
    {"product_id": "uuid1", "quantity": 2},
    {"product_id": "uuid2", "quantity": 1}
  ]
}

Response:
{
  "valid": true,
  "discount": {
    "id": "...",
    "promo_code": "SUMMER30",
    "discount_type": "percentage",
    "value": 30
  },
  "applicable_items": ["uuid1"],
  "estimated_discount": 45.00
}
```

## Benefits of This Design

✅ **Flexibility**: A discount can target products, categories, users, or combinations  
✅ **Scalability**: Easy to add more products/categories/users to existing discounts  
✅ **Performance**: Proper indexes on junction tables for fast queries  
✅ **Maintainability**: Clear separation of concerns with junction tables  
✅ **Business Logic**: Support for complex scenarios (VIP users, seasonal sales, product-specific promos)  
✅ **Audit Trail**: Track usage per user and globally  
✅ **Constraints**: Min purchase, max discount, usage limits, date ranges  

## Migration Strategy

1. Create new junction tables
2. Add new columns to discounts table
3. Migrate existing discount data if any
4. Update entities and repositories
5. Update use cases with validation logic
6. Update handlers and DTOs
7. Test thoroughly with various scenarios
