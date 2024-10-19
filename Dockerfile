# Build Stage
FROM golang:1.21-alpine AS builder

# Install necessary tools and dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    gcc \
    musl-dev

# Set the working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod tidy

# Copy the project files
COPY . .

# Build the Go application binary
RUN go build -o main ./main.go

# Test Stage
FROM golang:1.21-alpine AS tester

# Set the working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

# Copy the project files
COPY . .

# Run tests
CMD ["go", "test", "./..."]

# Final Stage: Use minimal image
FROM alpine:latest

# Install certificates for secure connections
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Command to run when the container starts
CMD ["./main"]
