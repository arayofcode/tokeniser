# Check only for .env file (For development)
# .env will not be used in production
-include $(wildcard .env)

BINARY_NAME=tokeniser
BINARY_PATH=$(CURDIR)/bin/

DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
FLYWAY_URL=jdbc:postgresql://$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)

.PHONY: all build run test clean dep vet format migrate-up migrate-down lint help

build: clean dep $(BINARY_PATH)
	@echo "Building $(BINARY_NAME)"
	@go build -v -o $(BINARY_PATH) ./...

$(BINARY_PATH):
	@mkdir -p $@

run: dep build migrate-up
	@echo "Running $(BINARY_NAME)"
	@DB=$(DB_URL) PASSPHRASE=$(PASSPHRASE) $(BINARY_PATH)/$(BINARY_NAME)

start:
	@docker compose up -d

start-clean:
	@docker compose up --build -d

test:
	@echo "Running tests"
	@DB=$(DB_URL) go test ./... -v -coverprofile=coverage.out

clean:
	@echo "Cleaning up"
	@go clean
	@rm -rf $(BINARY_PATH) coverage.out

dep:
	@echo "Verifying and downloading dependencies"
	@go mod tidy
	@go mod download
	@go mod verify

vet:
	@go vet -v ./...

format: vet
	@echo "Formatting code"
	@go fmt ./...

migrate-up:
	@echo "Running latest database migrations"
	@flyway -url=$(FLYWAY_URL) -user=$(POSTGRES_USER) -password=$(POSTGRES_PASSWORD) migrate

migrate-down:
	@echo "Need to find a workaround for this given Flyway Community doesn't provide this."

lint:
	@echo "Running linters..."
	@golangci-lint run ./...

help:
	@echo "Available commands:"
	@echo "  start         Start the application."
	@echo "  start-clean   Rebuild all images and start the application."
	@echo "  build         Build the application."
	@echo "  run           Run the application."
	@echo "  dep           Download and verify dependencies."
	@echo "  test          Run tests."
	@echo "  lint          Run lint check."
	@echo "  clean         Clean up the project."
	@echo "  vet           Run go vet."
	@echo "  format        Format the code."
	@echo "  migrate-up    Apply database migrations."
	@echo "  migrate-down  Undo the last database migration (Not yet implemented)."