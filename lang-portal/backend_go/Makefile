.PHONY: build run test clean

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run with seeding
run-with-seed:
	go run cmd/api/main.go -seed

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
