.PHONY: swagger run build test clean dev seed

# Generate swagger documentation
swagger:
	@echo "Generating swagger documentation..."
	@swag init -g ./cmd/server/main.go -o ./cmd/server/docs --parseDependency --parseInternal

# Run the application
run:
	@go run ./cmd/server/main.go

# Run with air (hot reload)
dev:
	@air

# Build the application
build:
	@echo "Building application..."
	@go build -o ./bin/paygo ./cmd/server/main.go

# Run tests
test:
	@go test -v ./...

# Seed the database with test data
seed:
	@go run ./cmd/seed/main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf ./tmp ./bin
	@go clean

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Build and run
start: build
	@./bin/paygo
