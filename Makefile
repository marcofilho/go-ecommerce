.PHONY: build run test clean help db-up db-down db-reset

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/api ./src/cmd/api
	@echo "Build complete! Binary located at bin/api"

# Start PostgreSQL database
db-up:
	@echo "Starting PostgreSQL database..."
	@docker-compose up -d
	@echo "Waiting for database to be ready..."
	@sleep 3
	@echo "Database is ready!"

# Stop PostgreSQL database
db-down:
	@echo "Stopping PostgreSQL database..."
	@docker-compose down
	@echo "Database stopped!"

# Reset database (stop, remove volumes, and start fresh)
db-reset:
	@echo "Resetting database..."
	@docker-compose down -v
	@docker-compose up -d
	@echo "Waiting for database to be ready..."
	@sleep 3
	@echo "Database reset complete!"

# Run migrations
migrate: db-up
	@echo "Running database migrations..."
	@go run ./src/cmd/migrate/main.go

# Run the application (starts database if not running)
run: db-up build
	@echo "Starting server..."
	@./bin/api

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Test API endpoints (requires server to be running)
test-api:
	@echo "Testing API endpoints..."
	@./test_api.sh

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean complete!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@echo "Dependencies installed!"

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev: db-up
	@air

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete!"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@golangci-lint run
	@echo "Lint complete!"

# Show help
help:
	@echo "Available commands:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Start database and run the application"
	@echo "  make db-up     - Start PostgreSQL database"
	@echo "  make db-down   - Stop PostgreSQL database"
	@echo "  make db-reset  - Reset database (fresh start)"
	@echo "  make migrate   - Run database migrations"
	@echo "  make test      - Run unit tests"
	@echo "  make test-api  - Test API endpoints (server must be running)"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make deps      - Install dependencies"
	@echo "  make dev       - Run with hot reload (requires air)"
	@echo "  make fmt       - Format code"
	@echo "  make lint      - Run linter"
	@echo "  make help      - Show this help message"

