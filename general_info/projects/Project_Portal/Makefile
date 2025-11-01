.PHONY: build start stop clean docker docker-stop test k8s-build k8s-deploy k8s-cleanup

# Build all microservices
build:
	@echo "Building all microservices..."
	@cd services/shared && go mod tidy
	@cd services/auth-service && go mod tidy && go build -o auth-main .
	@cd services/upload-service && go mod tidy && go build -o upload-main .
	@cd services/api-gateway && go mod tidy && go build -o gateway-main .
	@echo "Build complete!"

# Start all microservices
start: build
	@echo "Starting all microservices..."
	@./start-microservices.sh

# Stop all microservices
stop:
	@echo "Stopping all microservices..."
	@./stop-microservices.sh

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f services/auth-service/auth-main
	@rm -f services/upload-service/upload-main
	@rm -f services/api-gateway/gateway-main
	@echo "Clean complete!"

# Build and run with Docker Compose
docker:
	@echo "Building and starting with Docker Compose..."
	@docker-compose up --build

# Stop Docker containers
docker-stop:
	@echo "Stopping Docker containers..."
	@docker-compose down

# Build Docker images for Kubernetes
k8s-build:
	@echo "Building Docker images for Kubernetes..."
	@./k8s/build-images.sh

# Deploy to Kubernetes
k8s-deploy: k8s-build
	@echo "Deploying to Kubernetes..."
	@./k8s/deploy.sh

# Clean up Kubernetes deployment
k8s-cleanup:
	@echo "Cleaning up Kubernetes deployment..."
	@./k8s/cleanup.sh

# Test endpoints (requires services to be running)
test:
	@echo "Testing microservices endpoints..."
	@echo "Testing auth service..."
	@curl -s http://localhost:8082/register -X POST -H "Content-Type: application/json" -d '{"username":"test","password":"test"}' || echo "Auth service not responding"
	@echo "Testing API gateway..."
	@curl -s http://localhost:8081/ || echo "API gateway not responding"
	@echo "Test complete!"

# Development setup
dev: clean build
	@echo "Starting development environment..."
	@echo "Services will start on:"
	@echo "  - API Gateway: http://localhost:8081"
	@echo "  - Auth Service: http://localhost:8082"  
	@echo "  - Upload Service: http://localhost:8083"
	@./start-microservices.sh