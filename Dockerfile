# Stage 1: Build the Go application (godoctor server)
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod .
COPY go.sum .

# Download Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the godoctor server application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o godoctor ./cmd/godoctor

# Stage 2: Create a minimal image
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS requests if needed
RUN apk add --no-cache ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/godoctor .

# Expose the server port
EXPOSE 8080

# Command to run the server application
CMD ["./godoctor"]
