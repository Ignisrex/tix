# ============================
# 1) Builder image
# ============================

FROM golang:1.25-alpine AS builder

# Install tools needed for go mod (git) and HTTPS (certs)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy the entire core service source into the image
COPY . .

# Download dependencies (go.mod/go.sum must be in or under /app)
RUN go mod download

# Build the binary
# Our main.go is in ./cmd/main.go. We'll build it as /core
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /core ./cmd/main.go

# ============================
# 2) Runtime image
# ============================

FROM alpine:3.20

# Certificates so the app can make HTTPS requests
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy compiled binary from builder
COPY --from=builder /core /app/core

# Optional non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose port (documentation only - actual port comes from env var)
# PORT environment variable should be set via docker-compose or runtime
EXPOSE 8080

# Start the server
ENTRYPOINT ["/app/core"]
