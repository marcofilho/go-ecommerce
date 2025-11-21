.PHONY: up down logs test help

# Default target
.DEFAULT_GOAL := help

# Start all services (PostgreSQL + API)
up:
	@echo "Starting services..."
	@docker-compose up -d
	@echo "✓ Services started!"
	@echo "  API: http://localhost:8080"
	@echo "  Swagger: http://localhost:8080/swagger/index.html"

# Stop all services
down:
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

# Show help
help:
	@echo "Go E-Commerce API - Available commands:"
	@echo ""
	@echo "  make up     - Start all services (PostgreSQL + API)"
	@echo "  make down   - Stop all services"
	@echo "  make logs   - View service logs"
	@echo "  make test   - Run tests in Docker"
	@echo ""
	@echo "Other:"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make deps      - Install dependencies"
	@echo "  make help      - Show this help message"

