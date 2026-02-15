#!/bin/bash

# Test script for monitoring setup
# This script verifies that Prometheus and Grafana are working correctly

set -e

echo "🔍 Testing Monitoring Setup..."
echo "================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if a service is responding
check_service() {
    local name=$1
    local url=$2
    local timeout=${3:-10}
    
    echo -n "Checking $name... "
    
    if curl -s --max-time $timeout "$url" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ OK${NC}"
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        return 1
    fi
}

# Function to check if metrics endpoint is accessible
check_metrics() {
    echo -n "Checking application metrics endpoint... "
    
    if curl -s --max-time 10 "http://localhost:8080/metrics" | grep -q "http_requests_total"; then
        echo -e "${GREEN}✓ OK${NC}"
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        return 1
    fi
}

# Function to check if Prometheus is scraping metrics
check_prometheus_targets() {
    echo -n "Checking Prometheus targets... "
    
    if curl -s --max-time 10 "http://localhost:9090/api/v1/targets" | grep -q "UP"; then
        echo -e "${GREEN}✓ OK${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠ WARNING - No UP targets found${NC}"
        return 1
    fi
}

# Function to check if Grafana datasource is configured
check_grafana_datasource() {
    echo -n "Checking Grafana datasource... "
    
    # This is a basic check - in a real scenario you'd need authentication
    if curl -s --max-time 10 "http://localhost:3000/api/health" | grep -q "ok"; then
        echo -e "${GREEN}✓ OK${NC}"
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        return 1
    fi
}

# Main test execution
echo "Starting monitoring tests..."
echo ""

# Check if services are running
echo "1. Service Health Checks:"
check_service "Go Application" "http://localhost:8080/health"
check_service "Prometheus" "http://localhost:9090"
check_service "Grafana" "http://localhost:3000"
check_service "Jaeger" "http://localhost:16686"

echo ""
echo "2. Metrics and Monitoring:"
check_metrics
check_prometheus_targets
check_grafana_datasource

echo ""
echo "3. Quick Metrics Sample:"
echo "Fetching current metrics from application..."
curl -s "http://localhost:8080/metrics" | grep -E "(http_requests_total|http_request_duration_seconds)" | head -5

echo ""
echo "4. Access URLs:"
echo -e "${YELLOW}Grafana Dashboard:${NC} http://localhost:3000 (admin/admin)"
echo -e "${YELLOW}Prometheus:${NC} http://localhost:9090"
echo -e "${YELLOW}Jaeger Tracing:${NC} http://localhost:16686"
echo -e "${YELLOW}Application Metrics:${NC} http://localhost:8080/metrics"
echo -e "${YELLOW}API Documentation:${NC} http://localhost:8080/swagger/"

echo ""
echo "5. Useful Commands:"
echo -e "${YELLOW}View all containers:${NC} docker-compose ps"
echo -e "${YELLOW}View logs:${NC} docker-compose logs -f [service_name]"
echo -e "${YELLOW}Restart services:${NC} docker-compose restart"
echo -e "${YELLOW}Stop all:${NC} docker-compose down"

echo ""
echo "✅ Monitoring setup test completed!"
echo ""
echo "📊 Next Steps:"
echo "1. Open Grafana at http://localhost:3000"
echo "2. Login with admin/admin"
echo "3. The Go REST API dashboard should be automatically loaded"
echo "4. Generate some traffic to see metrics: curl http://localhost:8080/api/v1/health" 