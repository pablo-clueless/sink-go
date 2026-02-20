# sink-go Makefile
# Temporary email inbox service

BINARY_NAME=sink-go
MAIN_PATH=.
BUILD_DIR=bin
GO=go

# Version info (can be overridden)
VERSION?=0.1.0
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Ports
API_PORT?=8080
SMTP_PORT?=2525

.PHONY: all build run clean test lint fmt vet tidy deps dev docker-build docker-run help

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/ /'

## all: Build the application
all: build

## build: Compile the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Built: $(BUILD_DIR)/$(BINARY_NAME)"

## build-linux: Cross-compile for Linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

## build-windows: Cross-compile for Windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

## build-all: Build for all platforms
build-all: build build-linux build-windows

## run: Run the application
run:
	@echo "Starting sink-go (API: $(API_PORT), SMTP: $(SMTP_PORT))..."
	$(GO) run $(MAIN_PATH)

## dev: Run with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air

## test: Run tests
test:
	$(GO) test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@mkdir -p $(BUILD_DIR)
	$(GO) test -v -coverprofile=$(BUILD_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "Coverage report: $(BUILD_DIR)/coverage.html"

## bench: Run benchmarks
bench:
	$(GO) test -bench=. -benchmem ./...

## lint: Run golangci-lint (requires: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

## fmt: Format code
fmt:
	$(GO) fmt ./...
	@echo "Code formatted."

## vet: Run go vet
vet:
	$(GO) vet ./...

## tidy: Tidy and verify dependencies
tidy:
	$(GO) mod tidy
	$(GO) mod verify

## deps: Download dependencies
deps:
	$(GO) mod download

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	$(GO) clean
	@echo "Cleaned."

## docker-build: Build Docker image
docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest

## docker-run: Run Docker container
docker-run:
	docker run -p $(API_PORT):8080 -p $(SMTP_PORT):2525 $(BINARY_NAME):latest

## docker-compose: Run with docker-compose
docker-compose:
	docker-compose up --build

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test

## install: Install binary to GOPATH/bin
install:
	$(GO) install $(LDFLAGS) $(MAIN_PATH)

## send-test-email: Send a test email (requires curl and netcat/telnet)
send-test-email:
	@echo "Sending test email to test@sink.io.local..."
	@printf "EHLO localhost\r\nMAIL FROM:<sender@example.com>\r\nRCPT TO:<test@sink.io.local>\r\nDATA\r\nFrom: sender@example.com\r\nTo: test@sink.io.local\r\nSubject: Test Email\r\n\r\nThis is a test email body.\r\n.\r\nQUIT\r\n" | nc localhost $(SMTP_PORT)
	@echo "Test email sent!"
