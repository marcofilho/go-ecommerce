.PHONY: build run test clean help db-up db-down db-reset docker-build docker-up docker-down docker-logs docker-restart swagger migrate deps dev fmt lint

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/api ./src/cmd/api
	@echo "Build complete! Binary located at bin/api"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker-compose build
	@echo "Docker image built!"

# Start all services (PostgreSQL + API)
docker-up:
	@echo "Starting all services..."
	@docker-compose up -d
	@echo "Services started! API: http://localhost:8080"

# Stop all services
docker-down:
	@echo "Stopping all services..."
	@docker-compose down
	@echo "Services stopped!"

# View logs from all services
docker-logs:
	@docker-compose logs -f

# Restart all services
docker-restart: docker-down docker-up

# Start PostgreSQL only
db-up:
	@echo "Starting PostgreSQL database..."
	@docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 3
	@echo "Database is ready!"

# Stop PostgreSQL database
db-down:
	@echo "Stopping PostgreSQL database..."
	@docker-compose stop postgres
	@echo "Database stopped!"

# Reset database (stop, remove volumes, and start fresh)
db-reset:
	@echo "Resetting database..."
	@docker-compose down -v
	@docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 3
	@echo "Database reset complete!"

# Run migrations
migrate: db-up
	@echo "Running database migrations..."
	@go run ./src/cmd/migrate/main.go

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@$(HOME)/go/bin/swag init -g src/cmd/api/main.go -o docs
	@echo "Swagger docs generated! Visit http://localhost:8080/swagger/ after starting the server"

# Run the application (starts database if not running)
run: db-up build
	@echo "Starting server..."
	@./bin/api

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

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
	@echo ""
	@echo "Docker (Full Stack):"
	@echo "  make docker-build   - Build Docker image for API"
	@echo "  make docker-up      - Start all services (PostgreSQL + API)"
	@echo "  make docker-down    - Stop all services"
	@echo "  make docker-logs    - View logs from all services"
	@echo "  make docker-restart - Restart all services"
	@echo ""
	@echo "Local Development:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Start database and run application locally"
	@echo "  make swagger   - Generate Swagger documentation"
	@echo ""
	@echo "Database:"
	@echo "  make db-up     - Start PostgreSQL only"
	@echo "  make db-down   - Stop PostgreSQL only"
	@echo "  make db-reset  - Reset database (fresh start)"
	@echo ""
	@echo "Testing & Quality:"
	@echo "  make test      - Run unit tests"
	@echo "  make fmt       - Format code"
	@echo "  make lint      - Run linter"
	@echo ""
	@echo "Other:"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make deps      - Install dependencies"
	@echo "  make help      - Show this help message"

