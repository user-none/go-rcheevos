.PHONY: all build test vet fmt clean example submodule coverage

# Default target
all: submodule build vet

# Initialize/update git submodule
submodule:
	@if [ ! -f rcheevos/include/rc_client.h ]; then \
		echo "Initializing rcheevos submodule..."; \
		git submodule update --init --recursive; \
	fi

# Build the library
build: submodule
	go build .

# Build the example
example: submodule
	go build -o bin/example ./example

# Run tests
test: submodule
	go test -v ./...

# Run go vet
vet: submodule
	go vet ./...

# Format code
fmt:
	go fmt ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean .

# Run tests with coverage
coverage: submodule
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@rm -f coverage.out

# Run tests with coverage and open in browser
coverage-html: submodule
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	@rm -f coverage.out

# Update rcheevos submodule to latest
update-rcheevos:
	git submodule update --remote rcheevos
	@echo "Updated rcheevos to:"
	@cd rcheevos && git log -1 --oneline
