# First stage: Build the Go application
FROM golang:1.21-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Install dependencies
RUN go mod tidy

# Copy the Go source code from the auth_service directory
COPY . .

# Build the Go binary for the service
RUN go build -o auth_service .

# Second stage: Create a minimal image for running the service
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder image
COPY --from=builder /app/auth_service .
