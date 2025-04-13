# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy source code
COPY . .

# Download dependencies
RUN go mod download &&\
    # Build the application
    CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Final stage
FROM gcr.io/distroless/static-debian12 AS final

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server .
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Set environment variables
ENV GIN_MODE=release \
    MIGRATIONS_PATH=internal/database/migrations

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./server"]