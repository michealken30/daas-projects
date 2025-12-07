# Kubernetes Deployment Guide

## 📁 Structure

```
k8s_loan_shared/
├── postgres-configmap.yaml       # Database schemas and init scripts
├── postgres-secret.yaml          # Database credentials
├── postgres-pvc.yaml             # Persistent storage for PostgreSQL
├── postgres-statefulset.yaml    # PostgreSQL deployment
├── postgres-service.yaml         # PostgreSQL service
├── otel-collector-configmap.yaml # OTel Collector configuration
├── otel-collector-deployment.yaml # OTel Collector deployment
└── otel-collector-service.yaml   # OTel Collector service

government_api/k8s/
├── configmap.yaml                # Application configuration
├── secret.yaml                   # Application secrets
├── deployment.yaml               # Deployment with 2 replicas
└── service.yaml                  # ClusterIP service

loan_validator_portal/k8s/
├── configmap.yaml                # Application configuration
├── secret.yaml                   # Application secrets
├── deployment.yaml               # Deployment with 2 replicas
└── service.yaml                  # LoadBalancer service (external access)
```

## 🚀 Deployment Order

### 1. Deploy Infrastructure (PostgreSQL + OTel Collector)

```bash
cd /home/dapo/daas/daascohort3/general_info/projects

# Deploy PostgreSQL
kubectl apply -f k8s/postgres-secret.yaml
kubectl apply -f k8s/postgres-configmap.yaml
kubectl apply -f k8s/postgres-pvc.yaml
kubectl apply -f k8s/postgres-statefulset.yaml
kubectl apply -f k8s/postgres-service.yaml

# Wait for PostgreSQL to be ready
kubectl wait --for=condition=ready pod -l app=postgres --timeout=120s

# Deploy OTel Collector
kubectl apply -f k8s/otel-collector-configmap.yaml
kubectl apply -f k8s/otel-collector-deployment.yaml
kubectl apply -f k8s/otel-collector-service.yaml

# Wait for OTel Collector to be ready
kubectl wait --for=condition=ready pod -l app=otel-collector --timeout=60s
```

### 2. Deploy Government API

```bash
cd government_api

kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Wait for deployment
kubectl wait --for=condition=available deployment/government-api --timeout=120s
```

### 3. Deploy Loan Validator Portal

```bash
cd ../loan_validator_portal

kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Wait for deployment
kubectl wait --for=condition=available deployment/loan-validator-portal --timeout=120s
```

## ✅ Verify Deployment

```bash
# Check all pods
kubectl get pods

# Check services
kubectl get svc

# Get external IP for loan validator portal
kubectl get svc loan-validator-portal-service

# Check logs
kubectl logs -l app=loan-validator-portal --tail=50
kubectl logs -l app=government-loan-bank --tail=50
kubectl logs -l app=otel-collector --tail=50
```

## 🌐 Access Application

```bash
# Get the external IP or use port-forward
kubectl port-forward svc/loan-validator-portal-service 8080:8080

# Then access: http://localhost:8080
```

## 📊 Monitor in Datadog

Traces will automatically flow to: https://us5.datadoghq.com/apm/traces

Filter by:
- Environment: `production`
- Cluster: `daas-k8s-cluster`
- Service: `loan-validator-portal`

## 🔧 Scaling

```bash
# Scale government api
kubectl scale deployment government-api --replicas=3

# Scale loan validator portal
kubectl scale deployment loan-validator-portal --replicas=5
```

## 🗑️ Cleanup

```bash
cd /home/dapo/daas/daascohort3/general_info/projects/Project_Loan

# Delete applications
kubectl delete -f loan_validator_portal/k8s/
kubectl delete -f government_api/k8s/

# Delete infrastructure
kubectl delete -f k8s_loan_shared/otel-collector-service.yaml
kubectl delete -f k8s/otel-collector-deployment.yaml
kubectl delete -f k8s/otel-collector-configmap.yaml
kubectl delete -f k8s/postgres-service.yaml
kubectl delete -f k8s/postgres-statefulset.yaml
kubectl delete -f k8s/postgres-pvc.yaml
kubectl delete -f k8s/postgres-configmap.yaml
kubectl delete -f k8s/postgres-secret.yaml
```

## 🎯 Quick Deploy All

```bash
#!/bin/bash
cd /home/dapo/daas/daascohort3/general_info/projects/Project_Loan

# Infrastructure
kubectl apply -f k8s_loan_shared/

# Wait for infrastructure
kubectl wait --for=condition=ready pod -l app=postgres --timeout=120s
kubectl wait --for=condition=ready pod -l app=otel-collector --timeout=60s

# Applications
kubectl apply -f government_api/k8s/
kubectl apply -f loan_validator_portal/k8s/

# Wait for applications
kubectl wait --for=condition=available deployment/government-api --timeout=120s
kubectl wait --for=condition=available deployment/loan-validator-portal --timeout=120s

echo "✅ Deployment complete!"
kubectl get pods
kubectl get svc
```

## 📝 Notes

- PostgreSQL uses persistent storage (5Gi PVC)
- Both applications have 2 replicas for high availability
- OTel Collector sends traces to Datadog US5
- Loan Validator Portal exposed via LoadBalancer
- Government Bank is ClusterIP (internal only)
- Health checks configured for all services
- Resource limits set for optimal performance
