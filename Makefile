# BSV Custodial SDK Makefile

.PHONY: help build test clean run-example run-basic run-sharding run-transaction install deps

# Default target
help:
	@echo "BSV Custodial SDK - Available Commands:"
	@echo "======================================="
	@echo "  make help           - Show this help message"
	@echo "  make deps           - Download dependencies"
	@echo "  make build          - Build the library"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make run-example    - Run interactive example"
	@echo "  make run-basic      - Run basic usage example"
	@echo "  make run-sharding   - Run sharding example"
	@echo "  make run-transaction - Run transaction example"
	@echo "  make install        - Install the library"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	go mod download
	go mod tidy

# Build the library
build: deps
	@echo "🔨 Building BSV Custodial SDK..."
	go build ./...

# Run tests
test: deps
	@echo "🧪 Running tests..."
	go test -v ./tests/...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	go clean
	rm -f bsv-custodial-sdk

# Run interactive example
run-example: build
	@echo "🚀 Running interactive example..."
	go run cmd/main.go

# Run basic usage example
run-basic: build
	@echo "🚀 Running basic usage example..."
	go run examples/enhanced/enhanced_usage.go

# Run sharding example
run-sharding: build
	@echo "🚀 Running sharding example..."
	go run examples/sharding/sharding_example.go

# Run transaction example
run-transaction: build
	@echo "🚀 Running transaction example..."
	go run examples/transaction/transaction_example.go

# Install the library
install: build
	@echo "📦 Installing BSV Custodial SDK..."
	go install ./...

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Install with:"; \
		echo "   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.54.2"; \
	fi

# Run all examples
examples: run-basic run-sharding run-transaction

# Development setup
dev-setup: deps
	@echo "🛠️  Setting up development environment..."
	@echo "✅ Dependencies downloaded"
	@echo "✅ Ready for development!"
	@echo ""
	@echo "Next steps:"
	@echo "  make test          - Run tests"
	@echo "  make run-example   - Try interactive example"
	@echo "  make examples      - Run all examples"

# CI/CD pipeline
ci: deps fmt test build
	@echo "✅ CI pipeline completed successfully!"

# Release build
release: clean deps test build
	@echo "🚀 Building release..."
	@echo "✅ Release build completed!"
	@echo ""
	@echo "Release artifacts:"
	@echo "  - Library: ./pkg/"
	@echo "  - Examples: ./examples/"
	@echo "  - Tests: ./tests/"
