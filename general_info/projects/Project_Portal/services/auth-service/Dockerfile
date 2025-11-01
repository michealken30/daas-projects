FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY services/shared/ ./shared/

# Copy auth service
COPY services/auth-service/ ./auth-service/
WORKDIR /app/auth-service
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o main .

FROM scratch
COPY --from=builder /app/auth-service/main /main
EXPOSE 8082
ENTRYPOINT ["/main"]