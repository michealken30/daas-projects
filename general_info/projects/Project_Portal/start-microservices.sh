#!/bin/bash

echo "Starting microservices..."

# Function to check if port is available
check_port() {
    if lsof -Pi :$1 -sTCP:LISTEN -t >/dev/null ; then
        echo "Port $1 is already in use"
        return 1
    fi
    return 0
}

# Check required ports
echo "Checking ports..."
check_port 8081 || exit 1
check_port 8082 || exit 1  
check_port 8083 || exit 1

echo "All ports available. Starting services..."

# Start auth service in background
echo "Starting auth service on port 8082..."
cd services/auth-service
./auth-main &
AUTH_PID=$!
cd ../..

# Start upload service in background
echo "Starting upload service on port 8083..."
cd services/upload-service  
./upload-main &
UPLOAD_PID=$!
cd ../..

# Wait a moment for services to start
sleep 2

# Start API gateway (main service)
echo "Starting API gateway on port 8081..."
cd services/api-gateway
./gateway-main &
GATEWAY_PID=$!
cd ../..

echo ""
echo "All microservices started!"
echo "Auth Service PID: $AUTH_PID"
echo "Upload Service PID: $UPLOAD_PID"  
echo "API Gateway PID: $GATEWAY_PID"
echo ""
echo "Application available at: http://localhost:8081"
echo ""
echo "To stop all services, run: kill $AUTH_PID $UPLOAD_PID $GATEWAY_PID"

# Wait for any service to exit
wait