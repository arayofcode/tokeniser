BINARY_NAME=tokeniser
BINARY_PATH=bin/${BINARY_NAME}
DB_URL ?= $(DATABASE_URL)

.PHONY: build run test clean dep vet lint migrate-up migrate-down

build:
	@go build -v -o ${BINARY_PATH} ./...

run: dep build
	@./$(BINARY_PATH)

test:
	go test ./... -v -coverprofile=coverage.out

clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f $(BINARY_PATH)

dep:
	@go mod download
	@go mod verify

vet:
	@go vet -v ./...

format: vet

migrate-up:
	migrate -path database/migration -database "$(DB_URL)" up

migrate-down:
	migrate -path database/migration -database "$(DB_URL)" down
