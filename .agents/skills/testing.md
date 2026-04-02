---
name: testing
type: knowledge
version: 1.0.0
agent: CodeActAgent
triggers:
  - test
  - go test
  - table test
  - httptest
  - mock
---

# Testing — Go (testing + httptest)

## Table-Driven Tests

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name  string
        email string
        want  bool
    }{
        {"valid email", "alice@test.com", true},
        {"missing @", "invalid", false},
        {"empty", "", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := validateEmail(tt.email)
            if got != tt.want {
                t.Errorf("validateEmail(%q) = %v, want %v", tt.email, got, tt.want)
            }
        })
    }
}
```

## HTTP Handler Testing

```go
func TestCreateUser(t *testing.T) {
    mockService := &MockUserService{}
    handler := NewHandler(mockService)

    body := `{"name":"Alice","email":"alice@test.com"}`
    req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()

    handler.Create(rec, req)

    if rec.Code != http.StatusCreated {
        t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
    }
}
```

## Service Testing

```go
func TestUserService_Create(t *testing.T) {
    repo := &MockRepo{}
    svc := NewService(repo)

    user, err := svc.Create(context.Background(), CreateDTO{Name: "Alice", Email: "a@b.com"})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "Alice" {
        t.Errorf("name = %q, want Alice", user.Name)
    }
}
```

## Interface Mocking

```go
type MockUserService struct {
    CreateFunc func(ctx context.Context, dto CreateDTO) (*User, error)
}

func (m *MockUserService) Create(ctx context.Context, dto CreateDTO) (*User, error) {
    return m.CreateFunc(ctx, dto)
}
```

## Rules

- Table-driven tests for input variations.
- `httptest.NewRequest` + `httptest.NewRecorder` for handler tests.
- Interface-based mocking — no external mock library needed.
- `go test ./...` — runs all tests recursively.
- `go test -race ./...` — race condition detection.
