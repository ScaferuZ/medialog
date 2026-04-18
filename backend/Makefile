BIN_DIR := bin
BINARY := $(BIN_DIR)/server

.PHONY: dev build test clean migrate-up migrate-down migrate-create sqlc-generate db-reset

dev:
	SERVER_PORT?=8080 go run ./cmd/server

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BINARY) ./cmd/server

test:
	go test ./...

clean:
	rm -rf $(BIN_DIR)

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create:
	@read -p "Migration name: " name; \
		migrate create -ext sql -dir migrations -seq $$name

sqlc-generate:
	sqlc generate

db-reset: migrate-down migrate-up
