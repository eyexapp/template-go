# AGENTS.md — Go Clean Architecture Backend

## Project Identity

| Key | Value |
|-----|-------|
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

---

## Architecture — Clean Architecture (Strict Layers)

```
cmd/api/
└── main.go                  ← COMPOSITION ROOT: wiring, server start, graceful shutdown

internal/
├── config/config.go         ← INFRASTRUCTURE: koanf-based config loading
├── domain/                  ← DOMAIN (ZERO external deps — stdlib only)
│   ├── item.go              ← Entity struct + input DTOs + Validate()
│   └── errors.go            ← AppError + sentinel errors
├── repository/
│   ├── interfaces.go        ← CONTRACTS: Repository interfaces
│   └── postgres/
│       ├── postgres.go      ← DB connection + pool config
│       └── item_repo.go     ← sqlx implementation of interfaces
├── service/
│   └── item_service.go      ← BUSINESS LOGIC: accepts repo interface
└── handler/
    ├── handler.go           ← Shared handler struct
    ├── item_handler.go      ← HTTP handlers (CRUD)
    ├── health_handler.go    ← Health/readiness checks
    ├── middleware/           ← Logger, Recoverer
    └── response/json.go     ← JSON response helpers
```

### Dependency Direction (INWARD ONLY)

```
domain  ←  service  ←  handler  ←  cmd/api
  ↑          ↑
  └── repository (interfaces in repository/, implementations in postgres/)
```

### Strict Layer Rules

| Layer | Can Import From | NEVER Imports |
|-------|----------------|---------------|
| `domain/` | Go standard library ONLY | nothing from `internal/` |
| `repository/` | domain/ | service/, handler/ |
| `service/` | domain/, repository/ (interfaces only) | handler/ |
| `handler/` | service/, domain/ | repository/postgres/ |
| `cmd/api/` | everything (composition root) | — |

---

## Adding New Code — Where Things Go

### New Feature Checklist
1. **Domain**: `internal/domain/product.go` — entity struct + input struct + `Validate()` method
2. **Repository interface**: Add to `internal/repository/interfaces.go`
3. **Repository impl**: `internal/repository/postgres/product_repo.go`
4. **Service**: `internal/service/product_service.go` — accepts interface, returns domain types
5. **Handler**: `internal/handler/product_handler.go` — `RegisterProductRoutes(r chi.Router)`
6. **Wire in main.go**: Create repo → service → handler, call `RegisterProductRoutes`
7. **Migration**: `migrations/` via `make migrate-create`
8. **Tests**: Unit test domain validation + service logic + handler HTTP routing

### Domain Entity Pattern
```go
// internal/domain/product.go
package domain

import "time"

type Product struct {
    ID        string    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Price     float64   `json:"price" db:"price"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateProductInput struct {
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

func (c CreateProductInput) Validate() error {
    if c.Name == "" {
        return NewAppError(ErrValidation, "name is required")
    }
    if c.Price <= 0 {
        return NewAppError(ErrValidation, "price must be positive")
    }
    return nil
}
```

### Handler Pattern
```go
func (h *Handler) RegisterProductRoutes(r chi.Router) {
    r.Route("/api/v1/products", func(r chi.Router) {
        r.Post("/", h.CreateProduct)
        r.Get("/", h.ListProducts)
        r.Get("/{id}", h.GetProduct)
        r.Put("/{id}", h.UpdateProduct)
        r.Delete("/{id}", h.DeleteProduct)
    })
}
```

### Wiring in main.go
```go
// cmd/api/main.go — composition root (manual DI)
productRepo := postgres.NewProductRepo(db)
productService := service.NewProductService(productRepo, logger)
h := handler.NewHandler(productService, logger)
h.RegisterProductRoutes(r)
```

---

## Design & Architecture Principles

### Go Idioms
- `context.Context` as first parameter of every I/O function
- `(value, error)` returns — NEVER use panic for expected failures
- Short variable names for receivers (`s` for service, `r` for repo)
- Interfaces where consumed (not where implemented)

### Constructor Pattern
```go
func NewProductService(repo repository.ProductRepository, logger *slog.Logger) *ProductService {
    return &ProductService{repo: repo, logger: logger}
}
```

### No DI Framework
- `cmd/api/main.go` is the composition root
- Manual wiring: create repo → pass to service → pass to handler
- This is intentional — Go prefers explicit over magic

---

## Error Handling

### AppError + Sentinels
```go
// internal/domain/errors.go
var (
    ErrNotFound   = errors.New("not found")
    ErrConflict   = errors.New("conflict")
    ErrValidation = errors.New("validation")
)

type AppError struct {
    Code    error
    Message string
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Is(target error) bool { return errors.Is(e.Code, target) }
```

### Handler Error Mapping
```go
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
    product, err := h.productService.GetByID(r.Context(), chi.URLParam(r, "id"))
    if err != nil {
        response.Error(w, err) // AppError → correct HTTP status + JSON
        return
    }
    response.JSON(w, http.StatusOK, product)
}
```

---

## Code Quality

### Naming Conventions
| Artifact | Convention | Example |
|----------|-----------|---------|
| Package | `lowercase` (no underscores) | `domain`, `postgres` |
| File | `snake_case.go` | `item_handler.go` |
| Exported | `MixedCaps` | `ProductService` |
| Unexported | `mixedCaps` | `productRepo` |
| Interface | Verb + `-er` suffix | `ProductRepository` |
| Constructor | `New` + type name | `NewProductService` |
| Struct tags | `json` + `db` | `json:"name" db:"name"` |

### Code Style
- `gofmt` / `goimports` — non-negotiable formatting
- `golangci-lint` — enforced linting
- `make vet` before commit
- Short functions — if > 40 lines, extract
- Early returns — avoid deep nesting

---

## Testing Strategy

| Level | What | Where | Tool |
|-------|------|-------|------|
| Unit | Domain validation, service logic | `*_test.go` beside source | Go testing + testify |
| Handler | HTTP routing, status codes | `handler/*_test.go` | httptest.NewRecorder |
| Integration | Full API with DB | `tests/` | Docker Compose + testify |

### Table-Driven Tests — Mandatory Pattern
```go
func TestCreateProductInput_Validate(t *testing.T) {
    tests := []struct {
        name    string
        input   domain.CreateProductInput
        wantErr bool
    }{
        {"valid", domain.CreateProductInput{Name: "Test", Price: 10}, false},
        {"empty name", domain.CreateProductInput{Name: "", Price: 10}, true},
        {"zero price", domain.CreateProductInput{Name: "Test", Price: 0}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.input.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Manual Mocks (No frameworks)
```go
type mockProductRepo struct {
    products []domain.Product
    err      error
}

func (m *mockProductRepo) GetByID(ctx context.Context, id string) (*domain.Product, error) {
    if m.err != nil { return nil, m.err }
    // ...
}
```

---

## Security & Performance

### Security
- All user input validated in domain layer (`Validate()` methods)
- Never log sensitive data (tokens, passwords)
- Environment variables for secrets (koanf + .env)
- SQL queries via sqlx placeholders — NEVER string concatenation

### Performance
- sqlx named queries — prepared statement reuse
- `context.Context` with timeouts for all DB operations
- Graceful shutdown via `signal.NotifyContext` in main.go
- Multi-stage Docker build — minimal production image
- Connection pool tuning in `internal/repository/postgres/postgres.go`

---

## Commands

| Action | Command |
|--------|---------|
| Build | `make build` |
| Run | `make run` |
| Dev (hot reload) | `make dev` |
| Test (race) | `make test` |
| Coverage | `make test-coverage` |
| Lint | `make lint` |
| Format | `make fmt` |
| Migrate up | `make migrate-up` |
| Migrate down | `make migrate-down` |
| New migration | `make migrate-create` |
| Docker up | `make docker-up` |

---

## Prohibitions — NEVER Do These

1. **NEVER** import outer layers from inner layers (domain must not import handler)
2. **NEVER** use `panic` for expected errors — return `(value, error)`
3. **NEVER** use a DI framework — manual wiring in `cmd/api/main.go`
4. **NEVER** use mock frameworks — write manual mock structs
5. **NEVER** use `interface{}` / `any` when a concrete type is known
6. **NEVER** skip `context.Context` as first parameter in I/O functions
7. **NEVER** concatenate SQL strings — use sqlx placeholders only
8. **NEVER** use global variables for state — pass dependencies explicitly
9. **NEVER** use `init()` functions — explicit initialization in main
10. **NEVER** ignore linter warnings — `golangci-lint` is the law
