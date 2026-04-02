---
name: version-control
type: knowledge
version: 1.0.0
agent: CodeActAgent
triggers:
  - git
  - commit
  - ci
  - makefile
  - deploy
---

# Version Control — Go

## Commits (Conventional)

- `feat(user): add email verification`
- `fix(auth): handle expired refresh tokens`
- `refactor(handler): extract JSON helper package`

## Makefile

```makefile
.PHONY: run build test lint migrate

run:
	go run cmd/server/main.go

build:
	CGO_ENABLED=0 go build -o bin/server cmd/server/main.go

test:
	go test -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1
```

## CI Pipeline

1. `go mod download`
2. `golangci-lint run` — lint
3. `go test -race -coverprofile=coverage.out ./...` — test
4. `go build -o bin/server cmd/server/main.go` — build

## .gitignore

```
bin/
vendor/
.env
coverage.out
```

## Docker

```dockerfile
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server cmd/server/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=build /server /server
EXPOSE 8080
CMD ["/server"]
```

## Module Management

- `go mod tidy` — clean up unused dependencies.
- `go.sum` is committed — integrity verification.
- Vendor optional: `go mod vendor` for offline builds.
