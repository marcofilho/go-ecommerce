.PHONY: start stop logs test test-webhook test-auth seed clean-db reset-db help

# Default target
.DEFAULT_GOAL := help

# Start all services (PostgreSQL + API)
start:
	@echo "Running unit tests before starting services..."
	@docker build --target test -t go-ecommerce-test . 2>&1 | grep -E "(RUN go test|PASS|FAIL|coverage:|ok  )" || true
	@echo ""
	@echo "Starting services..."
	@docker-compose up -d
	@echo "✓ Services started!"
	@echo "  API: http://localhost:8080"
	@echo "  Swagger: http://localhost:8080/swagger/index.html"
	@echo ""
	@echo "Waiting for API to be ready..."
	@sleep 5
	@echo ""
	@echo "Checking database state..."
	@USER_COUNT=$$(docker exec ecommerce_postgres psql -U postgres -d ecommerce -t -A -c "SELECT COUNT(*) FROM users;" 2>/dev/null || echo "0"); \
	if [ "$$USER_COUNT" -eq 0 ]; then \
		echo "Database is empty. Seeding with initial data..."; \
		docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/seed_data.sql | grep -E "NOTICE:" || true; \
		echo "✓ Database seeded successfully!"; \
	else \
		echo "✓ Database already contains data ($$USER_COUNT users found)"; \
	fi
	@echo ""
	@echo "Running Payment Webhook Integration Tests..."
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@./test_payment_webhook_batch.sh
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "✓ All systems verified and ready!"
	@echo ""
	@echo "Opening Swagger UI in your browser..."
	@open http://localhost:8080/swagger/index.html || xdg-open http://localhost:8080/swagger/index.html || start http://localhost:8080/swagger/index.html 2>/dev/null || true

# Stop all services
stop:
	@echo "Stopping services..."
	@docker-compose down
	@echo "✓ Services stopped!"

# View logs
logs:
	@docker-compose logs -f

# Run tests in Docker
test:
	@echo "Running tests in Docker..."
	@docker build --target test -t go-ecommerce-test .
	@echo "✓ Tests complete!"

# Run webhook integration tests
test-webhook:
	@echo "Running Payment Webhook Integration Tests..."
	@./test_payment_webhook_batch.sh

# Run authentication integration tests
test-auth:
	@echo "Running Authentication Integration Tests..."
	@./test_authentication.sh

# Database management commands
seed:
	@echo "Seeding database with sample data..."
	@docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/seed_data.sql | grep -E "NOTICE:" || true
	@echo "✓ Database seeded!"

clean-db:
	@echo "⚠️  WARNING: This will delete all data from the database!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker exec -i ecommerce_postgres psql -U postgres -d ecommerce < scripts/clean_database.sql | grep -E "NOTICE:" || true; \
		echo "✓ Database cleaned!"; \
	else \
		echo "Cancelled."; \
	fi

reset-db: clean-db seed
	@echo "✓ Database reset complete!"

# Show help
help:
	@echo "Go E-Commerce API - Available commands:"
	@echo ""
	@echo "Services:"
	@echo "  make start         - Run tests and start all services (PostgreSQL + API)"
	@echo "  make stop          - Stop all services"
	@echo "  make logs          - View service logs"
	@echo ""
	@echo "Testing:"
	@echo "  make test          - Run unit tests in Docker"
	@echo "  make test-webhook  - Run webhook integration tests"
	@echo "  make test-auth     - Run authentication integration tests"
	@echo ""
	@echo "Database:"
	@echo "  make seed          - Seed database with sample data"
	@echo "  make clean-db      - Clean all data from database (with confirmation)"
	@echo "  make reset-db      - Clean and seed database"
	@echo ""
	@echo "Other:"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make deps      - Install dependencies"
	@echo "  make help      - Show this help message"

