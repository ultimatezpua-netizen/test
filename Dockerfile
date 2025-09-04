# Multi-stage build for Go application
FROM golang:1.21-alpine AS builder

# Install ca-certificates and git
RUN apk update && apk add --no-cache ca-certificates git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY *.go ./

# Build the application with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o flyctl .

# Final stage - minimal Alpine image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/flyctl .

# Create test JSON files for demonstration
RUN mkdir -p /tmp /opt && \
    echo '{"id":"session-12345","status":"active","config":{"name":"my-flyapp","version":"v1.0.2"}}' > /tmp/session.json && \
    echo '{"version":"v1.0.2","application":"my-flyapp"}' > /tmp/manifest.json && \
    echo '{"theme":"dark","language":"en-US"}' > /opt/customize.json

# Health check using the flyctl command
HEALTHCHECK CMD ./flyctl launch sessions finalize --session-path /tmp/session.json --manifest-path /tmp/manifest.json --from-file /opt/customize.json || exit 1

# Default command
CMD ["./flyctl", "--help"]