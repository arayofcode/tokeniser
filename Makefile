## ----------------------------------------------------------------------
## Here are the commands you can use for running, testing and developing
## the application. Please note that golangci-lint and flyway CLI need to
## be installed by the user manually. Refer to their docs.
## PS: I don't use .PHONY because it's not needed.
## ----------------------------------------------------------------------
start:		## Run application. Reuse any existing containers.
	$(call print-target)
	@echo "Running application using docker compose. Reusing any existing containers"
	@docker compose -f compose.yaml up --no-recreate -d

start-clean:	## Rebuild containers and run the application.
	$(call print-target)
	@echo "Running application using docker compose and rebuilding containers"
	@docker compose -f compose.yaml up --build -d

stop:		## Run any existing instance of the application.
	$(call print-target)
	@echo "Stopping any instances of the application running"
	@docker compose -f compose.yaml stop

dev:		## Run application in development mode.
# || true suppresses make: *** [dev] Error 130
	$(call print-target)
	@echo "Running application in dev environment"
	@docker compose up --watch || true

build:		## Build binary.
build: clean $(BINARY_PATH) install
	$(call print-target)
	@echo "Building $(BINARY_NAME)"
	@for os in $(PLATFORMS); do \
		echo "Building for $$os/$(ARCH)"; \
		GOARCH=$(ARCH) GOOS=$$os go build -o $(BINARY_PATH)/$(BINARY_NAME)-$$os; \
	done

run:		## Run the application.
run: build
	$(call print-target)
	@echo "Running $(BINARY_NAME)"
	@$(BINARY_PATH)/$(BINARY_NAME)-darwin || $(BINARY_PATH)/$(BINARY_NAME)-linux || $(BINARY_PATH)/$(BINARY_NAME)-windows

test:		## Run tests (DB should be up and running).
test: clean install format lint
	$(call print-target)
	@echo "Running tests"
	@go test ./... -v -coverprofile=coverage.out

test-ci:	## Run tests in CI using docker compose file
	$(call print-target)
	@echo "Running tests"
	@docker compose -f compose.test.yaml up app --build --remove-orphans

install:	## Install dependencies.
	$(call print-target)
	@echo "Verifying and downloading dependencies"
	@go mod tidy
	@go mod download
	@go mod verify

clean:		## Remove temporary files and previous binary.
	$(call print-target)
	@go clean
	@rm -rf $(BINARY_PATH) coverage.out

format:		## Format code.
	$(call print-target)
	@echo "Formatting code"
	@go fmt ./...

lint:		## Lint using golangci-lint.
	$(call print-target)
	@echo "Running linters..."
	@golangci-lint run ./...

migrate-up:	## Run database migrations using Flyway.
	$(call print-target)
	@echo "Running latest database migrations"
	@flyway -url=$(FLYWAY_URL) -user=$(POSTGRES_USER) -password=$(POSTGRES_PASSWORD) -locations=migrations/migration migrate

migrate-down:	## Undo the last database migration (Not yet implemented).
	$(call print-target)
	@echo "Need to find a workaround for this given Flyway Community doesn't provide this."

help:		## Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

BINARY_NAME=tokeniser
BINARY_PATH=$(CURDIR)/bin
PLATFORMS=darwin linux windows
ARCH=amd64
FLYWAY_URL=jdbc:postgresql://$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)

$(BINARY_PATH):
	@mkdir -p $@

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef