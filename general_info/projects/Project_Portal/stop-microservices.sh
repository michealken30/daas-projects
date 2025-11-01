#!/bin/bash

echo "Stopping all microservices..."

# Kill all Go processes running our services
pkill -f "auth-main"
pkill -f "upload-main" 
pkill -f "gateway-main"

echo "All microservices stopped!"