BINARY_NAME=tokeniser
BINARY_PATH=bin/${BINARY_NAME}
PLATFORMS=darwin linux windows
ARCH=amd64

.PHONY: build run test clean dep vet lint

build:
	@for os in $(PLATFORMS); do \
		echo "Building for $$os/$(ARCH)"; \
		GOARCH=$(ARCH) GOOS=$$os go build -o $(BINARY_PATH)-$$os; \
	done

run: build
	@./$(BINARY_PATH)-darwin || ./$(BINARY_PATH)-linux

test:
	go test ./... -v -coverprofile=coverage.out

clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f $(BINARY_PATH)-*

dep:
	@go mod download

vet:
	go vet ./...

format: vet
