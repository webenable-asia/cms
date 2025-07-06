# WebEnable CMS Backend Makefile

.PHONY: test test-verbose test-coverage build run clean lint fmt vet

# Test commands
test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-race:
	go test -race ./...

# Build commands
build:
	go build -o bin/webenable-cms-backend .

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/webenable-cms-backend-linux .

# Run commands
run:
	go run .

run-dev:
	air

# Code quality
lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

# Dependencies
deps:
	go mod download
	go mod tidy

# Clean
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker
docker-build:
	docker build -t webenable-cms-backend .

docker-run:
	docker-compose up --build

# Database
init-admin:
	go run scripts/init_admin.go

populate-db:
	go run scripts/populate_database.go

# Help
help:
	@echo "Available commands:"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-race     - Run tests with race detection"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  run-dev       - Run with air for development"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  init-admin    - Initialize admin user"
	@echo "  populate-db   - Populate database with sample data"