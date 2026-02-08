# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies for CGO (SQLite)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build with CGO enabled for SQLite support
# Use static linking for musl to work with distroless
RUN CGO_ENABLED=1 go build -ldflags="-s -w -linkmode external -extldflags '-static'" -o /api ./cmd/api

# Production stage - use scratch since we have a fully static binary
FROM scratch

# Copy CA certificates for HTTPS calls (Turnstile verification)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /api /api

# Copy migrations
COPY --from=builder /app/migrations /migrations

# Run as non-root (UID 65534 is nobody)
USER 65534:65534

EXPOSE 8080

ENTRYPOINT ["/api"]
