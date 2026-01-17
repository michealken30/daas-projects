# Multi-stage build for Government API
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy everything including vendor directory
COPY backend/ ./

# Build the application using vendored dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o government-api .

# Final stage
FROM alpine:3.18

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/government-api .

# Expose port
EXPOSE 8081

# Run the application
CMD ["./government-api"]
