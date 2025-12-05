.PHONY: start stop logs test test-webhook test-auth help

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

# Show help
help:
	@echo "Go E-Commerce API - Available commands:"
	@echo ""
	@echo "  make start         - Run tests and start all services (PostgreSQL + API)"
	@echo "  make stop          - Stop all services"
	@echo "  make logs          - View service logs"
	@echo "  make test          - Run unit tests in Docker"
	@echo "  make test-webhook  - Run webhook integration tests"
	@echo "  make test-auth     - Run authentication integration tests"
	@echo ""
	@echo "Other:"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make deps      - Install dependencies"
	@echo "  make help      - Show this help message"

