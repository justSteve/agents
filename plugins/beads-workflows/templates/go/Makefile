.PHONY: build test lint clean

# Build the project
build:
	go build -o bin/{{project_name}} ./cmd/{{project_name}}

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
coverage: test
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run

# Run go vet
vet:
	go vet ./...

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run the application
run:
	go run ./cmd/{{project_name}}
