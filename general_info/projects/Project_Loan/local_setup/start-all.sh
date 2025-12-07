#!/bin/bash

# Master startup script - starts everything needed for local development
# Run this from anywhere, it will find the correct paths

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

# Get script directory (resolves symlinks)
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║                                                  ║${NC}"
echo -e "${CYAN}║     Loan Validator System - Local Setup         ║${NC}"
echo -e "${CYAN}║                                                  ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════╝${NC}"
echo ""

# Change to script directory
cd "${SCRIPT_DIR}"

# Step 1: Start Docker services
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Step 1/5: Starting Docker Services${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

if docker ps | grep -q "loan_postgres"; then
    echo -e "${YELLOW}⚠ PostgreSQL already running${NC}"
else
    echo -e "${YELLOW}Starting PostgreSQL...${NC}"
    docker-compose -f docker-compose.local.yml up -d postgres
fi

if docker ps | grep -q "otel-collector"; then
    echo -e "${YELLOW}⚠ OTel Collector already running${NC}"
else
    echo -e "${YELLOW}Starting OpenTelemetry Collector...${NC}"
    docker-compose -f docker-compose.local.yml up -d otel-collector
fi

echo -e "${GREEN}✓ Docker services started${NC}"
echo ""

# Step 2: Wait for services to be ready
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Step 2/5: Waiting for Services${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo -e "${YELLOW}Waiting for PostgreSQL to be healthy...${NC}"
max_attempts=30
attempt=0
until docker exec loan_postgres pg_isready -U postgres > /dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ $attempt -gt $max_attempts ]; then
        echo -e "${RED}✗ PostgreSQL failed to start${NC}"
        exit 1
    fi
    sleep 1
done
echo -e "${GREEN}✓ PostgreSQL is ready${NC}"

echo -e "${YELLOW}Waiting for OTel Collector to start...${NC}"
attempt=0
until docker ps | grep -q "otel-collector"; do
    attempt=$((attempt + 1))
    if [ $attempt -gt $max_attempts ]; then
        echo -e "${RED}✗ OTel Collector failed to start${NC}"
        echo -e "${YELLOW}Checking logs...${NC}"
        docker logs otel-collector 2>&1 | tail -10
        exit 1
    fi
    sleep 1
done

# Give it a moment to initialize
sleep 2

# Check if OTLP port is accessible (more reliable than health endpoint)
if nc -z localhost 4317 2>/dev/null || docker exec otel-collector sh -c "exit 0" 2>/dev/null; then
    echo -e "${GREEN}✓ OTel Collector is ready${NC}"
else
    echo -e "${YELLOW}⚠ OTel Collector is running but ports may not be ready yet${NC}"
    echo -e "${YELLOW}  (This is usually fine - it will be ready when apps start)${NC}"
fi
echo ""

# Step 3: Setup databases
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Step 3/5: Setting Up Databases${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

chmod +x setup-databases.sh
./setup-databases.sh
echo ""

# Step 4: Start Government API
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Step 4/5: Starting Government API${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

chmod +x start-government-api.sh
./start-government-api.sh
echo ""

# Wait a moment for government API to be ready
sleep 2

# Step 5: Start Loan Validator Portal
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Step 5/5: Starting Loan Validator Portal${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

chmod +x start-loan-validator.sh
./start-loan-validator.sh
echo ""

# Final status check
sleep 2

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                                                  ║${NC}"
echo -e "${GREEN}║           🎉 ALL SYSTEMS READY! 🎉               ║${NC}"
echo -e "${GREEN}║                                                  ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════╝${NC}"
echo ""

# Display summary
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}SERVICE SUMMARY${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Check services
services_ok=true

if docker ps | grep -q "loan_postgres"; then
    echo -e "${GREEN}✓${NC} PostgreSQL         (Docker, Port 5432)"
else
    echo -e "${RED}✗${NC} PostgreSQL         (Not Running)"
    services_ok=false
fi

if docker ps | grep -q "otel-collector"; then
    echo -e "${GREEN}✓${NC} OTel Collector     (Docker, Port 4317)"
else
    echo -e "${RED}✗${NC} OTel Collector     (Not Running)"
    services_ok=false
fi

if [ -f "logs/gov-api.pid" ] && ps -p $(cat logs/gov-api.pid) > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Government API     (Go Process, Port 8082)"
else
    echo -e "${RED}✗${NC} Government API     (Not Running)"
    services_ok=false
fi

if [ -f "logs/validator.pid" ] && ps -p $(cat logs/validator.pid) > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Loan Validator     (Go Process, Port 8080)"
else
    echo -e "${RED}✗${NC} Loan Validator     (Not Running)"
    services_ok=false
fi

echo ""
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}ACCESS POINTS${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "  ${YELLOW}Application:${NC}      http://localhost:8080"
echo -e "  ${YELLOW}Government API:${NC}   http://localhost:8082"
echo -e "  ${YELLOW}Datadog APM:${NC}      https://us5.datadoghq.com/apm/traces"
echo ""

echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}LOGS${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "  ${YELLOW}Government API:${NC}   tail -f logs/government-api.log"
echo -e "  ${YELLOW}Loan Validator:${NC}   tail -f logs/loan-validator-portal.log"
echo -e "  ${YELLOW}OTel Collector:${NC}   docker logs -f otel-collector"
echo -e "  ${YELLOW}PostgreSQL:${NC}       docker logs -f loan_postgres"
echo ""

echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}NEXT STEPS${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "  1. ${GREEN}Open browser:${NC} http://localhost:8080"
echo -e "  2. ${GREEN}Register${NC} with any of these names:"
echo -e "     • Oladapo Babalola"
echo -e "     • John Doe"
echo -e "     • Jane Smith"
echo -e "  3. ${GREEN}Login${NC} and click ${YELLOW}\"Validate Your Loan\"${NC}"
echo -e "  4. ${GREEN}View traces${NC} in Datadog APM"
echo ""
echo -e "  ${YELLOW}To stop everything:${NC} ./teardown.sh"
echo ""

if [ "$services_ok" = false ]; then
    echo -e "${RED}⚠ WARNING: Some services failed to start${NC}"
    echo -e "${RED}  Check logs in the logs/ directory${NC}"
    echo -e "${RED}  Run ./teardown.sh and try again${NC}"
    echo ""
    exit 1
fi

echo -e "${GREEN}🚀 Ready to use! Open http://localhost:8080 in your browser${NC}"
echo ""
