.PHONY: build run test lint vet fmt clean migrate-up migrate-down migrate-create docker-up docker-down

# Go parameters
BINARY_NAME=myapp
MAIN_PATH=./cmd/api

# Build
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run
run:
	go run $(MAIN_PATH)

# Run with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	air

# Test
test:
	go test ./... -v -race -count=1

test-coverage:
	go test ./... -v -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Lint & Format
lint:
	golangci-lint run ./...

vet:
	go vet ./...

fmt:
	gofmt -s -w .

# Migrate (requires golang-migrate CLI)
MIGRATE_DSN ?= "postgres://postgres:postgres@localhost:5432/myapp?sslmode=disable"

migrate-up:
	migrate -path migrations -database $(MIGRATE_DSN) up

migrate-down:
	migrate -path migrations -database $(MIGRATE_DSN) down 1

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Docker
docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app

# Clean
clean:
	rm -rf bin/ tmp/ coverage.out coverage.html
