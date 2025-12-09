-- ============================================
-- Go E-Commerce - Database Seed Script
-- ============================================
-- This script populates the database with initial sample data
-- Usage: docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/seed_data.sql

BEGIN;

-- ============================================
-- 1. USERS (3 records) - UUID, no timestamps
-- ============================================
INSERT INTO
    users (
        id,
        email,
        password_hash,
        name,
        role,
        active
    )
VALUES (
        '11111111-1111-1111-1111-111111111111',
        'admin@ecommerce.com',
        '$2a$10$N9qo8uLOickgx2ZMRZoMye7I925HvEMhLmYh8HWHlKWkflGfwYvAe',
        'System Administrator',
        'admin',
        true
    ),
    (
        '22222222-2222-2222-2222-222222222222',
        'john.doe@example.com',
        '$2a$10$N9qo8uLOickgx2ZMRZoMye7I925HvEMhLmYh8HWHlKWkflGfwYvAe',
        'John Doe',
        'customer',
        true
    ),
    (
        '33333333-3333-3333-3333-333333333333',
        'jane.smith@example.com',
        '$2a$10$N9qo8uLOickgx2ZMRZoMye7I925HvEMhLmYh8HWHlKWkflGfwYvAe',
        'Jane Smith',
        'customer',
        true
    ) ON CONFLICT (email) DO NOTHING;

-- ============================================
-- 2. CATEGORIES (3 records)
-- ============================================
INSERT INTO
    categories (id, name)
VALUES (
        'a1111111-1111-1111-1111-111111111111',
        'Electronics'
    ),
    (
        'b2222222-2222-2222-2222-222222222222',
        'Computers'
    ),
    (
        'c3333333-3333-3333-3333-333333333333',
        'Gaming'
    ) ON CONFLICT (name) DO NOTHING;

-- ============================================
-- 3. PRODUCTS (3 records)
-- ============================================
INSERT INTO
    products (
        id,
        name,
        description,
        price,
        quantity
    )
VALUES (
        'd4444444-4444-4444-4444-444444444444',
        'MacBook Pro 16" M3 Max',
        'Powerful laptop with M3 Max chip, 36GB RAM, 1TB SSD. Perfect for developers and creative professionals.',
        3499.00,
        25
    ),
    (
        'e5555555-5555-5555-5555-555555555555',
        'Logitech MX Master 3S Mouse',
        'Advanced wireless mouse with 8K DPI sensor. Rechargeable battery lasts 70 days.',
        99.99,
        150
    ),
    (
        'f6666666-6666-6666-6666-666666666666',
        'PlayStation 5 Digital Edition',
        'Next-gen console with 825GB SSD, ray tracing, 4K gaming at 120fps.',
        449.99,
        40
    ) ON CONFLICT (id) DO NOTHING;

-- ============================================
-- 4. PRODUCT VARIANTS (3 records)
-- ============================================
INSERT INTO
    product_variants (
        id,
        product_id,
        variant_name,
        variant_value,
        quantity,
        price_override
    )
VALUES (
        '77777777-7777-7777-7777-777777777777',
        'd4444444-4444-4444-4444-444444444444',
        'Color',
        'Space Black',
        15,
        NULL
    ),
    (
        '88888888-8888-8888-8888-888888888888',
        'd4444444-4444-4444-4444-444444444444',
        'Color',
        'Silver',
        10,
        NULL
    ),
    (
        '99999999-9999-9999-9999-999999999999',
        'e5555555-5555-5555-5555-555555555555',
        'Color',
        'Graphite',
        80,
        NULL
    ) ON CONFLICT (id) DO NOTHING;

-- ============================================
-- 5. PRODUCT_CATEGORIES (Junction - 6 records)
-- ============================================
INSERT INTO
    product_categories (product_id, category_id)
VALUES (
        'd4444444-4444-4444-4444-444444444444',
        'a1111111-1111-1111-1111-111111111111'
    ),
    (
        'd4444444-4444-4444-4444-444444444444',
        'b2222222-2222-2222-2222-222222222222'
    ),
    (
        'e5555555-5555-5555-5555-555555555555',
        'a1111111-1111-1111-1111-111111111111'
    ),
    (
        'e5555555-5555-5555-5555-555555555555',
        'b2222222-2222-2222-2222-222222222222'
    ),
    (
        'f6666666-6666-6666-6666-666666666666',
        'a1111111-1111-1111-1111-111111111111'
    ),
    (
        'f6666666-6666-6666-6666-666666666666',
        'c3333333-3333-3333-3333-333333333333'
    ) ON CONFLICT (product_id, category_id) DO NOTHING;

-- ============================================
-- 6. ORDERS (3 records) - customer_id is bigint
-- ============================================
INSERT INTO
    orders (
        id,
        customer_id,
        status,
        payment_status,
        total_price
    )
VALUES (
        '10101010-0101-0101-0101-010101010101',
        1,
        'completed',
        'paid',
        3598.99
    ),
    (
        '20202020-0202-0202-0202-020202020202',
        2,
        'pending',
        'unpaid',
        449.99
    ),
    (
        '30303030-0303-0303-0303-030303030303',
        1,
        'completed',
        'paid',
        99.99
    ) ON CONFLICT (id) DO NOTHING;

-- ============================================
-- 7. ORDER_ITEMS (4 records) - id is auto-increment
-- ============================================
INSERT INTO
    order_items (
        order_id,
        product_id,
        variant_id,
        quantity,
        price
    )
VALUES (
        '10101010-0101-0101-0101-010101010101',
        'd4444444-4444-4444-4444-444444444444',
        '77777777-7777-7777-7777-777777777777',
        1,
        3499.00
    ),
    (
        '10101010-0101-0101-0101-010101010101',
        'e5555555-5555-5555-5555-555555555555',
        '99999999-9999-9999-9999-999999999999',
        1,
        99.99
    ),
    (
        '20202020-0202-0202-0202-020202020202',
        'f6666666-6666-6666-6666-666666666666',
        NULL,
        1,
        449.99
    ),
    (
        '30303030-0303-0303-0303-030303030303',
        'e5555555-5555-5555-5555-555555555555',
        '99999999-9999-9999-9999-999999999999',
        1,
        99.99
    );

-- ============================================
-- 8. WEBHOOK_LOGS (3 records)
-- ============================================
INSERT INTO
    webhook_logs (
        id,
        order_id,
        transaction_id,
        payment_status,
        status,
        retry_count,
        next_retry_at,
        raw_payload,
        processed_at
    )
VALUES (
        '80808080-0808-0808-0808-080808080808',
        '10101010-0101-0101-0101-010101010101',
        'txn_abc123xyz789',
        'paid',
        'completed',
        0,
        NULL,
        '{"order_id":"10101010-0101-0101-0101-010101010101","transaction_id":"txn_abc123xyz789","payment_status":"paid","amount":3598.99}',
        NOW()
    ),
    (
        '90909090-0909-0909-0909-090909090909',
        '20202020-0202-0202-0202-020202020202',
        'txn_def456uvw012',
        'unpaid',
        'pending',
        0,
        NOW() + INTERVAL '1 hour',
        '{"order_id":"20202020-0202-0202-0202-020202020202","transaction_id":"txn_def456uvw012","payment_status":"unpaid","amount":449.99}',
        NULL
    ),
    (
        'a0a0a0a0-0a0a-0a0a-0a0a-0a0a0a0a0a0a',
        '30303030-0303-0303-0303-030303030303',
        'txn_ghi789rst345',
        'paid',
        'completed',
        0,
        NULL,
        '{"order_id":"30303030-0303-0303-0303-030303030303","transaction_id":"txn_ghi789rst345","payment_status":"paid","amount":99.99}',
        NOW()
    ) ON CONFLICT (transaction_id) DO NOTHING;

COMMIT;

-- ============================================
-- VERIFICATION
-- ============================================
DO $$
DECLARE
  user_count INT;
  category_count INT;
  product_count INT;
  variant_count INT;
  pc_count INT;
  order_count INT;
  item_count INT;
  webhook_count INT;
BEGIN
  SELECT COUNT(*) INTO user_count FROM users WHERE email LIKE '%@ecommerce.com' OR email LIKE '%@example.com';
  SELECT COUNT(*) INTO category_count FROM categories WHERE name IN ('Electronics', 'Computers', 'Gaming');
  SELECT COUNT(*) INTO product_count FROM products WHERE name LIKE '%MacBook%' OR name LIKE '%Logitech%' OR name LIKE '%PlayStation%';
  SELECT COUNT(*) INTO variant_count FROM product_variants WHERE variant_name = 'Color';
  SELECT COUNT(*) INTO pc_count FROM product_categories;
  SELECT COUNT(*) INTO order_count FROM orders;
  SELECT COUNT(*) INTO item_count FROM order_items;
  SELECT COUNT(*) INTO webhook_count FROM webhook_logs;
  
  RAISE NOTICE '========================================';
  RAISE NOTICE 'Database Seeding Complete!';
  RAISE NOTICE '========================================';
  RAISE NOTICE 'Sample data inserted:';
  RAISE NOTICE '  Users:              %', user_count;
  RAISE NOTICE '  Categories:         %', category_count;
  RAISE NOTICE '  Products:           %', product_count;
  RAISE NOTICE '  Product Variants:   %', variant_count;
  RAISE NOTICE '  Product-Categories: %', pc_count;
  RAISE NOTICE '  Orders:             %', order_count;
  RAISE NOTICE '  Order Items:        %', item_count;
  RAISE NOTICE '  Webhook Logs:       %', webhook_count;
  RAISE NOTICE '========================================';
  RAISE NOTICE 'Sample Credentials (password: password123):';
  RAISE NOTICE '  Admin:     admin@ecommerce.com';
  RAISE NOTICE '  Customer1: john.doe@example.com';
  RAISE NOTICE '  Customer2: jane.smith@example.com';
  RAISE NOTICE '========================================';
END $$;