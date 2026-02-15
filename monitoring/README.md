# Monitoring Setup with Grafana and Prometheus

This directory contains the monitoring configuration for the Go REST API template using Grafana and Prometheus.

## Overview

The monitoring stack consists of:

- **Prometheus**: Time-series database for storing metrics
- **Grafana**: Visualization and dashboard platform
- **Go Application**: Exposes metrics at `/metrics` endpoint

## Architecture

```sh
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Grafana   │◄───│  Prometheus │◄───│ Go App     │
│   :3000     │    │   :9090     │    │ :8080      │
└─────────────┘    └─────────────┘    └─────────────┘
```

## Metrics Collected

The Go application exposes the following Prometheus metrics:

- `http_requests_total`: Total number of HTTP requests (counter)
- `http_request_duration_seconds`: HTTP request duration (histogram)

These metrics are labeled with:

- `code`: HTTP status code
- `method`: HTTP method (GET, POST, etc.)
- `path`: Request path

## Getting Started

### 1. Start the Monitoring Stack

```bash
# Start all services including monitoring
docker-compose up -d

# Or start only monitoring services
docker-compose up -d prometheus grafana
```

### 2. Access the Services

- **Grafana**: http://localhost:3000
  - Username: `admin`
  - Password: `admin`
- **Prometheus**: http://localhost:9090
- **Go Application**: http://localhost:8080
- **Application Metrics**: http://localhost:8080/metrics

### 3. Import Dashboard

The Go REST API dashboard is automatically provisioned and includes:

- Request rate over time
- 95th percentile response time
- Success/error rates by status code
- Requests by HTTP method
- Requests by status code

## Configuration Files

### Prometheus Configuration (`monitoring/prometheus/prometheus.yml`)

- Scrapes metrics from the Go application every 5 seconds
- Stores data for 200 hours
- Includes Redis monitoring (optional)

### Grafana Configuration

- **Datasources**: Automatically configured to connect to Prometheus
- **Dashboards**: Pre-configured dashboard for Go REST API monitoring
- **Provisioning**: Automatic setup of datasources and dashboards

## Customization

### Adding Custom Metrics

To add custom metrics to your Go application:

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    customCounter = promauto.NewCounter(prometheus.CounterOpts{
        Name: "my_custom_counter",
        Help: "Description of my custom counter",
    })
)
```

### Creating Custom Dashboards

1. Access Grafana at http://localhost:3000
2. Create a new dashboard
3. Add panels with Prometheus queries
4. Export the dashboard JSON and place it in `monitoring/grafana/dashboards/`

### Alerting

To set up alerts:

1. Configure alerting rules in Prometheus
2. Set up alert manager
3. Configure notification channels in Grafana

## Useful Prometheus Queries

### Request Rate

```promql
rate(http_requests_total[5m])
```

### Error Rate

```promql
rate(http_requests_total{code=~"5.."}[5m])
```

### 95th Percentile Response Time

```promql
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

### Success Rate

```promql
sum(rate(http_requests_total{code=~"2.."}[5m])) / sum(rate(http_requests_total[5m]))
```

## Troubleshooting

### Prometheus Not Scraping Metrics

1. Check if the Go application is running: `docker-compose ps`
2. Verify metrics endpoint: `curl http://localhost:8080/metrics`
3. Check Prometheus targets: http://localhost:9090/targets

### Grafana Can't Connect to Prometheus

1. Verify Prometheus is running: `docker-compose ps prometheus`
2. Check Prometheus logs: `docker-compose logs prometheus`
3. Verify network connectivity between containers

### Dashboard Not Loading

1. Check if the dashboard JSON is valid
2. Verify the datasource is properly configured
3. Check Grafana logs: `docker-compose logs grafana`

## Production Considerations

### Security

- Change default Grafana credentials
- Use environment variables for sensitive configuration
- Consider using reverse proxy with authentication
- Enable HTTPS for production deployments

### Performance

- Adjust Prometheus retention period based on storage requirements
- Configure appropriate scrape intervals
- Monitor Prometheus and Grafana resource usage
- Consider using Prometheus federation for large deployments

### High Availability

- Use external storage for Prometheus data
- Set up Prometheus federation
- Configure Grafana with external database
- Use load balancers for multiple instances

## Additional Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Prometheus Client Go](https://github.com/prometheus/client_golang)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)
