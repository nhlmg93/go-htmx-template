.PHONY: all build run clean

# Default target
all: build

# Build the application executable into the build directory
build:
	@mkdir -p build
	go build -o build/foo ./main.go

# Run the application
run: build
	./build/foo

dev:
	air

# Clean up generated files
clean:
	rm -rf build
