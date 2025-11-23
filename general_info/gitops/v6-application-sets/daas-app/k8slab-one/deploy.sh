#!/bin/bash

echo "Deploying microservices to Kubernetes..."

# Apply common resources first
echo "Applying common resources..."
kubectl apply -f k8s/00-common.yaml

# Wait for PVC to be bound
echo "Waiting for PVC to be ready..."
kubectl wait --for=condition=Bound pvc/uploads-pvc --timeout=60s

# Apply services
echo "Deploying auth service..."
kubectl apply -f k8s/01-auth-service.yaml

echo "Deploying upload service..."
kubectl apply -f k8s/02-upload-service.yaml

echo "Deploying API gateway..."
kubectl apply -f k8s/03-api-gateway.yaml

# Wait for deployments to be ready
echo "Waiting for deployments to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/auth-service
kubectl wait --for=condition=available --timeout=300s deployment/upload-service
kubectl wait --for=condition=available --timeout=300s deployment/api-gateway


# Apply network policies
echo "Applying network policies..."
kubectl apply -f k8s/05-network-policies.yaml

echo ""
echo "Deployment complete!"
echo ""
echo "Check status with:"
echo "  kubectl get pods"
echo "  kubectl get services"
echo ""
echo "Access the application:"
echo "  kubectl port-forward service/api-gateway 8081:8081"
echo "  Then visit: http://localhost:8081"
echo ""
echo "Or if using LoadBalancer:"
echo "  kubectl get service api-gateway"