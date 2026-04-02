---
name: code-quality
type: knowledge
version: 1.0.0
agent: CodeActAgent
triggers:
  - clean code
  - naming
  - lint
  - error handling
  - interface
---

# Code Quality — Go

## Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Package | lowercase, short | `user`, `httputil` |
| Exported func | PascalCase | `CreateUser()` |
| Unexported | camelCase | `validateEmail()` |
| Interface | -er suffix | `UserReader`, `Validator` |
| Struct | PascalCase | `UserService` |
| Constant | PascalCase or camelCase | `MaxRetries` |
| Error var | Err prefix | `ErrNotFound` |
| File | snake_case | `user_handler.go` |

## Error Handling

```go
// Sentinel errors
var (
    ErrNotFound    = errors.New("not found")
    ErrConflict    = errors.New("already exists")
    ErrUnauthorized = errors.New("unauthorized")
)

// Wrap errors with context
func (s *Service) Create(ctx context.Context, dto CreateDTO) (*User, error) {
    user, err := s.repo.Create(ctx, dto)
    if err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }
    return user, nil
}

// Check errors
if errors.Is(err, ErrNotFound) {
    httputil.Error(w, http.StatusNotFound, "user not found")
}
```

## Interface Pattern

```go
// Small, focused interfaces
type UserReader interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

type UserWriter interface {
    Create(ctx context.Context, dto CreateDTO) (*User, error)
    Update(ctx context.Context, id string, dto UpdateDTO) (*User, error)
}

type UserRepository interface {
    UserReader
    UserWriter
}
```

## Logging — slog

```go
slog.Info("user created", "user_id", user.ID, "email", user.Email)
slog.Error("failed to create user", "error", err, "email", dto.Email)
```

- Structured logging with `slog` (stdlib).
- Key-value pairs — no string formatting.

## Linting — golangci-lint

```yaml
# .golangci.yml
linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - unused
    - gosec
    - gocritic
```

- `golangci-lint run` — comprehensive linting.
- `gofmt` / `goimports` for formatting.
