# --- Stage 1: Build ---
FROM golang:1.25.4-alpine AS go-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/app/main.go


# --- Stage 2: Final Image ---
FROM alpine:latest
WORKDIR /app

# Install redis-cli for running test scenarios
RUN apk add --no-cache redis

COPY --from=go-builder /app/myapp .
COPY --from=go-builder /app/migrations ./migrations

EXPOSE 3000
CMD ["./myapp"]
