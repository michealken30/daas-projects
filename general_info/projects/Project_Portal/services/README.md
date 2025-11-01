# Microservices Architecture

This directory contains the microservices version of the Image Upload application.

## Architecture

The application has been split into the following services:

1. **Auth Service** (Port 8082)
   - Handles user registration and login
   - Generates JWT tokens
   - Manages user authentication

2. **Upload Service** (Port 8083)
   - Handles file uploads
   - Validates JWT tokens
   - Serves uploaded static files
   - Provides user profile information

3. **API Gateway** (Port 8081)
   - Routes requests to appropriate microservices
   - Serves the frontend HTML
   - Acts as a single entry point for clients

4. **Shared Module**
   - Common types and utilities
   - JWT validation and generation
   - Password hashing functions

## Running the Application

### Option 1: Docker Compose (Recommended)
```bash
docker-compose up --build
```

### Option 2: Local Development
1. Build all services:
```bash
./build-microservices.sh
```

2. Start all services:
```bash
./start-microservices.sh
```

### Option 3: Manual Start
1. Start each service in separate terminals:
```bash
# Terminal 1 - Auth Service
cd auth-service && go run main.go

# Terminal 2 - Upload Service  
cd upload-service && go run main.go

# Terminal 3 - API Gateway
cd api-gateway && go run main.go
```

## Service Endpoints

### Auth Service (8082)
- POST /register - Register new user (JSON)
- POST /login - User login (JSON)
- POST /register-form - Register new user (Form)
- POST /login-form - User login (Form)

### Upload Service (8083)
- POST /upload - Upload file (requires auth)
- GET /profile - Get user profile (requires auth)
- GET /uploads/* - Serve static files

### API Gateway (8081)
- GET / - Serve frontend HTML
- POST /auth/login - Proxy to auth service
- POST /auth/register - Proxy to auth service
- POST /api/* - Proxy to appropriate services
- GET /uploads/* - Proxy to upload service

## Benefits of Microservices Architecture

1. **Separation of Concerns**: Each service has a single responsibility
2. **Independent Scaling**: Services can be scaled independently
3. **Technology Diversity**: Each service can use different technologies
4. **Fault Isolation**: Failure in one service doesn't bring down the entire application
5. **Independent Deployment**: Services can be deployed independently

## Development Notes

- All services use the shared module for common functionality
- JWT tokens are used for authentication between services
- The API Gateway handles request routing and load balancing
- Services communicate via HTTP REST APIs