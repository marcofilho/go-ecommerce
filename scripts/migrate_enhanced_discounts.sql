-- Enhanced Discount System Migration
-- This migration adds support for:
-- - Product-specific discounts
-- - Category-wide discounts
-- - User-specific discounts
-- - Enhanced discount constraints (min purchase, max discount, usage limits, date ranges)

-- ============================================
-- Step 1: Add new columns to discounts table
-- ============================================

ALTER TABLE discounts
ADD COLUMN IF NOT EXISTS min_purchase_amount DECIMAL(10, 2),
ADD COLUMN IF NOT EXISTS max_discount_amount DECIMAL(10, 2),
ADD COLUMN IF NOT EXISTS usage_limit INTEGER,
ADD COLUMN IF NOT EXISTS usage_count INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS valid_from TIMESTAMP,
ADD COLUMN IF NOT EXISTS valid_until TIMESTAMP,
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- Add index on deleted_at for soft deletes
CREATE INDEX IF NOT EXISTS idx_discounts_deleted_at ON discounts (deleted_at);

-- ============================================
-- Step 2: Create junction tables
-- ============================================

-- Junction table: Discount <-> Products
CREATE TABLE IF NOT EXISTS discount_products (
    discount_id UUID NOT NULL,
    product_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (discount_id, product_id),
    CONSTRAINT fk_discount_products_discount FOREIGN KEY (discount_id) REFERENCES discounts (id) ON DELETE CASCADE,
    CONSTRAINT fk_discount_products_product FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_discount_products_discount ON discount_products (discount_id);

CREATE INDEX IF NOT EXISTS idx_discount_products_product ON discount_products (product_id);

-- Junction table: Discount <-> Categories
CREATE TABLE IF NOT EXISTS discount_categories (
    discount_id UUID NOT NULL,
    category_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (discount_id, category_id),
    CONSTRAINT fk_discount_categories_discount FOREIGN KEY (discount_id) REFERENCES discounts (id) ON DELETE CASCADE,
    CONSTRAINT fk_discount_categories_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_discount_categories_discount ON discount_categories (discount_id);

CREATE INDEX IF NOT EXISTS idx_discount_categories_category ON discount_categories (category_id);

-- Junction table: Discount <-> Users (with per-user usage tracking)
CREATE TABLE IF NOT EXISTS discount_users (
    discount_id UUID NOT NULL,
    user_id UUID NOT NULL,
    usage_count INTEGER DEFAULT 0,
    usage_limit INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (discount_id, user_id),
    CONSTRAINT fk_discount_users_discount FOREIGN KEY (discount_id) REFERENCES discounts (id) ON DELETE CASCADE,
    CONSTRAINT fk_discount_users_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_discount_users_discount ON discount_users (discount_id);

CREATE INDEX IF NOT EXISTS idx_discount_users_user ON discount_users (user_id);

-- ============================================
-- Step 3: Add comments for documentation
-- ============================================

COMMENT ON
TABLE discount_products IS 'Junction table linking discounts to specific products';

COMMENT ON
TABLE discount_categories IS 'Junction table linking discounts to product categories';

COMMENT ON
TABLE discount_users IS 'Junction table linking discounts to specific users with per-user usage limits';

COMMENT ON COLUMN discounts.min_purchase_amount IS 'Minimum order total required to apply this discount';

COMMENT ON COLUMN discounts.max_discount_amount IS 'Maximum discount amount (caps percentage discounts)';

COMMENT ON COLUMN discounts.usage_limit IS 'Global limit on how many times this discount can be used';

COMMENT ON COLUMN discounts.usage_count IS 'Current number of times this discount has been used';

COMMENT ON COLUMN discounts.valid_from IS 'Discount becomes active at this date/time';

COMMENT ON COLUMN discounts.valid_until IS 'Discount expires at this date/time';

-- ============================================
-- Step 4: Sample data for testing (optional)
-- ============================================

-- Example 1: Site-wide 10% discount for all users
DO $$
DECLARE
    discount_id UUID := 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
BEGIN
    INSERT INTO discounts (id, promo_code, discount_type, value, active, min_purchase_amount, valid_from, valid_until)
    VALUES (discount_id, 'WELCOME10', 'percentage', 10, true, 50.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '30 days')
    ON CONFLICT (id) DO NOTHING;
    -- No products/categories/users = applies to all
END $$;

-- Example 2: $20 off specific product for VIP users
DO $$
DECLARE
    discount_id UUID := 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb';
    macbook_id UUID := 'd4444444-4444-4444-4444-444444444444';
    admin_user UUID;
BEGIN
    -- Get first admin user
    SELECT id INTO admin_user FROM users WHERE role = 'admin' LIMIT 1;
    
    IF admin_user IS NOT NULL THEN
        INSERT INTO discounts (id, promo_code, discount_type, value, active, usage_limit, valid_until)
        VALUES (discount_id, 'VIP20', 'amount', 20, true, 100, CURRENT_TIMESTAMP + INTERVAL '60 days')
        ON CONFLICT (id) DO NOTHING;
        
        -- Link to specific product
        INSERT INTO discount_products (discount_id, product_id)
        VALUES (discount_id, macbook_id)
        ON CONFLICT DO NOTHING;
        
        -- Link to specific user
        INSERT INTO discount_users (discount_id, user_id, usage_limit)
        VALUES (discount_id, admin_user, 5)
        ON CONFLICT DO NOTHING;
    END IF;
END $$;

-- Example 3: 15% off entire Electronics category
DO $$
DECLARE
    discount_id UUID := 'cccccccc-cccc-cccc-cccc-cccccccccccc';
    electronics_id UUID := 'a1111111-1111-1111-1111-111111111111';
BEGIN
    INSERT INTO discounts (id, promo_code, discount_type, value, active, max_discount_amount, valid_until)
    VALUES (discount_id, 'TECH15', 'percentage', 15, true, 100.00, CURRENT_TIMESTAMP + INTERVAL '90 days')
    ON CONFLICT (id) DO NOTHING;
    
    -- Link to Electronics category
    INSERT INTO discount_categories (discount_id, category_id)
    VALUES (discount_id, electronics_id)
    ON CONFLICT DO NOTHING;
END $$;

-- ============================================
-- Verification queries
-- ============================================

-- Show all discounts with their associations
DO $$
BEGIN
    RAISE NOTICE '=== Discount System Migration Complete ===';
    RAISE NOTICE 'Tables created:';
    RAISE NOTICE '  - discount_products';
    RAISE NOTICE '  - discount_categories';
    RAISE NOTICE '  - discount_users';
    RAISE NOTICE '';
    RAISE NOTICE 'Enhanced columns added to discounts table';
    RAISE NOTICE 'Sample discounts created (if data exists)';
END $$;