# Image Upload Application - Microservices Architecture

A modern, scalable image upload application built with Go microservices, featuring user authentication, file uploads, and a responsive web interface.

## 🏗️ Architecture Overview

This application follows a microservices architecture pattern:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │   Auth Service  │
│   (HTML/JS)     │───▶│   Port: 8081    │───▶│   Port: 8082    │
└─────────────────┘    │                 │    └─────────────────┘
                       │                 │    ┌─────────────────┐
                       │                 │───▶│ Upload Service  │
                       └─────────────────┘    │   Port: 8083    │
                                              └─────────────────┘
```

### Services

- **🚪 API Gateway** (Port 8081): Routes requests, serves frontend, handles CORS
- **🔐 Auth Service** (Port 8082): User registration, login, JWT token management
- **📁 Upload Service** (Port 8083): File uploads, static file serving, user profiles
- **📚 Shared Module**: Common utilities, types, and JWT validation

## 📂 Project Structure

```
Project_Portal/
├── frontend/                   # Frontend templates and assets
│   └── templates/
│       └── index.html
├── services/                   # Microservices
│   ├── api-gateway/           # Request routing and frontend serving
│   ├── auth-service/          # Authentication service
│   ├── upload-service/        # File upload service
│   ├── shared/                # Common utilities and types
│   └── README.md
├── k8s/                       # Kubernetes deployment manifests
│   ├── 00-common.yaml         # ConfigMaps, PVs, PVCs
│   ├── 01-auth-service.yaml   # Auth service deployment
│   ├── 02-upload-service.yaml # Upload service deployment
│   ├── 03-api-gateway.yaml    # API gateway deployment
│   ├── 04-autoscaling.yaml    # Horizontal Pod Autoscalers
│   ├── 05-network-policies.yaml # Network security policies
│   ├── build-images.sh        # Build Docker images
│   ├── deploy.sh              # Deploy to Kubernetes
│   ├── cleanup.sh             # Clean up deployment
│   └── README.md
├── docker-compose.yml         # Docker Compose configuration
├── Makefile                   # Build automation
└── *.sh                       # Management scripts
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Make
- Kubernetes cluster (for K8s deployment)

### Local Development

1. **Clone and navigate to the project**
```bash
cd /path/to/Project_Portal
```

2. **Start all services**
```bash
make dev
```

3. **Access the application**
- Frontend: http://localhost:8081
- Auth Service: http://localhost:8082
- Upload Service: http://localhost:8083

### Docker Deployment

```bash
# Start with Docker Compose
make docker

# Stop services
make docker-stop
```

### Kubernetes Deployment

```bash
# Build and deploy to Kubernetes
make k8s-deploy

# Access via port-forward
kubectl port-forward service/api-gateway 8081:8081

# Clean up
make k8s-cleanup
```

## 📋 Available Commands

### Development
```bash
make build          # Build all services
make start           # Start all services
make stop            # Stop all services
make dev             # Clean, build, and start
make clean           # Clean build artifacts
make test            # Test endpoints
```

### Docker
```bash
make docker          # Build and start with Docker Compose
make docker-stop     # Stop Docker containers
```

### Kubernetes
```bash
make k8s-build       # Build Docker images for K8s
make k8s-deploy      # Deploy to Kubernetes
make k8s-cleanup     # Clean up K8s deployment
```

## 🌟 Features

### Application Features
- ✅ **User Authentication**: Registration and login with JWT tokens
- ✅ **File Upload**: Secure image upload with unique naming
- ✅ **Responsive UI**: Modern web interface with HTMX
- ✅ **Static File Serving**: Efficient image serving
- ✅ **Session Management**: JWT-based authentication

### Architecture Features
- ✅ **Microservices**: Decoupled, independently scalable services
- ✅ **API Gateway**: Single entry point with request routing
- ✅ **Health Checks**: Service monitoring and auto-recovery
- ✅ **Load Balancing**: Distributed traffic handling
- ✅ **Auto Scaling**: Kubernetes HPA for dynamic scaling
- ✅ **Network Security**: Kubernetes network policies
- ✅ **Persistent Storage**: Volume-mounted file persistence
- ✅ **Container Ready**: Docker and Kubernetes deployments

## 🔧 Configuration

### Environment Variables
- `GIN_MODE`: Set to "release" for production
- `JWT_SECRET`: Secret key for JWT token signing
- `AUTH_SERVICE_URL`: Auth service endpoint
- `UPLOAD_SERVICE_URL`: Upload service endpoint

### Service Ports
- **API Gateway**: 8081 (external access)
- **Auth Service**: 8082 (internal)
- **Upload Service**: 8083 (internal)

## 📊 Monitoring and Observability

### Health Endpoints
- API Gateway: `GET /`
- Auth Service: `POST /register` (dummy endpoint)
- Upload Service: `GET /uploads/`

### Kubernetes Monitoring
```bash
# Check pod status
kubectl get pods

# View service logs
kubectl logs -f deployment/api-gateway

# Monitor resource usage
kubectl top pods

# Check autoscaling
kubectl get hpa
```

## 🔒 Security

### Network Policies
- Auth and Upload services only accept traffic from API Gateway
- API Gateway can communicate with all backend services
- External traffic only reaches API Gateway

### Authentication
- JWT-based authentication with configurable secret
- Secure password hashing with SHA256
- Token validation middleware

## 🎯 Benefits of This Architecture

1. **Scalability**: Each service can scale independently
2. **Reliability**: Service isolation prevents cascading failures
3. **Maintainability**: Clear separation of concerns
4. **Flexibility**: Services can be updated independently
5. **Performance**: Optimized request routing and caching
6. **Security**: Network policies and authentication isolation
7. **DevOps Ready**: Container-native with K8s support

## 📚 Documentation

- [Services Documentation](services/README.md)
- [Kubernetes Deployment Guide](k8s/README.md)

## 🤝 Development Workflow

1. **Make changes** to individual services
2. **Build and test** locally with `make dev`
3. **Test with Docker** using `make docker`
4. **Deploy to K8s** with `make k8s-deploy`
5. **Monitor and scale** as needed

This architecture provides a solid foundation for a production-ready, scalable web application! 🚀