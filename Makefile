.PHONY: build run test clean install-deps help

# Default target
help:
	@echo "Available targets:"
	@echo "  build         - Build the scraper binary"
	@echo "  run           - Run the scraper"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts and output"
	@echo "  install-deps  - Install Go dependencies"

# Build the application
build:
	@echo "Building scraper..."
	go build -o bin/scraper ./cmd/scraper

# Run the application
run: build
	@echo "Running scraper..."
	./bin/scraper

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts and output
clean:
	@echo "Cleaning..."
	rm -rf bin output
	go clean

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	go mod download
