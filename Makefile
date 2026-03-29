BINARY  := tztui
CMD     := ./cmd/tztui
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build test lint install clean

build:
	go build $(LDFLAGS) -o $(BINARY) $(CMD)

test:
	go test ./...

lint:
	go vet ./...

install:
	go install $(LDFLAGS) $(CMD)

clean:
	rm -f $(BINARY)
