# --- Stage 1: Tailwind ---
FROM node:24.11.0-alpine AS tailwind-builder
WORKDIR /app
COPY package*.json ./
RUN npm install
# Copying public/ so Tailwind can see your input.css
# Copying internal/ so Tailwind can scan your .templ files for classes
COPY public/ ./public/
COPY internal/ ./internal/
RUN npx @tailwindcss/cli -i ./public/css/input.css -o ./public/css/output.css


# --- Stage 2: Go & Templ ---
FROM golang:1.25.4-alpine AS go-builder
WORKDIR /app
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source (internal, cmd, migrations, etc.)
COPY . .

# Copy the CSS from the previous stage
COPY --from=tailwind-builder /app/public/css/output.css ./public/css/output.css

# Remove stale generated files so templ regenerates them fresh
RUN find . -name "*_templ.go" -delete

# Run templ generate (scans your .templ files in internal/ui)
RUN templ generate

# Build the binary using your specific path
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/app/main.go


# --- Stage 3: Final Image ---
FROM alpine:latest
WORKDIR /app

# Copy the binary
COPY --from=go-builder /app/myapp .
# Copy public assets (images, js, and the built css)
COPY --from=go-builder /app/public ./public
# Copy migrations if your app runs them on startup
COPY --from=go-builder /app/migrations ./migrations

# Optional: Copy the initial DB if it's not managed by a volume
# COPY --from=go-builder /app/data ./data 

EXPOSE 3000
CMD ["./myapp"]