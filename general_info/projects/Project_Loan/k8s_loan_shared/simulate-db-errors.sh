#!/bin/bash

# Database Stop/Start Script for Error Simulation
# This script stops/starts PostgreSQL to simulate database failures in Datadog

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo ""
echo -e "${BLUE}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                                                  ║${NC}"
echo -e "${BLUE}║     Database Error Simulation                   ║${NC}"
echo -e "${BLUE}║                                                  ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════════╝${NC}"
echo ""

# Check if we're using Kubernetes or local Docker
if kubectl get pods -l app=postgres &> /dev/null 2>&1; then
    echo -e "${GREEN}✓ Using Kubernetes environment${NC}"
    ENVIRONMENT="k8s"
    POD_NAME="postgres-0"
else
    echo -e "${GREEN}✓ Using local Docker environment${NC}"
    ENVIRONMENT="docker"
    CONTAINER_NAME="loan_postgres"
fi

echo ""
echo -e "${YELLOW}What would you like to do?${NC}"
echo ""
echo "  1) Stop PostgreSQL (simulate database outage)"
echo "  2) Start PostgreSQL (restore database)"
echo "  3) Restart PostgreSQL (brief interruption)"
echo "  4) Scale down StatefulSet to 0 (K8s only - complete shutdown)"
echo "  5) Scale up StatefulSet to 1 (K8s only - restore)"
echo "  6) Exit"
echo ""

read -p "Select option (1-6): " option

case $option in
    1)
        echo ""
        echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${RED}Stopping PostgreSQL${NC}"
        echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
        
        if [ "$ENVIRONMENT" = "k8s" ]; then
            echo -e "${YELLOW}Stopping PostgreSQL in pod...${NC}"
            kubectl exec $POD_NAME -- sh -c "pg_ctl -D /var/lib/postgresql/data/pgdata stop -m fast" || true
        else
            echo -e "${YELLOW}Stopping PostgreSQL container...${NC}"
            docker stop $CONTAINER_NAME
        fi
        
        echo -e "${RED}✓ PostgreSQL stopped${NC}"
        echo ""
        echo -e "${YELLOW}⚠ Applications will now encounter database connection errors${NC}"
        echo -e "${YELLOW}⚠ Check Datadog APM for error traces${NC}"
        echo ""
        echo "To restore: Run this script again and select option 2"
        ;;
        
    2)
        echo ""
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${GREEN}Starting PostgreSQL${NC}"
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
        
        if [ "$ENVIRONMENT" = "k8s" ]; then
            echo -e "${YELLOW}Starting PostgreSQL in pod...${NC}"
            kubectl exec $POD_NAME -- sh -c "pg_ctl -D /var/lib/postgresql/data/pgdata start" || \
            kubectl delete pod $POD_NAME --grace-period=0 --force
            
            echo -e "${YELLOW}Waiting for pod to be ready...${NC}"
            kubectl wait --for=condition=ready pod/$POD_NAME --timeout=60s
        else
            echo -e "${YELLOW}Starting PostgreSQL container...${NC}"
            docker start $CONTAINER_NAME
            
            echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
            sleep 5
        fi
        
        echo -e "${GREEN}✓ PostgreSQL started${NC}"
        echo -e "${GREEN}✓ Database is now accessible${NC}"
        ;;
        
    3)
        echo ""
        echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${YELLOW}Restarting PostgreSQL${NC}"
        echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
        
        if [ "$ENVIRONMENT" = "k8s" ]; then
            kubectl delete pod $POD_NAME --grace-period=0 --force
            echo -e "${YELLOW}Waiting for pod to restart...${NC}"
            kubectl wait --for=condition=ready pod/$POD_NAME --timeout=60s
        else
            docker restart $CONTAINER_NAME
            sleep 5
        fi
        
        echo -e "${GREEN}✓ PostgreSQL restarted${NC}"
        ;;
        
    4)
        if [ "$ENVIRONMENT" != "k8s" ]; then
            echo -e "${RED}✗ This option is only available in Kubernetes environment${NC}"
            exit 1
        fi
        
        echo ""
        echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${RED}Scaling Down PostgreSQL StatefulSet${NC}"
        echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
        
        kubectl scale statefulset postgres --replicas=0
        
        echo -e "${YELLOW}Waiting for pod to terminate...${NC}"
        kubectl wait --for=delete pod/$POD_NAME --timeout=60s 2>/dev/null || true
        
        echo -e "${RED}✓ PostgreSQL completely shut down${NC}"
        echo -e "${YELLOW}⚠ All applications will fail with database errors${NC}"
        echo ""
        echo "To restore: Run this script and select option 5"
        ;;
        
    5)
        if [ "$ENVIRONMENT" != "k8s" ]; then
            echo -e "${RED}✗ This option is only available in Kubernetes environment${NC}"
            exit 1
        fi
        
        echo ""
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${GREEN}Scaling Up PostgreSQL StatefulSet${NC}"
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
        
        kubectl scale statefulset postgres --replicas=1
        
        echo -e "${YELLOW}Waiting for pod to be ready...${NC}"
        kubectl wait --for=condition=ready pod/$POD_NAME --timeout=120s
        
        echo -e "${GREEN}✓ PostgreSQL restored${NC}"
        echo -e "${GREEN}✓ Database is now accessible${NC}"
        ;;
        
    6)
        echo "Cancelled."
        exit 0
        ;;
        
    *)
        echo -e "${RED}✗ Invalid option${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Monitoring in Datadog${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo "After the database change, check Datadog APM:"
echo ""
echo "  1. Go to: https://us5.datadoghq.com/apm/traces"
echo "  2. Look for error traces with:"
echo "     • Red error indicators"
echo "     • HTTP 500 status codes"
echo "     • Database connection errors"
echo "     • Error messages in span details"
echo ""
echo "  3. Common errors you'll see:"
echo "     • 'Failed to ping database'"
echo "     • 'connection refused'"
echo "     • 'dial tcp: connect: connection refused'"
echo ""

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}Test the Error${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

if [ "$ENVIRONMENT" = "k8s" ]; then
    GOVT_API_URL=$(kubectl get svc government-api-service -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "<LOAD_BALANCER>")
    echo "Try accessing the API to trigger errors:"
    echo ""
    echo "  curl 'http://${GOVT_API_URL}:8081/api/customer?first_name=John&last_name=Doe'"
    echo ""
    echo "Or use the Loan Validator Portal UI:"
    PORTAL_URL=$(kubectl get svc loan-validator-portal-service -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "<LOAD_BALANCER>")
    echo "  http://${PORTAL_URL}:8080"
else
    echo "Try accessing the API to trigger errors:"
    echo ""
    echo "  curl 'http://localhost:8082/api/customer?first_name=John&last_name=Doe'"
    echo ""
    echo "Or use the Loan Validator Portal:"
    echo "  http://localhost:8080"
fi

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                                                  ║${NC}"
echo -e "${GREEN}║     ✓ Database Control Complete!                ║${NC}"
echo -e "${GREEN}║                                                  ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════╝${NC}"
echo ""
