# BubbleTable Development Makefile

.PHONY: help test test-verbose test-coverage benchmark lint fmt vet clean build examples doc

# Default target
help: ## Show this help message
	@echo "BubbleTable Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

# Testing
test: ## Run all tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-coverage: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

benchmark: ## Run benchmarks
	go test -bench=. -benchmem ./...

# Code quality
lint: ## Run linter (requires golangci-lint)
	golangci-lint run

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

# Build
build: ## Build the library and examples
	go build ./...
	@echo "Building examples..."
	@cd examples/basic && go build
	@cd examples/custom_theme && go build
	@cd examples/headless && go build

# Examples
examples: ## Run all examples
	@echo "Running basic example..."
	@cd examples/basic && go run main.go
	@echo "Running custom theme example..."
	@cd examples/custom_theme && go run main.go
	@echo "Running headless example..."
	@cd examples/headless && go run main.go

# Documentation
doc: ## Generate and serve documentation
	godoc -http=:6060
	@echo "Documentation available at http://localhost:6060/pkg/github.com/anurag-roy/bubbletable/"

# Cleanup
clean: ## Clean build artifacts and coverage files
	rm -f coverage.out coverage.html
	rm -f examples/basic/basic
	rm -f examples/custom_theme/custom_theme
	rm -f examples/headless/headless

# Development setup
setup: ## Install development dependencies
	go mod tidy
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Release preparation
pre-release: fmt vet lint test-coverage ## Run all checks before release
	@echo "All checks passed! Ready for release." 