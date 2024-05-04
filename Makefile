## ----------------------------------------------------------------------
## Here are the commands you can use for running, testing and developing
## the application. Please note that golangci-lint and flyway CLI need to
## be installed by you until I create an init script. Refer to docs.
## PS: I don't use .PHONY because it's not needed.
## ----------------------------------------------------------------------

build:		## Build binary.
build: clean $(BINARY_PATH) install
	$(call print-target)
	@echo "Building $(BINARY_NAME)"
	@for os in $(PLATFORMS); do \
		echo "Building for $$os/$(ARCH)"; \
		GOARCH=$(ARCH) GOOS=$$os go build -o $(BINARY_PATH)/$(BINARY_NAME)-$$os; \
	done

run:		## Run the application
run: build
	$(call print-target)
	@echo "Running $(BINARY_NAME)"
	@$(BINARY_PATH)/$(BINARY_NAME)-darwin || $(BINARY_PATH)/$(BINARY_NAME)-linux || $(BINARY_PATH)/$(BINARY_NAME)-windows

test:		## Run tests (DB should be up and running)
test: clean install format lint
	$(call print-target)
	@echo "Running tests"
	@go test ./... -v -coverprofile=coverage.out

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

lint:		## Lint using golangci-lint (Need to install it).
	$(call print-target)
	@echo "Running linters..."
	@golangci-lint run ./...

migrate-up:	## Run database migrations using Flyway (Need to install it).
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