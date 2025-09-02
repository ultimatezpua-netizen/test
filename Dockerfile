# Build stage
FROM golang:1.21-bullseye AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o flyctl .

# Final stage
FROM debian:bullseye-slim

# Install ca-certificates for HTTPS
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/flyctl .

# Create non-root user
RUN useradd -m appuser
RUN chown -R appuser:appuser /app

USER appuser

# Expose port (required for Fly.io)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ./flyctl --help > /dev/null || exit 1

# For Fly.io deployment, start the HTTP server
CMD ["./flyctl", "server"]