.PHONY: help build install uninstall clean run fmt deps tidy

BINARY_NAME=zap
GOPATH=$(shell go env GOPATH)
INSTALL_PATH=$(GOPATH)/bin/$(BINARY_NAME)

.DEFAULT_GOAL := help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .
	@echo "✓ Build complete: ./$(BINARY_NAME)"

install: ## Install the CLI tool to GOPATH/bin
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	go install .
	@echo "✓ $(BINARY_NAME) installed successfully!"
	@echo "  Run '$(BINARY_NAME) --help' to get started"

uninstall: ## Uninstall the CLI tool from GOPATH/bin
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_PATH)
	@echo "✓ $(BINARY_NAME) uninstalled"

clean: ## Clean build artifacts and test cache
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@echo "✓ Clean complete"

run: ## Run the application (use ARGS="your args" to pass arguments)
	@go run . $(ARGS)

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@echo "✓ Dependencies downloaded"

tidy: ## Tidy go.mod and go.sum
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "✓ Dependencies tidied"
