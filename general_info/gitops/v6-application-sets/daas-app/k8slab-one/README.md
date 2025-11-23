# Kubernetes Deployment Guide

This directory contains Kubernetes manifests to deploy the microservices application to a Kubernetes cluster.

## Prerequisites

- Kubernetes cluster (minikube, kind, or cloud provider)
- kubectl configured to connect to your cluster
- Docker for building images

## Quick Start

### 1. Build Docker Images
```bash
chmod +x k8s/build-images.sh
./k8s/build-images.sh
```

### 2. Deploy to Kubernetes
```bash
chmod +x k8s/deploy.sh
./k8s/deploy.sh
```

### 3. Access the Application
```bash
# Port forward to access locally
kubectl port-forward service/api-gateway 8081:8081

# Then visit: http://localhost:8081
```

### 4. Clean Up
```bash
chmod +x k8s/cleanup.sh
./k8s/cleanup.sh
```

## Manifest Files

### `00-common.yaml`
- **ConfigMap**: Shared configuration for all services
- **PersistentVolume**: Storage for uploaded files
- **PersistentVolumeClaim**: Storage claim for upload service

### `01-auth-service.yaml`
- **Deployment**: Auth service with 2 replicas
- **Service**: Internal ClusterIP service on port 8082
- **Health Checks**: Liveness and readiness probes

### `02-upload-service.yaml`
- **Deployment**: Upload service with 2 replicas and persistent storage
- **Service**: Internal ClusterIP service on port 8083
- **Volume Mount**: Persistent storage for uploads

### `03-api-gateway.yaml`
- **Deployment**: API Gateway with 2 replicas
- **Service**: LoadBalancer service for external access on port 8081
- **Ingress**: Optional ingress configuration


### `04-network-policies.yaml`
- **NetworkPolicy**: Security policies restricting inter-service communication
- Auth/Upload services only accept traffic from API gateway
- API gateway can communicate with all services

## Architecture in Kubernetes

```
┌─────────────────┐
│   LoadBalancer  │ ← External Traffic
│   (api-gateway) │
└─────────┬───────┘
          │
┌─────────▼───────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │    │   Auth Service  │    │ Upload Service  │
│   Deployment    │───▶│   Deployment    │    │   Deployment    │
│   (2-5 pods)    │    │   (2-10 pods)   │    │   (2-10 pods)   │
└─────────────────┘    └─────────────────┘    └─────────┬───────┘
                                                         │
                                              ┌─────────▼───────┐
                                              │ PersistentVolume│
                                              │   (uploads)     │
                                              └─────────────────┘
```

## Features

- **High Availability**: Multiple replicas for each service
- **Auto-Scaling**: Automatic scaling based on resource usage
- **Persistent Storage**: Uploaded files survive pod restarts
- **Security**: Network policies restrict communication
- **Health Checks**: Kubernetes monitors service health
- **Load Balancing**: Traffic distributed across replicas

## Monitoring

```bash
# Check pod status
kubectl get pods

# Check services
kubectl get services

# Check autoscaling status
kubectl get hpa

# View logs
kubectl logs -f deployment/api-gateway
kubectl logs -f deployment/auth-service
kubectl logs -f deployment/upload-service

# Check resource usage
kubectl top pods
```

## Troubleshooting

### Pods not starting
```bash
kubectl describe pod <pod-name>
kubectl logs <pod-name>
```

### Service not accessible
```bash
kubectl get endpoints
kubectl describe service api-gateway
```

### Storage issues
```bash
kubectl get pv
kubectl get pvc
kubectl describe pvc uploads-pvc
```

### Network issues
```bash
kubectl get networkpolicies
kubectl describe networkpolicy <policy-name>
```