.PHONY: build run clean install uninstall fmt lint test coverage help

# Variables
BINARY_NAME=MemoryAnalyzer
GO_FILES=$(shell find . -type f -name '*.go')
VERSION ?= 1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LD_FLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

## help: Display this help message
help:
	@echo "$(GREEN)Memory Analyzer - Makefile Commands$(NC)"
	@echo ""
	@echo "$(YELLOW)Build Commands:$(NC)"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the application
build: fmt
	@echo "$(GREEN)[*] Building $(BINARY_NAME)...$(NC)"
	@go build -v $(LD_FLAGS) -o $(BINARY_NAME) .
	@echo "$(GREEN)[✓] Build complete: ./$(BINARY_NAME)$(NC)"

## build-linux: Build for Linux (AMD64)
build-linux:
	@echo "$(GREEN)[*] Building for Linux AMD64...$(NC)"
	@GOOS=linux GOARCH=amd64 go build -v $(LD_FLAGS) -o $(BINARY_NAME)-linux .
	@echo "$(GREEN)[✓] Linux build complete: ./$(BINARY_NAME)-linux$(NC)"

## build-darwin-intel: Build for macOS Intel (AMD64)
build-darwin-intel:
	@echo "$(GREEN)[*] Building for macOS Intel...$(NC)"
	@GOOS=darwin GOARCH=amd64 go build -v $(LD_FLAGS) -o $(BINARY_NAME)-darwin-amd64 .
	@echo "$(GREEN)[✓] macOS Intel build complete: ./$(BINARY_NAME)-darwin-amd64$(NC)"

## build-darwin-arm: Build for macOS ARM (M1/M2/M3)
build-darwin-arm:
	@echo "$(GREEN)[*] Building for macOS ARM...$(NC)"
	@GOOS=darwin GOARCH=arm64 go build -v $(LD_FLAGS) -o $(BINARY_NAME)-darwin-arm64 .
	@echo "$(GREEN)[✓] macOS ARM build complete: ./$(BINARY_NAME)-darwin-arm64$(NC)"

## build-all: Build for all supported platforms
build-all: build-linux build-darwin-intel build-darwin-arm build
	@echo "$(GREEN)[✓] All platform builds complete$(NC)"

## run: Build and run the application
run: build
	@echo "$(GREEN)[*] Running $(BINARY_NAME)...$(NC)"
	@./$(BINARY_NAME)

## clean: Remove build artifacts
clean:
	@echo "$(YELLOW)[*] Cleaning up...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME)-*
	@rm -f *.out
	@rm -f coverage.html
	@go clean -cache -modcache
	@echo "$(GREEN)[✓] Cleanup complete$(NC)"

## install: Install the binary to system PATH
install: build
	@echo "$(GREEN)[*] Installing $(BINARY_NAME)...$(NC)"
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)[✓] Installed to /usr/local/bin/$(BINARY_NAME)$(NC)"
	@echo "$(GREEN)[✓] Run '$(BINARY_NAME)' from anywhere$(NC)"

## uninstall: Remove the binary from system PATH
uninstall:
	@echo "$(YELLOW)[*] Uninstalling $(BINARY_NAME)...$(NC)"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)[✓] Uninstalled$(NC)"

## fmt: Format code
fmt:
	@echo "$(YELLOW)[*] Formatting code...$(NC)"
	@go fmt ./...
	@gofmt -w .
	@echo "$(GREEN)[✓] Code formatted$(NC)"

## deps: Download and verify dependencies
deps:
	@echo "$(YELLOW)[*] Downloading dependencies...$(NC)"
	@go mod download
	@go mod verify
	@echo "$(GREEN)[✓] Dependencies verified$(NC)"

## tidy: Tidy Go modules
tidy:
	@echo "$(YELLOW)[*] Running go mod tidy...$(NC)"
	@go mod tidy
	@echo "$(GREEN)[✓] Go modules tidied$(NC)"

## version: Show version information
version:
	@echo "$(GREEN)Memory Analyzer - Version Information$(NC)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"


# Default target
.DEFAULT_GOAL := help
