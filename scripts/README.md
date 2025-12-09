# Database Scripts

This folder contains SQL scripts for database management and initialization.

## Available Scripts

### 1. `seed_data.sql` - Database Seeding Script

Populates the database with initial sample data for testing and development.

**What it creates:**
- 3 users (1 admin, 2 customers)
- 3 product categories (Electronics, Computers & Accessories, Gaming)
- 3 products (MacBook Pro, Logitech Mouse, PlayStation 5)
- 3 product variants (different colors)
- 6 product-category relationships
- 3 orders (2 completed, 1 pending)
- 4 order items
- 3 webhook logs

**Usage:**

```bash
# Using Docker (recommended)
docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/seed_data.sql

# Direct PostgreSQL connection
psql -U postgres -d ecommerce -f scripts/seed_data.sql

# Using make (if added to Makefile)
make seed
```

**Sample Credentials:**
- Admin: `admin@ecommerce.com` / `password123`
- Customer 1: `john.doe@example.com` / `password123`
- Customer 2: `jane.smith@example.com` / `password123`

**Sample Data Highlights:**
- **MacBook Pro 16" M3 Max** - $3,499.00 (25 in stock)
  - Variants: Space Black, Silver
  - Categories: Electronics, Computers & Accessories
  
- **Logitech MX Master 3S Mouse** - $99.99 (150 in stock)
  - Variant: Graphite
  - Categories: Electronics, Computers & Accessories
  
- **Sony PlayStation 5 Digital** - $449.99 (40 in stock)
  - Categories: Electronics, Gaming

---

### 2. `clean_database.sql` - Database Cleanup Script

⚠️ **WARNING: This script deletes ALL data from the database!**

Truncates all tables and resets sequences. Use with caution!

**What it does:**
- Removes all records from all tables
- Resets auto-increment sequences (users.id)
- Maintains table structure and constraints
- Provides verification summary

**Usage:**

```bash
# Using Docker (recommended)
docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/clean_database.sql

# Direct PostgreSQL connection
psql -U postgres -d ecommerce -f scripts/clean_database.sql

# Using make (if added to Makefile)
make clean-db
```

**Safety Notes:**
- Always backup your database before running cleanup
- Never run this script on production databases
- The script is wrapped in a transaction (BEGIN/COMMIT)
- Verification query confirms all tables are empty

---

## Common Workflows

### Fresh Start (Clean + Seed)

```bash
# Clean the database
docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/clean_database.sql

# Seed with sample data
docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/seed_data.sql
```

### Reset Development Database

```bash
# Stop the API
make stop

# Clean and seed
docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/clean_database.sql
docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/seed_data.sql

# Start the API
make start
```

### Verify Database State

```bash
# Check record counts
docker exec ecommerce_postgres psql -U postgres -d ecommerce -c "
  SELECT 
    (SELECT COUNT(*) FROM users) as users,
    (SELECT COUNT(*) FROM categories) as categories,
    (SELECT COUNT(*) FROM products) as products,
    (SELECT COUNT(*) FROM product_variants) as variants,
    (SELECT COUNT(*) FROM orders) as orders,
    (SELECT COUNT(*) FROM order_items) as items;
"
```

---

## Adding to Makefile

Consider adding these commands to your `Makefile` for easier access:

```makefile
.PHONY: seed clean-db reset-db

seed: ## Seed database with sample data
	@echo "Seeding database..."
	@docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/seed_data.sql

clean-db: ## Clean all data from database
	@echo "⚠️  WARNING: This will delete all data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/clean_database.sql; \
	fi

reset-db: clean-db seed ## Clean database and seed with fresh data
	@echo "✓ Database reset complete!"
```

Then you can use:
```bash
make seed      # Seed database
make clean-db  # Clean database (with confirmation)
make reset-db  # Clean and seed
```

---

## Troubleshooting

### Connection Refused
```bash
# Make sure PostgreSQL container is running
docker ps | grep postgres

# Start containers if needed
make start
```

### Permission Denied
```bash
# Ensure you're using the correct PostgreSQL user
docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/seed_data.sql
```

### Duplicate Key Errors
The seed script includes `ON CONFLICT DO NOTHING` clauses, so it's safe to run multiple times. However, if you want truly fresh data, run the clean script first:

```bash
docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/clean_database.sql
docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/seed_data.sql
```

---

## Notes

- Scripts use fixed UUIDs for reproducibility and easier testing
- All passwords are hashed with bcrypt (password: `password123`)
- Transaction IDs use realistic format (e.g., `txn_stripe_abc123xyz789`)
- Sample data includes realistic product descriptions and prices
- Foreign key relationships are properly maintained
- Scripts include verification queries and success messages
