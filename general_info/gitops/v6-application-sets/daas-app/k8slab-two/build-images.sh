#!/bin/bash

echo "Building Docker images for Kubernetes deployment..."

# Registry name
REGISTRY="daasnigeria/daasrepo"

# Build all images
echo "Building auth-service image..."
docker build -t $REGISTRY:auth-service -f services/auth-service/Dockerfile .

echo "Building upload-service image..."
docker build -t $REGISTRY:upload-service -f services/upload-service/Dockerfile .

echo "Building api-gateway image..."
docker build -t $REGISTRY:api-gateway -f services/api-gateway/Dockerfile .

echo ""
echo "Docker images built successfully!"
echo ""
echo "Pushing images to registry..."
docker push $REGISTRY:auth-service
docker push $REGISTRY:upload-service
docker push $REGISTRY:api-gateway

echo ""
echo "List of images:"
docker images | grep "$REGISTRY"
echo ""
echo "All images built and pushed to registry successfully!"
echo "Next steps:"
echo "1. Make sure your Kubernetes cluster is running"
echo "2. Run: chmod +x k8s/deploy.sh"
echo "3. Run: ./k8s/deploy.sh"