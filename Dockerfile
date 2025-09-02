# Build stage
FROM golang:1.21-alpine AS builder

RUN apk update && apk add --no-cache ca-certificates git
RUN update-ca-certificates

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o flyctl .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/flyctl .

# Create necessary directories
RUN mkdir -p /tmp /opt

# Make the binary executable
RUN chmod +x ./flyctl

EXPOSE 8080

CMD ["./flyctl"]