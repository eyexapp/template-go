---
name: security-performance
type: knowledge
version: 1.0.0
agent: CodeActAgent
triggers:
  - security
  - performance
  - goroutine
  - connection pool
  - sql injection
  - rate limit
---

# Security & Performance — Go

## Performance

### Goroutines

```go
// Concurrent operations
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error { return fetchUsers(ctx) })
g.Go(func() error { return fetchOrders(ctx) })
if err := g.Wait(); err != nil { ... }
```

- Goroutines are lightweight (~2KB each).
- Use `errgroup` for coordinated concurrent work.
- Use `context` for cancellation and timeouts.

### Connection Pooling

```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(5 * time.Minute)
```

### Profiling

```go
import _ "net/http/pprof"
// Access: http://localhost:6060/debug/pprof/
```

- `go tool pprof` for CPU/memory profiling.
- `go test -bench=.` for benchmarks.
- `-race` flag for race condition detection.

### JSON Performance

- `encoding/json` for standard use.
- `json-iterator/go` or `sonic` for high-throughput APIs.
- Use struct tags: `json:"name,omitempty"`.

## Security

### SQL Injection Prevention

```go
// ALWAYS parameterized — sqlx handles this
db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)

// NEVER string concatenation
// db.Query("SELECT * FROM users WHERE id = " + id)  // VULNERABLE
```

### Input Validation

```go
type CreateUserDTO struct {
    Name  string `json:"name" validate:"required,min=1,max=100"`
    Email string `json:"email" validate:"required,email"`
}

validate := validator.New()
if err := validate.Struct(dto); err != nil {
    httputil.Error(w, http.StatusBadRequest, err.Error())
    return
}
```

### JWT Authentication

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
        claims, err := validateJWT(token)
        if err != nil {
            httputil.Error(w, http.StatusUnauthorized, "invalid token")
            return
        }
        ctx := context.WithValue(r.Context(), userKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Rate Limiting

```go
limiter := rate.NewLimiter(rate.Every(time.Second), 10) // 10 req/sec

func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "rate limited", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### Secrets

- Environment variables — `os.Getenv("SECRET_KEY")`.
- `godotenv` for local `.env` loading.
- Never log secrets — use `slog` with explicit field selection.
