# Multi-stage build for Government API
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o government-api .

# Final stage - use distroless
FROM gcr.io/distroless/static-debian11:nonroot

WORKDIR /app

# Copy binary (distroless already has nonroot user 65532)
COPY --from=builder --chown=65532:65532 /app/government-api ./government-api

EXPOSE 8081

USER 65532:65532

CMD ["/app/government-api"]