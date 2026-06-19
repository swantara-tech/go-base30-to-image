# jSignature Base30 Converter - Makefile

# Variables
BINARY_NAME=jsign-convert
VERSION=1.0.0
BUILD_DIR=build

# Go commands
GO=go
GOFMT=gofmt
GOLINT=golangci-lint

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Default target
.PHONY: all
all: clean test build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for all platforms..."
	@$(GO) env -w GOOS=linux GOARCH=amd64
	@$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@$(GO) env -w GOOS=windows GOARCH=amd64
	@$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@$(GO) env -w GOOS=darwin GOARCH=amd64
	@$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@$(GO) env -u GOOS GOARCH
	@echo "Build complete for all platforms"

# Run tests with coverage
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -v -coverprofile=coverage.out ./...
	@echo "Tests complete"

# Run tests and open coverage report
.PHONY: test-coverage
test-coverage: test
	@echo "Generating coverage report..."
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Check test coverage percentage
.PHONY: test-coverage-check
test-coverage-check: test
	@echo "Checking coverage..."
	@$(GO) tool cover -func=coverage.out | findstr /C:"total:"

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	$(GOLINT) run
	@echo "Lint complete"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@echo "Format complete"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@$(GO) clean
	-@if exist "$(BUILD_DIR)" rmdir /s /q "$(BUILD_DIR)"
	-@del /q coverage.out coverage.html $(BINARY_NAME).exe $(BINARY_NAME)
	@echo "Clean complete"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GO) mod tidy
	@echo "Dependencies installed"

# Run the application (development)
.PHONY: run
run:
	@echo "Running application..."
	$(GO) run . $(ARGS)

# Create example files
.PHONY: examples
examples:
	@echo "Creating example files..."
	@if not exist "examples" mkdir "examples"
	@echo "image/jsignature;base30,5K247669cffhlo1vmhc9852" > examples/signature.txt
	@echo "id,signature" > examples/signatures.csv
	@echo "1,image/jsignature;base30,5K247669cffhlo1vmhc9852" >> examples/signatures.csv
	@echo "2,image/jsignature;base30,7L358770dglimp2wnid0963" >> examples/signatures.csv
	@echo "Examples created in examples/ directory"

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build              - Build the binary"
	@echo "  build-all          - Build for multiple platforms"
	@echo "  test               - Run tests with coverage"
	@echo "  test-coverage      - Run tests and generate HTML coverage report"
	@echo "  test-coverage-check - Check test coverage percentage"
	@echo "  lint               - Run linter"
	@echo "  fmt                - Format code"
	@echo "  clean              - Clean build artifacts"
	@echo "  deps               - Install dependencies"
	@echo "  run                - Run the application (use ARGS for arguments)"
	@echo "  examples           - Create example files"
	@echo "  help               - Show this help"
