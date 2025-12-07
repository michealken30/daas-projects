#!/bin/bash

# Start Government API in background

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"
LOG_DIR="${SCRIPT_DIR}/logs"

# Create log directory
mkdir -p "${LOG_DIR}"

echo -e "${YELLOW}Starting Government API...${NC}"

# Set environment variables and start
cd "${PROJECT_ROOT}/government_api/backend"

export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=government_loan_db
export PORT=8082

# Download dependencies if needed
if [ ! -d "vendor" ]; then
    echo -e "${YELLOW}  Downloading Go dependencies...${NC}"
    go mod download > /dev/null 2>&1
fi

# Start in background
nohup go run main.go > "${LOG_DIR}/government-api.log" 2>&1 &
GOV_API_PID=$!
echo $GOV_API_PID > "${LOG_DIR}/gov-api.pid"

# Wait for it to start
sleep 3

# Check if it's running
if ps -p $GOV_API_PID > /dev/null; then
    echo -e "${GREEN}✓ Government API started (PID: ${GOV_API_PID})${NC}"
    echo -e "${GREEN}  Listening on http://localhost:8082${NC}"
    echo -e "${YELLOW}  Logs: ${LOG_DIR}/government-api.log${NC}"
else
    echo -e "${RED}✗ Failed to start Government API${NC}"
    echo -e "${RED}  Check logs: ${LOG_DIR}/government-api.log${NC}"
    exit 1
fi
