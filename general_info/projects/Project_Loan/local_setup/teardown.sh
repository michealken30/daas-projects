#!/bin/bash

# Teardown script - stops all services and cleans up

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║                                                  ║${NC}"
echo -e "${CYAN}║     Loan Validator System - Teardown            ║${NC}"
echo -e "${CYAN}║                                                  ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════╝${NC}"
echo ""

cd "${SCRIPT_DIR}"

# Stop Go applications
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Stopping Go Applications${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

if [ -f "logs/validator.pid" ]; then
    VALIDATOR_PID=$(cat logs/validator.pid)
    if ps -p $VALIDATOR_PID > /dev/null 2>&1; then
        echo -e "${YELLOW}Stopping Loan Validator Portal (PID: ${VALIDATOR_PID})...${NC}"
        kill $VALIDATOR_PID 2>/dev/null || true
        sleep 1
        # Force kill if still running
        if ps -p $VALIDATOR_PID > /dev/null 2>&1; then
            kill -9 $VALIDATOR_PID 2>/dev/null || true
        fi
        echo -e "${GREEN}✓ Loan Validator Portal stopped${NC}"
    else
        echo -e "${YELLOW}⚠ Loan Validator Portal not running${NC}"
    fi
    rm -f logs/validator.pid
else
    echo -e "${YELLOW}⚠ Loan Validator Portal PID file not found${NC}"
fi

if [ -f "logs/gov-api.pid" ]; then
    GOV_API_PID=$(cat logs/gov-api.pid)
    if ps -p $GOV_API_PID > /dev/null 2>&1; then
        echo -e "${YELLOW}Stopping Government API (PID: ${GOV_API_PID})...${NC}"
        kill $GOV_API_PID 2>/dev/null || true
        sleep 1
        # Force kill if still running
        if ps -p $GOV_API_PID > /dev/null 2>&1; then
            kill -9 $GOV_API_PID 2>/dev/null || true
        fi
        echo -e "${GREEN}✓ Government API stopped${NC}"
    else
        echo -e "${YELLOW}⚠ Government API not running${NC}"
    fi
    rm -f logs/gov-api.pid
else
    echo -e "${YELLOW}⚠ Government API PID file not found${NC}"
fi

# Kill any remaining go run processes for these apps
echo -e "${YELLOW}Checking for any remaining Go processes...${NC}"
pkill -f "go run main.go" 2>/dev/null || true
pkill -f "go run main_instrumented.go" 2>/dev/null || true

# Also kill any processes listening on our ports
lsof -ti:8080 2>/dev/null | xargs kill -9 2>/dev/null || true
lsof -ti:8081 2>/dev/null | xargs kill -9 2>/dev/null || true
lsof -ti:8082 2>/dev/null | xargs kill -9 2>/dev/null || true

echo -e "${GREEN}✓ All Go processes cleaned up${NC}"

echo ""

# Stop Docker services
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Stopping Docker Services${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

if docker ps | grep -q "otel-collector"; then
    echo -e "${YELLOW}Stopping OTel Collector...${NC}"
    docker-compose -f docker-compose.local.yml stop otel-collector
    echo -e "${GREEN}✓ OTel Collector stopped${NC}"
else
    echo -e "${YELLOW}⚠ OTel Collector not running${NC}"
fi

if docker ps | grep -q "loan_postgres"; then
    echo -e "${YELLOW}Stopping PostgreSQL...${NC}"
    docker-compose -f docker-compose.local.yml stop postgres
    echo -e "${GREEN}✓ PostgreSQL stopped${NC}"
else
    echo -e "${YELLOW}⚠ PostgreSQL not running${NC}"
fi

echo ""

# Optionally remove containers
echo -e "${YELLOW}Do you want to remove Docker containers? (y/N):${NC} "
read -r -n 1 REMOVE_CONTAINERS
echo ""

if [[ $REMOVE_CONTAINERS =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Removing Docker containers...${NC}"
    docker-compose -f docker-compose.local.yml down
    echo -e "${GREEN}✓ Docker containers removed${NC}"
else
    echo -e "${YELLOW}⚠ Docker containers stopped but not removed${NC}"
    echo -e "${YELLOW}  To start again without database recreation, just run ./start-all.sh${NC}"
fi

echo ""

# Optionally remove volumes
echo -e "${YELLOW}Do you want to remove database volume? (data will be lost) (y/N):${NC} "
read -r -n 1 REMOVE_VOLUMES
echo ""

if [[ $REMOVE_VOLUMES =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Removing database volume...${NC}"
    docker volume rm loan_postgres_data 2>/dev/null || echo -e "${YELLOW}⚠ Volume already removed or doesn't exist${NC}"
    echo -e "${GREEN}✓ Database volume removed${NC}"
    echo -e "${RED}⚠ All database data has been deleted${NC}"
else
    echo -e "${YELLOW}⚠ Database volume preserved${NC}"
fi

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                                                  ║${NC}"
echo -e "${GREEN}║         ✓ TEARDOWN COMPLETE ✓                   ║${NC}"
echo -e "${GREEN}║                                                  ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════╝${NC}"
echo ""

echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}FINAL STATUS${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Check what's still running
if docker ps | grep -q "loan_postgres\|otel-collector"; then
    echo -e "${YELLOW}⚠ Some Docker containers are still running:${NC}"
    docker ps --filter "name=loan_postgres" --filter "name=otel-collector" --format "  • {{.Names}} ({{.Status}})"
else
    echo -e "${GREEN}✓ All Docker containers stopped${NC}"
fi

if pgrep -f "go run main" > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠ Some Go processes are still running:${NC}"
    ps aux | grep "go run main" | grep -v grep | awk '{print "  • PID " $2}'
else
    echo -e "${GREEN}✓ All Go processes stopped${NC}"
fi

echo ""
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}NEXT STEPS${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "  ${GREEN}To start everything again:${NC}"
echo -e "  ./start-all.sh"
echo ""
echo -e "  ${GREEN}To check what's running:${NC}"
echo -e "  docker ps"
echo -e "  ps aux | grep 'go run'"
echo ""
echo -e "  ${GREEN}To manually remove everything:${NC}"
echo -e "  docker-compose -f docker-compose.local.yml down -v"
echo ""

echo -e "${GREEN}✓ Teardown complete!${NC}"
echo ""
