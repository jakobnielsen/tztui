BINARY  := tztui
CMD     := ./cmd/tztui
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
DIST    := dist

.PHONY: build test lint install clean \
        build-linux-amd64 build-linux-arm64 \
        build-darwin-amd64 build-darwin-arm64 \
        build-all help

.DEFAULT_GOAL := help

## help: Show this help message.
help:
	@echo "Usage: make <target>"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ { printf "  %-24s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""

build: ## Build for the host platform → ./tztui
	go build $(LDFLAGS) -o $(BINARY) $(CMD)

# ── Cross-compilation targets ─────────────────────────────────────────────────

build-linux-amd64: ## Cross-compile for Linux / x86-64.
	GOOS=linux  GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-amd64   $(CMD)

build-linux-arm64: ## Cross-compile for Linux / ARM64.
	GOOS=linux  GOARCH=arm64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-arm64   $(CMD)

build-darwin-amd64: ## Cross-compile for macOS / x86-64.
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-darwin-amd64  $(CMD)

build-darwin-arm64: ## Cross-compile for macOS / Apple Silicon.
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-darwin-arm64  $(CMD)

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 ## Build all platform/arch combinations into dist/.

# ── Standard targets ─────────────────────────────────────────────────────────

test: ## Run all tests.
	go test ./...

lint: ## Run go vet.
	go vet ./...

install: ## Install tztui to GOPATH/bin.
	go install $(LDFLAGS) $(CMD)

clean: ## Remove build artefacts (binary + dist/).
	rm -f $(BINARY)
	rm -rf $(DIST)
