#!/bin/bash

echo "Cleaning up Kubernetes deployment..."

# Delete in reverse order
echo "Removing network policies..."
kubectl delete -f k8s/05-network-policies.yaml --ignore-not-found=true

echo "Removing autoscaling..."
kubectl delete -f k8s/04-autoscaling.yaml --ignore-not-found=true

echo "Removing API gateway..."
kubectl delete -f k8s/03-api-gateway.yaml --ignore-not-found=true

echo "Removing upload service..."
kubectl delete -f k8s/02-upload-service.yaml --ignore-not-found=true

echo "Removing auth service..."
kubectl delete -f k8s/01-auth-service.yaml --ignore-not-found=true

echo "Removing common resources..."
kubectl delete -f k8s/00-common.yaml --ignore-not-found=true

echo ""
echo "Cleanup complete!"
echo ""
echo "Verify with:"
echo "  kubectl get pods"
echo "  kubectl get services"
echo "  kubectl get pvc"