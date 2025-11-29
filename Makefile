BINARY_NAME=cutter
API_BINARY=devops-cutter-api
VERSION=0.1.0
BUILD_DIR=build

.PHONY: all
all: build

.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: build
build: build-api build-cli

.PHONY: build-api
build-api:
	@echo "Building API server..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(API_BINARY) cmd/api/main.go

.PHONY: build-cli
build-cli:
	@echo "Building CLI..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) cmd/cutter/main.go

.PHONY: install-cli
install-cli: build-cli
	@echo "Installing cutter to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

.PHONY: run-api
run-api:
	go run cmd/api/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover ./...

.PHONY: test-coverage-report
test-coverage-report:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-coverage-func
test-coverage-func:
	@echo "Generating function coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make deps                  - Download dependencies"
	@echo "  make build                 - Build all binaries"
	@echo "  make build-cli             - Build CLI only"
	@echo "  make build-api             - Build API server only"
	@echo "  make install-cli           - Install CLI to /usr/local/bin"
	@echo "  make run-api               - Run API server locally"
	@echo "  make test                  - Run tests"
	@echo "  make test-coverage         - Run tests with coverage summary"
	@echo "  make test-coverage-report  - Generate HTML coverage report"
	@echo "  make test-coverage-func    - Show per-function coverage"
	@echo "  make clean                 - Clean build artifacts and coverage files"