# Woody's Backend Makefile
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
BINARY_NAME=woodys-backend
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./main.go

# Build flags
BUILD_FLAGS=-v
LDFLAGS=-ldflags "-s -w"

# Linting
GOLANGCI_LINT=golangci-lint
GOLANGCI_VERSION=v1.55.2

.PHONY: all build run test clean deps lint fmt vet help install-tools check 

# Default target
all: clean deps lint build

# Help target
help: ## Show this help message
	@echo "Woody's Backend - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'
	@echo ""

# Build the application
build: ## Build the application binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "✓ Build completed: $(BINARY_PATH)"

# Run the application
run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run $(MAIN_PATH)

# Run the built binary
run-bin: build ## Run the built binary
	@echo "Running $(BINARY_PATH)..."
	$(BINARY_PATH)

# TODO: Write test xd
# # Run tests
# test: ## Run tests
# 	@echo "Running tests..."
# 	$(GOTEST) -v -race -coverprofile=coverage.out ./...
# 	@echo "✓ Tests completed"
#
# # Run tests with coverage report
# test-coverage: test ## Run tests and show coverage report
# 	@echo "Generating coverage report..."
# 	$(GOCMD) tool cover -html=coverage.out -o coverage.html
# 	@echo "✓ Coverage report generated: coverage.html"
#

# Clean build artifacts
clean: ## Clean build artifacts and cache
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✓ Cleaned"

# Install/update dependencies
deps: ## Download and install dependencies
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✓ Dependencies updated"

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@if command -v goimports >/dev/null 2>&1; then \
		echo "Running goimports..."; \
		goimports -w .; \
	fi
	@echo "✓ Code formatted"

# Run go vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...
	@echo "✓ go vet completed"

# Comprehensive linting
lint: ## Run all linters
	@echo "Running linters..."
	@if command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		$(GOLANGCI_LINT) run --timeout=5m ./...; \
	else \
		echo "⚠ golangci-lint not installed. Run 'make install-tools' to install it"; \
		echo "Running basic linters instead..."; \
		$(MAKE) vet; \
	fi
	@echo "✓ Linting completed"

# Install development tools
install-tools: ## Install development tools (golangci-lint, gosec, goimports)
	@echo "Installing development tools..."
	@if ! command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $(GOLANGCI_VERSION); \
	fi
	@if ! command -v gosec >/dev/null 2>&1; then \
		echo "Installing gosec..."; \
		$(GOCMD) install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	fi
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "Installing goimports..."; \
		$(GOCMD) install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@echo "✓ Development tools installed"

# Pre-commit checks
check: fmt lint vet ## Run all pre-commit checks
	@echo "✓ All checks passed"
