.PHONY: build build-darwin-amd64 build-darwin-arm64 build-all install release clean test help

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BINARY_NAME = day-night-cycle
INSTALL_PATH = /usr/local/bin
BIN_DIR = bin
LDFLAGS = -ldflags "-X main.Version=$(VERSION)"

# Default target
all: build

## build: Build binary for current platform
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/day-night-cycle
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME)"

## build-darwin-amd64: Build for macOS Intel (amd64)
build-darwin-amd64:
	@echo "Building $(BINARY_NAME)-darwin-amd64 $(VERSION)..."
	@mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/day-night-cycle
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64"

## build-darwin-arm64: Build for macOS Apple Silicon (arm64)
build-darwin-arm64:
	@echo "Building $(BINARY_NAME)-darwin-arm64 $(VERSION)..."
	@mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/day-night-cycle
	@echo "Built: $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64"

## build-all: Build for all supported platforms
build-all: build-darwin-amd64 build-darwin-arm64
	@echo "All builds complete!"

## install: Install binary to system path
install: build
	@echo "Installing to $(INSTALL_PATH)..."
	install -m 755 $(BIN_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installed: $(INSTALL_PATH)/$(BINARY_NAME)"

## release: Create GitHub release with binaries
release: build-all
	@echo "Creating GitHub release $(VERSION)..."
	@if [ "$(VERSION)" = "dev" ]; then \
		echo "Error: Cannot release 'dev' version. Create a git tag first."; \
		exit 1; \
	fi
	@if ! command -v gh >/dev/null 2>&1; then \
		echo "Error: GitHub CLI (gh) is not installed"; \
		echo "Install it with: brew install gh"; \
		exit 1; \
	fi
	gh release create $(VERSION) \
		$(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 \
		$(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 \
		--title "$(VERSION)" \
		--generate-notes
	@echo "Release $(VERSION) created successfully!"

## clean: Remove built binaries
clean:
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)
	@echo "Clean complete"

## test: Run tests
test:
	@echo "Running tests..."
	go test -v ./...

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	go fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## mod: Tidy go.mod
mod:
	@echo "Tidying go.mod..."
	go mod tidy

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/  /'
