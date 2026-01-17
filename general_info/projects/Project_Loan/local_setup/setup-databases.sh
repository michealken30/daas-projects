#!/bin/bash

# Database setup script for local development
# Creates databases and runs schema migrations

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}=================================${NC}"
echo -e "${BLUE}Database Setup${NC}"
echo -e "${BLUE}=================================${NC}"
echo ""

CONTAINER_NAME="loan_postgres"
DB_USER="postgres"
VALIDATOR_DB="loan_validator_db"
GOVERNMENT_DB="government_loan_db"

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

# Wait for PostgreSQL to be ready
echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
max_attempts=30
attempt=0
until docker exec ${CONTAINER_NAME} pg_isready -U ${DB_USER} > /dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ $attempt -gt $max_attempts ]; then
        echo -e "${RED}✗ PostgreSQL failed to start${NC}"
        exit 1
    fi
    echo -e "${YELLOW}  Attempt $attempt/$max_attempts...${NC}"
    sleep 1
done
echo -e "${GREEN}✓ PostgreSQL is ready${NC}"
echo ""

# Create databases
echo -e "${BLUE}Creating databases...${NC}"

docker exec ${CONTAINER_NAME} psql -U ${DB_USER} -tc "SELECT 1 FROM pg_database WHERE datname = '${VALIDATOR_DB}'" | grep -q 1 || \
    docker exec ${CONTAINER_NAME} psql -U ${DB_USER} -c "CREATE DATABASE ${VALIDATOR_DB};"
echo -e "${GREEN}✓ Database ${VALIDATOR_DB} ready${NC}"

docker exec ${CONTAINER_NAME} psql -U ${DB_USER} -tc "SELECT 1 FROM pg_database WHERE datname = '${GOVERNMENT_DB}'" | grep -q 1 || \
    docker exec ${CONTAINER_NAME} psql -U ${DB_USER} -c "CREATE DATABASE ${GOVERNMENT_DB};"
echo -e "${GREEN}✓ Database ${GOVERNMENT_DB} ready${NC}"

echo ""
echo -e "${BLUE}Running schema migrations...${NC}"

# Run Loan Validator Portal schema
echo -e "${YELLOW}Setting up ${VALIDATOR_DB} schema...${NC}"
docker exec -i ${CONTAINER_NAME} psql -U ${DB_USER} -d ${VALIDATOR_DB} < "${PROJECT_ROOT}/loan_validator_portal/database/schema.sql"
echo -e "${GREEN}✓ Schema for ${VALIDATOR_DB} applied${NC}"

# Run Government API schema
echo -e "${YELLOW}Setting up ${GOVERNMENT_DB} schema...${NC}"
docker exec -i ${CONTAINER_NAME} psql -U ${DB_USER} -d ${GOVERNMENT_DB} < "${PROJECT_ROOT}/government_api/database/schema.sql"
echo -e "${GREEN}✓ Schema for ${GOVERNMENT_DB} applied (with sample data)${NC}"

echo ""
echo -e "${GREEN}=================================${NC}"
echo -e "${GREEN}Database setup complete!${NC}"
echo -e "${GREEN}=================================${NC}"
