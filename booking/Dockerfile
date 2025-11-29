# ============================
# 1) Builder image
# ============================

FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /booking ./cmd/main.go

# ============================
# 2) Runtime image
# ============================

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /booking /app/booking

RUN adduser -D -g '' appuser
USER appuser

ENV PORT=8081
EXPOSE 8081

ENTRYPOINT ["/app/booking"]


