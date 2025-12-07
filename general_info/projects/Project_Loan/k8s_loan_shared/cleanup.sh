#!/bin/bash

# Kubernetes Cleanup Script
# This script removes all deployed components

set -e

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "=================================================="
echo "Cleaning up Loan Validator System from Kubernetes"
echo "=================================================="
echo ""

read -p "Are you sure you want to delete all resources? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "Cleanup cancelled."
    exit 0
fi

echo ""
echo "🗑️  Deleting applications..."

# Delete loan validator portal
kubectl delete -f "$PROJECT_ROOT/loan_validator_portal/k8s/deployment.yaml" --ignore-not-found=true
kubectl delete -f "$PROJECT_ROOT/loan_validator_portal/k8s/service.yaml" --ignore-not-found=true
kubectl delete -f "$PROJECT_ROOT/loan_validator_portal/k8s/configmap.yaml" --ignore-not-found=true
kubectl delete -f "$PROJECT_ROOT/loan_validator_portal/k8s/secret.yaml" --ignore-not-found=true

# Delete government API
kubectl delete -f "$PROJECT_ROOT/government_api/k8s/deployment.yaml" --ignore-not-found=true
kubectl delete -f "$PROJECT_ROOT/government_api/k8s/service.yaml" --ignore-not-found=true
kubectl delete -f "$PROJECT_ROOT/government_api/k8s/configmap.yaml" --ignore-not-found=true
kubectl delete -f "$PROJECT_ROOT/government_api/k8s/secret.yaml" --ignore-not-found=true

echo "✅ Applications deleted"
echo ""

echo "🗑️  Deleting infrastructure..."

# Delete OTel Collector
kubectl delete -f "$SCRIPT_DIR/otel-collector-service.yaml" --ignore-not-found=true
kubectl delete -f "$SCRIPT_DIR/otel-collector-deployment.yaml" --ignore-not-found=true
kubectl delete -f "$SCRIPT_DIR/otel-collector-configmap.yaml" --ignore-not-found=true

# Delete PostgreSQL
kubectl delete -f "$SCRIPT_DIR/postgres-service.yaml" --ignore-not-found=true
kubectl delete -f "$SCRIPT_DIR/postgres-statefulset.yaml" --ignore-not-found=true
kubectl delete -f "$SCRIPT_DIR/postgres-pvc.yaml" --ignore-not-found=true
kubectl delete -f "$SCRIPT_DIR/postgres-configmap.yaml" --ignore-not-found=true
kubectl delete -f "$SCRIPT_DIR/postgres-secret.yaml" --ignore-not-found=true

# Wait for resources to be deleted
echo "⏳ Waiting for resources to be fully deleted..."
kubectl wait --for=delete pod -l app=loan-validator-portal --timeout=60s 2>/dev/null || true
kubectl wait --for=delete pod -l app=government-api --timeout=60s 2>/dev/null || true
kubectl wait --for=delete pod -l app=postgres --timeout=60s 2>/dev/null || true
kubectl wait --for=delete pod -l app=otel-collector --timeout=60s 2>/dev/null || true

echo "✅ Infrastructure deleted"
echo ""

echo "=================================================="
echo "✅ Cleanup Complete!"
echo "=================================================="
echo ""

kubectl get pods
echo ""
