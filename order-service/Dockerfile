# Dockerfile for order-service
FROM golang:1.23 as builder

WORKDIR /app

# Copy go module files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o order-service .

# Create a minimal production image
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/order-service .

# Create an empty .env file if needed (to solve the COPY error)
RUN touch .env

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./order-service"]