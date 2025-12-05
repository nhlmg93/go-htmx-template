.PHONY: all build run clean

# Default target
all: build

# Build the application executable into the build directory
build:
	@mkdir -p build
	@go build -o build/main ./main.go

# Run the application
run: build
	./build/main

dev:
	@air

generate:
	@sqlc generate

migrate:
	@mkdir -p build
	@touch  build/dev.db
	@sqlite3 build/dev.db < schema.sql

# Clean up generated files
clean:
	@rm -rf build
