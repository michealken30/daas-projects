FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy templates first
COPY frontend/templates/ ./templates/

# Copy api gateway
COPY services/api-gateway/ ./api-gateway/
WORKDIR /app/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o main .

FROM scratch
COPY --from=builder /app/api-gateway/main /main
COPY --from=builder /app/templates/ /templates/
EXPOSE 8081
ENTRYPOINT ["/main"]