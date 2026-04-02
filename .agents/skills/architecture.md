---
name: architecture
type: knowledge
version: 1.0.0
agent: CodeActAgent
triggers:
  - architecture
  - chi
  - clean architecture
  - handler
  - repository
  - middleware
---

# Architecture — Go (Chi + Clean Architecture)

## Chi Router

```go
func main() {
    cfg := config.Load()
    db := database.Connect(cfg.DatabaseURL)
    defer db.Close()

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{cfg.FrontendURL}}))

    r.Route("/api", func(r chi.Router) {
        r.Mount("/users", user.NewHandler(user.NewService(user.NewRepo(db))).Routes())
        r.Mount("/auth", auth.NewHandler(auth.NewService(db)).Routes())
    })

    slog.Info("server starting", "port", cfg.Port)
    http.ListenAndServe(":"+cfg.Port, r)
}
```

## Clean Architecture

```
cmd/
└── server/
    └── main.go         ← Entry point
internal/
├── config/             ← Environment config
├── database/           ← DB connection, migrations
├── user/               ← Feature package (self-contained)
│   ├── handler.go      ← HTTP handlers
│   ├── service.go      ← Business logic
│   ├── repository.go   ← Database queries (sqlx)
│   ├── model.go        ← Domain types
│   └── dto.go          ← Request/Response DTOs
├── auth/
├── middleware/          ← Auth, logging middleware
└── pkg/                ← Shared utilities
    ├── httputil/        ← JSON response helpers
    └── validate/        ← Input validation
```

## Feature Package Pattern

```go
// internal/user/handler.go
type Handler struct {
    service *Service
}

func NewHandler(s *Service) *Handler {
    return &Handler{service: s}
}

func (h *Handler) Routes() chi.Router {
    r := chi.NewRouter()
    r.Get("/", h.List)
    r.Post("/", h.Create)
    r.Get("/{id}", h.GetByID)
    r.Put("/{id}", h.Update)
    r.Delete("/{id}", h.Delete)
    return r
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    var dto CreateUserDTO
    if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
        httputil.Error(w, http.StatusBadRequest, "invalid JSON")
        return
    }
    user, err := h.service.Create(r.Context(), dto)
    if err != nil {
        httputil.HandleError(w, err)
        return
    }
    httputil.JSON(w, http.StatusCreated, user)
}
```

## Database — sqlx

```go
func (r *Repo) FindByID(ctx context.Context, id string) (*User, error) {
    var user User
    err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
    if err == sql.ErrNoRows {
        return nil, ErrNotFound
    }
    return &user, err
}
```

## Rules

- Feature packages: each domain is self-contained (`internal/user/`).
- Handler → Service → Repository — clean dependency flow.
- Context propagation: pass `context.Context` through all layers.
- No global state — dependency injection via constructors.
