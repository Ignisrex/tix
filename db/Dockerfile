# db/Dockerfile
FROM golang:1.25-alpine

# Install deps (optional but nice to have)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Install goose (v3 CLI uses env vars for driver/DB string)
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy migrations into the image
COPY migrations ./migrations

# These will be provided by docker-compose, but we set sensible defaults
ENV GOOSE_DRIVER=postgres
ENV GOOSE_DBSTRING=postgres://postgres:password@db:5432/tix_db?sslmode=disable

# Default entrypoint: run all migrations "up"
# You can override command in docker-compose if needed (e.g. `down`, `redo`, etc.)
ENTRYPOINT [ "sh", "-c", "goose -dir ./migrations up" ]