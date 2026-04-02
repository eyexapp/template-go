# AI Coding Instructions — MyApp (Go Clean Architecture Backend)

## Project Overview

This is a **Go Clean Architecture** backend template with layered separation:

- **cmd/api/** — Entry point, HTTP server bootstrap, dependency wiring
- **internal/domain/** — Entities, value objects, domain errors (zero external deps)
- **internal/repository/** — Repository interfaces + PostgreSQL implementations (sqlx)
- **internal/service/** — Business logic layer
- **internal/handler/** — HTTP handlers (Chi router), middleware, JSON response helpers
- **internal/config/** — Configuration loading (koanf, .env + environment variables)

## Architecture Rules (MUST FOLLOW)

1. **Dependency direction**: domain ← service ← handler ← cmd/api (inner layers NEVER import outer layers)
2. **`internal/domain/` has ZERO external dependencies** — only Go standard library
3. **Service layer accepts repository interfaces**, not concrete implementations
4. **`cmd/api/main.go` is the composition root** — manual dependency wiring (no DI framework)
5. **Use Go `(value, error)` returns** — never use generic Result types or panic for expected failures
6. **All packages live under `internal/`** — Go compiler enforces access control

## Code Conventions

- **Naming**: Go standard — `MixedCaps` exported, `mixedCaps` unexported, short receiver names
- **Errors**: Use `AppError` struct with sentinel errors (`ErrNotFound`, `ErrConflict`, `ErrValidation`)
- **Interfaces**: Define where they're consumed (repository interfaces in `internal/repository/`)
- **Context**: First parameter of every I/O function: `func Foo(ctx context.Context, ...) error`
- **Struct tags**: `json` for API, `db` for sqlx, `koanf` for config
- **Testing**: Table-driven tests, `testify/assert`, manual mock structs (no mock frameworks)

## Adding a New Feature (Checklist)

1. **Domain**: Create entity struct in `internal/domain/<entity>.go` with JSON + DB tags
2. **Domain**: Add validation method on the input struct (`func (c CreateInput) Validate() error`)
3. **Repository**: Add interface in `internal/repository/interfaces.go`
4. **Repository**: Implement in `internal/repository/postgres/<entity>_repo.go`
5. **Service**: Create `internal/service/<entity>_service.go` — accepts repo interface, returns domain types
6. **Handler**: Create `internal/handler/<entity>_handler.go` — register routes via `RegisterXxxRoutes(r chi.Router)`
7. **Main**: Wire in `cmd/api/main.go` — create repo → service → handler, call `RegisterXxxRoutes`
8. **Migration**: Create SQL files in `migrations/` using `make migrate-create`
9. **Tests**: Unit tests for domain validation + service logic, handler tests for HTTP routing/parsing

## Technology Stack

| Component | Technology |
|---|---|
| Language | Go 1.26 |
| Router | Chi v5 |
| Database | PostgreSQL + sqlx |
| Config | koanf v2 (.env + env vars) |
| Logging | slog (standard library) |
| Migrations | golang-migrate |
| Testing | Go testing + testify/assert |
| Linting | golangci-lint |
| Hot Reload | air |
| Container | Docker (multi-stage) + Docker Compose |

## Key Commands

```bash
make build          # Build binary to bin/myapp
make run            # Run the API server
make dev            # Run with hot reload (air)
make test           # Run all tests with race detector
make test-coverage  # Run tests with coverage report
make lint           # Run golangci-lint
make vet            # Run go vet
make fmt            # Format all Go files
make migrate-up     # Apply pending migrations
make migrate-down   # Rollback last migration
make migrate-create # Create a new migration pair
make docker-up      # Start app + PostgreSQL via Docker Compose
make docker-down    # Stop Docker Compose
```

## Project Structure

```
go/
├── cmd/api/main.go              # Entry point, wiring, graceful shutdown
├── internal/
│   ├── config/config.go         # koanf-based config loading
│   ├── domain/
│   │   ├── item.go              # Item entity + input DTOs + validation
│   │   └── errors.go            # AppError, sentinel errors
│   ├── repository/
│   │   ├── interfaces.go        # Repository interface definitions
│   │   └── postgres/
│   │       ├── postgres.go      # DB connection + pool config
│   │       └── item_repo.go     # sqlx implementation
│   ├── service/
│   │   └── item_service.go      # Business logic
│   └── handler/
│       ├── handler.go           # Shared handler struct
│       ├── item_handler.go      # CRUD endpoints
│       ├── health_handler.go    # Health + readiness checks
│       ├── middleware/           # Logger, Recoverer
│       └── response/json.go     # JSON response helpers
├── migrations/                  # SQL migration files
├── Dockerfile                   # Multi-stage build
├── docker-compose.yml           # App + PostgreSQL
├── Makefile                     # All build/test/run commands
├── go.mod                       # Module definition
└── .env.example                 # Environment variable template
```

## Common Patterns

### Error Handling
```go
// Service returns domain errors
item, err := s.repo.GetByID(ctx, id)
if err != nil {
    return nil, err // AppError propagates to handler
}

// Handler maps errors to HTTP responses
item, err := h.ItemService.GetByID(r.Context(), id)
if err != nil {
    response.Error(w, err) // AppError → correct status code + JSON
    return
}
```

### Route Registration
```go
func (h *Handler) RegisterItemRoutes(r chi.Router) {
    r.Route("/api/v1/items", func(r chi.Router) {
        r.Post("/", h.CreateItem)
        r.Get("/", h.ListItems)
        r.Get("/{id}", h.GetItem)
        r.Put("/{id}", h.UpdateItem)
        r.Delete("/{id}", h.DeleteItem)
    })
}
```

### Dependency Wiring (manual, in main.go)
```go
itemRepo := postgres.NewItemRepo(db)
itemService := service.NewItemService(itemRepo, logger)
h := handler.NewHandler(itemService, logger)
h.RegisterItemRoutes(r)
```
