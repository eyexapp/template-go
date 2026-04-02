# -- Build stage --
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/myapp ./cmd/api

# -- Run stage --
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/myapp .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["./myapp"]
