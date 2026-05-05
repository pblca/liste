.PHONY: build clean install test

BINARY_NAME=liste
VERSION?=0.1.0

# Build for current platform
build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BINARY_NAME) .

# Build for all platforms
build-all:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME)-windows-amd64.exe .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME)-linux-amd64 .

# Install to GOPATH/bin
install:
	go install -ldflags="-s -w -X main.version=$(VERSION)" .

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-*
