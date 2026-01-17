#!/bin/bash

# Start Loan Validator Portal in background

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"
LOG_DIR="${SCRIPT_DIR}/logs"

# Create log directory
mkdir -p "${LOG_DIR}"

echo -e "${YELLOW}Starting Loan Validator Portal...${NC}"

# Set environment variables and start
cd "${PROJECT_ROOT}/loan_validator_portal/backend"

export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=loan_validator_db
export PORT=8080
export GOV_BANK_URL=http://localhost:8082
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
export ENV=development

# Download dependencies if needed
if [ ! -d "vendor" ]; then
    echo -e "${YELLOW}  Downloading Go dependencies...${NC}"
    go mod download > /dev/null 2>&1
fi

# Start in background
nohup go run main_instrumented.go > "${LOG_DIR}/loan-validator-portal.log" 2>&1 &
VALIDATOR_PID=$!
echo $VALIDATOR_PID > "${LOG_DIR}/validator.pid"

# Wait for it to start
sleep 3

# Check if it's running
if ps -p $VALIDATOR_PID > /dev/null; then
    echo -e "${GREEN}✓ Loan Validator Portal started (PID: ${VALIDATOR_PID})${NC}"
    echo -e "${GREEN}  Listening on http://localhost:8080${NC}"
    echo -e "${YELLOW}  Logs: ${LOG_DIR}/loan-validator-portal.log${NC}"
else
    echo -e "${RED}✗ Failed to start Loan Validator Portal${NC}"
    echo -e "${RED}  Check logs: ${LOG_DIR}/loan-validator-portal.log${NC}"
    exit 1
fi
