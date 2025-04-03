package middleware

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

// TestPrometheusMiddlewareRegisterMetrics tests the registration of metrics
func TestPrometheusMiddlewareRegisterMetrics(t *testing.T) {
	tests := []struct {
		name         string
		setupMetrics func() (*prometheusMiddleware, *prometheus.Registry)
		expectPanics bool
		panicMessage string
	}{
		{
			name: "Successfully register metrics",
			setupMetrics: func() (*prometheusMiddleware, *prometheus.Registry) {
				reg := prometheus.NewRegistry()
				m := &prometheusMiddleware{
					reg: reg,
					request: prometheus.NewCounterVec(
						prometheus.CounterOpts{
							Name: "http_requests_total",
							Help: "Total number of HTTP requests",
						},
						[]string{"method", "path", "status"},
					),
					latency: prometheus.NewHistogramVec(
						prometheus.HistogramOpts{
							Name: "http_request_duration_seconds",
							Help: "HTTP request latency in seconds",
						},
						[]string{"method", "path"},
					),
				}
				return m, reg
			},
			expectPanics: false,
		},
		{
			name: "Panic on duplicate registration",
			setupMetrics: func() (*prometheusMiddleware, *prometheus.Registry) {
				reg := prometheus.NewRegistry()
				m := &prometheusMiddleware{
					reg: reg,
					request: prometheus.NewCounterVec(
						prometheus.CounterOpts{
							Name: "duplicate_metric",
							Help: "Duplicate metric",
						},
						[]string{"method"},
					),
					latency: prometheus.NewHistogramVec(
						prometheus.HistogramOpts{
							Name: "duplicate_metric",
							Help: "Duplicate metric",
						},
						[]string{"method"},
					),
				}
				return m, reg
			},
			expectPanics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware, _ := tt.setupMetrics()

			if tt.expectPanics {
				assert.Panics(t, func() {
					middleware.registerMetrics()
				})
			} else {
				assert.NotPanics(t, func() {
					middleware.registerMetrics()
				})
			}
		})
	}
}
