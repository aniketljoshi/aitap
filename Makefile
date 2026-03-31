BINARY=aitap
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT)"

.PHONY: build run clean test install

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/aitap

run: build
	./bin/$(BINARY)

install:
	CGO_ENABLED=0 go install $(LDFLAGS) ./cmd/aitap

test:
	go test ./... -v

clean:
	rm -rf bin/

# Cross-compile for releases
release:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY)-darwin-arm64 ./cmd/aitap
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY)-darwin-amd64 ./cmd/aitap
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY)-linux-amd64 ./cmd/aitap
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY)-linux-arm64 ./cmd/aitap
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY)-windows-amd64.exe ./cmd/aitap

# Demo: run proxy and a test curl side by side
demo: build
	@echo "Starting aitap..."
	@echo "In another terminal, run:"
	@echo "  export HTTP_PROXY=http://127.0.0.1:9119"
	@echo "  curl http://localhost:11434/api/chat -d '{\"model\":\"llama3\",\"messages\":[{\"role\":\"user\",\"content\":\"hello\"}],\"stream\":false}'"
	./bin/$(BINARY)
