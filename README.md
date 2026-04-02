# MyApp — Go Clean Architecture Backend Template

A production-ready **Go** backend template using **Clean Architecture**, Chi router, PostgreSQL (sqlx), structured logging (slog), and Docker support.

## Architecture

```
cmd/api/main.go              → Entry point, wiring, graceful shutdown
internal/
├── config/                  → Configuration (koanf, .env + env vars)
├── domain/                  → Entities, errors (zero external deps)
├── repository/              → Interfaces + PostgreSQL (sqlx) implementations
├── service/                 → Business logic
└── handler/                 → HTTP handlers (Chi), middleware, response helpers
migrations/                  → SQL migration files (golang-migrate)
```

**Dependency flow**: `domain` ← `service` ← `handler` ← `cmd/api`

## Tech Stack

- **Go 1.26** with `internal/` package privacy
- **Chi v5** — lightweight, idiomatic HTTP router
- **PostgreSQL + sqlx** — raw SQL with type-safe scanning
- **koanf v2** — configuration from `.env` + environment variables
- **slog** — structured logging (Go standard library)
- **golang-migrate** — SQL-based database migrations
- **testify** — assertion library for tests
- **golangci-lint** — linter aggregator
- **air** — hot reload for development
- **Docker** — multi-stage build + Docker Compose

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [PostgreSQL 15+](https://www.postgresql.org/) (or use Docker Compose)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) (for migrations)

## Getting Started

```bash
# 1. Clone and enter the project
cd go

# 2. Copy environment config
cp .env.example .env

# 3. Start PostgreSQL (via Docker Compose)
make docker-up

# 4. Apply database migrations
make migrate-up

# 5. Run the server
make run

# Server starts at http://localhost:8080
# Health check: GET http://localhost:8080/health
```

## Key Commands

```bash
make build           # Build binary → bin/myapp
make run             # Run the API server
make dev             # Run with hot reload (air)
make test            # Run all tests with -race
make test-coverage   # Tests with HTML coverage report
make lint            # Run golangci-lint
make vet             # Run go vet
make fmt             # Format all Go files
make migrate-up      # Apply pending migrations
make migrate-down    # Rollback last migration
make migrate-create  # Create new migration pair (prompts for name)
make docker-up       # Start app + PostgreSQL containers
make docker-down     # Stop containers
make clean           # Remove build artifacts
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness check |
| GET | `/readiness` | Readiness check (includes DB ping) |
| POST | `/api/v1/items` | Create a new item |
| GET | `/api/v1/items` | List items (paginated: `?page=1&page_size=20`) |
| GET | `/api/v1/items/{id}` | Get item by ID |
| PUT | `/api/v1/items/{id}` | Update an item |
| DELETE | `/api/v1/items/{id}` | Delete an item |

## Adding a New Feature

1. **Domain** — Create entity + input structs in `internal/domain/`
2. **Repository** — Add interface to `internal/repository/interfaces.go`, implement in `postgres/`
3. **Service** — Create service in `internal/service/`, accepts repo interface
4. **Handler** — Create handler in `internal/handler/`, register routes via `RegisterXxxRoutes()`
5. **Main** — Wire dependencies in `cmd/api/main.go`
6. **Migration** — `make migrate-create`, write SQL, `make migrate-up`
7. **Tests** — Unit tests for domain + service, handler tests with httptest

## Project Conventions

- **`(value, error)` returns** — Go's native error handling, no generic Result types
- **`AppError` struct** — structured errors with HTTP status codes and sentinel errors
- **Manual dependency wiring** — no DI framework, everything in `main.go`
- **`internal/` packages** — Go compiler prevents external access
- **Table-driven tests** — `t.Run()` with subtests
- **Configuration** — `.env` file + environment variable overrides

## Docker

```bash
# Start everything (app + PostgreSQL)
docker compose up -d --build

# View logs
docker compose logs -f app

# Stop
docker compose down
```

## License

MIT