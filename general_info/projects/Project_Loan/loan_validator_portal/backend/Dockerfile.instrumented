# Multi-stage build for Loan Validator Portal (Instrumented with OpenTelemetry)
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy all backend files including vendor
COPY backend/ ./

# Build the instrumented application using vendored dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o loan-validator-portal main_instrumented.go

# Final stage
FROM alpine:3.18

# Set workdir to /app/backend so ../frontend/templates/* resolves correctly
WORKDIR /app/backend

# Copy the binary from builder
COPY --from=builder /app/loan-validator-portal .

# Copy templates to /app/frontend/templates/ (relative path ../frontend/templates/* from /app/backend)
COPY frontend/templates /app/frontend/templates/

# Expose port
EXPOSE 8080

# Run the instrumented application
CMD ["./loan-validator-portal"]
