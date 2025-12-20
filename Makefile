.PHONY: build test run docker-up docker-down migrate

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run tests
test:
	go test -v ./...

# Run the application
run:
	go run ./cmd/server/main.go

# Start Docker services
docker-up:
	docker-compose up -d postgres redis

# Stop Docker services
docker-down:
	docker-compose down

# Run database migrations
migrate:
	@echo "Please run migrations manually using psql or your preferred migration tool"
	@echo "Example: psql -h localhost -U postgres -d erh_safety -f database/migrations/001_initial_schema.up.sql"

# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

