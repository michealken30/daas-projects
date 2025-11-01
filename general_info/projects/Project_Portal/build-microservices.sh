#!/bin/bash

echo "Building microservices locally..."

# Build shared module
echo "Building shared module..."
cd services/shared
go mod tidy
cd ../..

# Build auth service
echo "Building auth service..."
cd services/auth-service
go mod tidy
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-main .
cd ../..

# Build upload service  
echo "Building upload service..."
cd services/upload-service
go mod tidy
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o upload-main .
cd ../..

# Build API gateway
echo "Building API gateway..."
cd services/api-gateway
go mod tidy
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway-main .
cd ../..

echo "Local build complete!"
echo ""
echo "To run the services:"
echo "1. Start auth service: cd services/auth-service && ./auth-main"
echo "2. Start upload service: cd services/upload-service && ./upload-main" 
echo "3. Start API gateway: cd services/api-gateway && ./gateway-main"
echo ""
echo "Or use Docker Compose: docker-compose up --build"