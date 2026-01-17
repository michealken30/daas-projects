# Multi-stage build for Loan Validator Portal (Instrumented with OpenTelemetry)
FROM golang:1.21-alpine AS builder

# Go app lives in backend
WORKDIR /app/backend

# Copy only backend (includes vendor)
COPY backend/ .

# Build the instrumented application
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o loan-validator-portal main_instrumented.go

# Final stage
FROM alpine:3.18

# Runtime working directory
WORKDIR /app/backend

# Copy binary
COPY --from=builder /app/backend/loan-validator-portal .

# Copy templates (path expected by the app)
COPY frontend/templates /app/frontend/templates/

EXPOSE 8080

CMD ["./loan-validator-portal"]
