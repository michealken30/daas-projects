#!/bin/bash

# Kubernetes Deployment Script for Loan Validator System
# This script deploys all components to a Kubernetes cluster

set -e

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "=================================================="
echo "Deploying Loan Validator System to Kubernetes"
echo "=================================================="
echo ""

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found. Please install kubectl first."
    exit 1
fi

# Check if cluster is accessible
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Cannot connect to Kubernetes cluster. Please check your kubeconfig."
    exit 1
fi

echo "✅ Connected to Kubernetes cluster"
echo ""

# Deploy infrastructure
echo "📦 Deploying Infrastructure (PostgreSQL + OTel Collector)..."
echo ""

kubectl apply -f "$SCRIPT_DIR/postgres-secret.yaml"
kubectl apply -f "$SCRIPT_DIR/postgres-configmap.yaml"
kubectl apply -f "$SCRIPT_DIR/postgres-pvc.yaml"
kubectl apply -f "$SCRIPT_DIR/postgres-statefulset.yaml"
kubectl apply -f "$SCRIPT_DIR/postgres-service.yaml"

echo "⏳ Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres --timeout=180s

echo "✅ PostgreSQL ready"
echo ""

kubectl apply -f "$SCRIPT_DIR/otel-collector-configmap.yaml"
kubectl apply -f "$SCRIPT_DIR/otel-collector-deployment.yaml"
kubectl apply -f "$SCRIPT_DIR/otel-collector-service.yaml"

echo "⏳ Waiting for OTel Collector to be ready..."
kubectl wait --for=condition=ready pod -l app=otel-collector --timeout=60s

echo "✅ OTel Collector ready"
echo ""

# Deploy applications
echo "📦 Deploying Government API..."
kubectl apply -f "$PROJECT_ROOT/government_api/k8s/"

echo "⏳ Waiting for Government API to be ready..."
kubectl wait --for=condition=available deployment/government-api --timeout=120s

echo "✅ Government API ready"
echo ""

echo "📦 Deploying Loan Validator Portal..."
kubectl apply -f "$PROJECT_ROOT/loan_validator_portal/k8s/"

echo "⏳ Waiting for Loan Validator Portal to be ready..."
kubectl wait --for=condition=available deployment/loan-validator-portal --timeout=120s

echo "✅ Loan Validator Portal ready"
echo ""

# Display status
echo "=================================================="
echo "✅ Deployment Complete!"
echo "=================================================="
echo ""

echo "📊 Cluster Status:"
echo ""
kubectl get pods
echo ""
kubectl get svc
echo ""

# Get access information
echo "=================================================="
echo "🌐 Access Information"
echo "=================================================="
echo ""

EXTERNAL_IP=$(kubectl get svc loan-validator-portal-service -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "pending")

if [ "$EXTERNAL_IP" = "pending" ] || [ -z "$EXTERNAL_IP" ]; then
    echo "⏳ External IP is pending. Use port-forward to access:"
    echo ""
    echo "   kubectl port-forward svc/loan-validator-portal-service 8080:8080"
    echo ""
    echo "   Then visit: http://localhost:8080"
else
    echo "🌐 Application URL: http://$EXTERNAL_IP:8080"
fi

echo ""
echo "📊 Datadog Traces: https://us5.datadoghq.com/apm/traces"
echo ""
echo "=================================================="
echo ""

echo "📝 Useful Commands:"
echo ""
echo "  View logs:"
echo "    kubectl logs -l app=loan-validator-portal --tail=50"
echo "    kubectl logs -l app=government-api --tail=50"
echo "    kubectl logs -l app=otel-collector --tail=50"
echo ""
echo "  Scale deployments:"
echo "    kubectl scale deployment loan-validator-portal --replicas=3"
echo "    kubectl scale deployment government-api --replicas=3"
echo ""
echo "  Port forward (if LoadBalancer pending):"
echo "    kubectl port-forward svc/loan-validator-portal-service 8080:8080"
echo ""
echo "=================================================="
