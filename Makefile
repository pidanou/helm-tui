# Makefile for Helmtui

# Go build settings
BINARY_NAME = helmtui
SOURCE_DIR = .
BUILD_DIR = ./bin

# Go commands
GO = go
GOFMT = gofmt
GOTEST = go test
GOBUILD = go build

# Default target
all: build

# Clean build directory
clean:
	rm -rf $(BUILD_DIR)

# Format Go source files
fmt:
	$(GOFMT) -w $(SOURCE_DIR)

# Build the project
build: clean fmt
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(SOURCE_DIR)

# Run the project
run: build
	$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	$(GOTEST) ./...

# Install dependencies
deps:
	$(GO) mod tidy

# Help message
help:
	@echo "Available targets:"
	@echo "  clean    - Clean the build directory"
	@echo "  fmt      - Format the Go source code"
	@echo "  build    - Build the project"
	@echo "  run      - Run the built binary"
	@echo "  test     - Run tests"
	@echo "  deps     - Install and tidy dependencies"
	@echo "  help     - Show this help message"

.PHONY: all clean fmt build run test deps help
