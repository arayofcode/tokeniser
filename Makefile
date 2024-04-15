# Check only for .env file (For development)
# .env will not be used in production
include $(wildcard .env)

.DEFAULT_GOAL := help

BINARY_NAME=tokeniser
BINARY_PATH=$(CURDIR)/bin
PLATFORMS=darwin linux windows
ARCH=amd64
FLYWAY_URL=jdbc:postgresql://$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)

$(BINARY_PATH):
	@mkdir -p $@

.PHONY: build
build: clean dep format lint test $(BINARY_PATH)
	$(call print-target)
	@echo "Building $(BINARY_NAME)"
	@for os in $(PLATFORMS); do \
		echo "Building for $$os/$(ARCH)"; \
		GOARCH=$(ARCH) GOOS=$$os go build -o $(BINARY_PATH)/$(BINARY_NAME)-$$os; \
	done

.PHONY: run
run: build migrate-up
	$(call print-target)
	@echo "Running $(BINARY_NAME)"
	@$(BINARY_PATH)/$(BINARY_NAME)-darwin || $(BINARY_PATH)/$(BINARY_NAME)-linux || $(BINARY_PATH)/$(BINARY_NAME)-windows

.PHONY: start
start:
	$(call print-target)
	@docker compose up --remove-orphans -d

.PHONY: start-build
start-build:
	$(call print-target)
	@docker compose up --build --remove-orphans --force-recreate -d 

.PHONY: stop
stop:
	$(call print-target)
	@docker compose stop

.PHONY: dev
dev:
	$(call print-target)
	@docker compose -f compose.dev.yaml up -d

.PHONY: dev-build
dev-build:
	$(call print-target)
	@docker compose -f compose.dev.yaml up --build -d

.PHONY: dev-watch
dev-watch:
	$(call print-target)
	@docker compose -f compose.dev.yaml watch

.PHONY: stop-dev
stop-dev:
	$(call print-target)
	@docker compose -f compose.dev.yaml stop

.PHONY: test
test: lint migrate-up
	$(call print-target)
	@echo "Running tests"
	@go test ./... -v -coverprofile=coverage.out

.PHONY: clean
clean:
	$(call print-target)
	@echo "Cleaning up"
	@go clean
	@rm -rf $(BINARY_PATH) coverage.out

.PHONY: dep
dep:
	$(call print-target)
	@echo "Verifying and downloading dependencies"
	@go mod tidy
	@go mod download
	@go mod verify

.PHONY: vet
vet:
	$(call print-target)
	@go vet -v ./...

.PHONY: format
format: vet
	$(call print-target)
	@echo "Formatting code"
	@go fmt ./...

.PHONY: migrate-up
migrate-up:
	$(call print-target)
	@echo "Setting up database"
	@docker compose -f compose.yaml up db -d
	@echo "Running latest database migrations"
	@flyway -url=$(FLYWAY_URL) -user=$(POSTGRES_USER) -password=$(POSTGRES_PASSWORD) -locations=migrations/migration migrate

migrate-up-dev:
	$(call print-target)
	@echo "Setting up database"
	@docker compose -f compose.dev.yaml up db-dev -d
	@echo "Running latest database migrations"
	@flyway -url=$(FLYWAY_URL) -user=$(POSTGRES_USER) -password=$(POSTGRES_PASSWORD) -locations=migrations/migration migrate

.PHONY: migrate-down
migrate-down:
	$(call print-target)
	@echo "Need to find a workaround for this given Flyway Community doesn't provide this."

.PHONY: lint
lint: format
	$(call print-target)
	@echo "Running linters..."
	@golangci-lint run ./...

.PHONY: help
help:
	@echo "Available commands:\n"
	@echo "Following commands use docker compose:"
	@echo "  start         Start the application."
	@echo "  start-build   Rebuild all images and start the application."
	@echo "  stop		   Stop the running containers for production app."
	@echo "  dev		   Run the local development setup."
	@echo "  dev-build	   Run the local development setup by rebuilding containers."
	@echo "  dev-watch	   Run the local development setup. Allows hot reloading."
	@echo "  stop-dev	   Stop the local development setup."

	@echo "\nFollowing commands run within your shell. Setup environment variables before proceeding:"
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

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef
