FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY services/shared/ ./shared/

# Copy upload service
COPY services/upload-service/ ./upload-service/
WORKDIR /app/upload-service
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o main .

FROM scratch
COPY --from=builder /app/upload-service/main /main
EXPOSE 8083
ENTRYPOINT ["/main"]