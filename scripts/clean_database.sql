-- ============================================
-- Go E-Commerce - Database Cleanup Script
-- ============================================
-- WARNING: This script will DELETE ALL DATA from the database!
-- Use with caution, especially in production environments.
--
-- Usage:
--   psql -U postgres -d ecommerce -f scripts/clean_database.sql
-- Or:
--   docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/clean_database.sql
--
-- This script:
-- 1. Disables foreign key checks temporarily (if needed)
-- 2. Truncates all tables in the correct order (respecting foreign keys)
-- 3. Resets sequences for auto-incrementing columns
-- 4. Re-enables foreign key checks
-- ============================================

BEGIN;

-- Display warning message
DO $$
BEGIN
  RAISE NOTICE '========================================';
  RAISE NOTICE 'WARNING: Database Cleanup in Progress';
  RAISE NOTICE '========================================';
  RAISE NOTICE 'All data will be permanently deleted!';
  RAISE NOTICE '========================================';
END $$;

-- ============================================
-- TRUNCATE ALL TABLES
-- ============================================
-- Using CASCADE to handle foreign key constraints
-- Tables are truncated in reverse dependency order for clarity

-- Step 1: Delete dependent data first (child tables)
TRUNCATE TABLE webhook_logs CASCADE;

TRUNCATE TABLE order_items CASCADE;

TRUNCATE TABLE orders CASCADE;

TRUNCATE TABLE product_categories CASCADE;

TRUNCATE TABLE product_variants CASCADE;

-- Step 2: Delete main entity tables
TRUNCATE TABLE products CASCADE;

TRUNCATE TABLE categories CASCADE;

TRUNCATE TABLE users CASCADE;

-- ============================================
-- RESET SEQUENCES
-- ============================================
-- Note: Users table uses UUID, no auto-increment sequence to reset

-- ============================================
-- VERIFICATION
-- ============================================
-- Check that all tables are empty
DO $$ DECLARE user_count INT;

category_count INT;

product_count INT;

variant_count INT;

pc_count INT;

order_count INT;

item_count INT;

webhook_count INT;

total_records INT;

BEGIN
SELECT COUNT(*) INTO user_count
FROM users;

SELECT COUNT(*) INTO category_count
FROM categories;

SELECT COUNT(*) INTO product_count
FROM products;

SELECT COUNT(*) INTO variant_count
FROM product_variants;

SELECT COUNT(*) INTO pc_count
FROM product_categories;

SELECT COUNT(*) INTO order_count
FROM orders;

SELECT COUNT(*) INTO item_count
FROM order_items;

SELECT COUNT(*) INTO webhook_count
FROM webhook_logs;

total_records := user_count + category_count + product_count + variant_count + pc_count + order_count + item_count + webhook_count;

RAISE NOTICE '========================================';

RAISE NOTICE 'Database Cleanup Complete!';

RAISE NOTICE '========================================';

RAISE NOTICE 'Remaining Records:';

RAISE NOTICE '  Users:              %',
user_count;

RAISE NOTICE '  Categories:         %',
category_count;

RAISE NOTICE '  Products:           %',
product_count;

RAISE NOTICE '  Product Variants:   %',
variant_count;

RAISE NOTICE '  Product-Categories: %',
pc_count;

RAISE NOTICE '  Orders:             %',
order_count;

RAISE NOTICE '  Order Items:        %',
item_count;

RAISE NOTICE '  Webhook Logs:       %',
webhook_count;

RAISE NOTICE '----------------------------------------';

RAISE NOTICE '  TOTAL:              %',
total_records;

RAISE NOTICE '========================================';

IF total_records = 0 THEN RAISE NOTICE 'SUCCESS: All tables are empty!';

ELSE RAISE WARNING 'WARNING: Some tables still contain data!';

END IF;

RAISE NOTICE '========================================';

END $$;

COMMIT;

-- ============================================
-- OPTIONAL: Vacuum and Analyze
-- ============================================
-- Uncomment the following lines to optimize the database after cleanup
-- This reclaims storage and updates statistics
-- VACUUM ANALYZE users;
-- VACUUM ANALYZE categories;
-- VACUUM ANALYZE products;
-- VACUUM ANALYZE product_variants;
-- VACUUM ANALYZE product_categories;
-- VACUUM ANALYZE orders;
-- VACUUM ANALYZE order_items;
-- VACUUM ANALYZE webhook_logs;