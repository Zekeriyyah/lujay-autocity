# Makefile for AutoCity Go Backend

# --- Variables ---
BINARY_NAME=autocity
MAIN_PATH=./cmd/autocity_server
BUILD_DIR=./build
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOLINT=golangci-lint
GOMOD=$(GOCMD) mod
GOTIDY=$(GOMOD) tidy

# --- Default target ---
.PHONY: help
help:  
	@echo "AutoCity Backend Makefile"
	@echo "Usage:"
	@echo "  make build     - Build the application binary"
	@echo "  make run       - Run the application (assumes dependencies are met)"
	@echo "  make test      - Run tests with verbose output"
	@echo "  make lint      - Run golangci-lint (requires golangci-lint installed)"
	@echo "  make tidy      - Run go mod tidy to manage dependencies"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make help      - Show this help message"
	@echo ""

# --- Build ---
.PHONY: build
build:  
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete. Binary located at $(BUILD_DIR)/$(BINARY_NAME)"

# --- Run ---
.PHONY: run
run: build  
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)


# --- Lint ---
.PHONY: lint
lint:  
	@echo "Running linter..."
	$(GOLINT) run ./...

# --- Tidy (Manage Dependencies) ---
.PHONY: tidy
tidy: 
	@echo "Tidying Go modules..."
	$(GOTIDY)

# --- Clean ---
.PHONY: clean
clean:  
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN)
