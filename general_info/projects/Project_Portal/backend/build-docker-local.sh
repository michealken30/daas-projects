#!/bin/bash

echo "Building Go binary locally..."

# Copy templates to local directory for build
cp -r ../frontend/templates ./

# Build the binary locally 
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

echo "Building Docker image with pre-built binary..."
docker build -t backend:latest .

# Clean up
rm -f main
rm -rf ./templates

echo "Docker build complete!"