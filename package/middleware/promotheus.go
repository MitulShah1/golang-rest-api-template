// Package middleware provides HTTP middleware components for the application.
// It includes authentication, CORS, logging, and telemetry middleware.
package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var dflBuckets = []float64{0.3, 1.0, 2.5, 5.0}

const (
	httpRequestsCount    = "http_requests_total"
	httpRequestsDuration = "http_request_duration_seconds"
	defaultSubsystem     = "golang_rest_api_template"
)

// Config specifies options how to create new PrometheusMiddleware.
type Config struct {
	// Namespace is components of the fully-qualified name of the Metric (created by joining Namespace,Subsystem and Name components with "_")
	// Optional
	Namespace string

	// Subsystem is components of the fully-qualified name of the Metric (created by joining Namespace,Subsystem and Name components with "_")
	// Defaults to: "golang_rest_api_template"
	Subsystem string

	// Buckets specifies an custom buckets to be used in request histograpm.
	Buckets []float64

	// If DoNotUseRequestPathFor404 is true, all 404 responses (due to non-matching route) will have the same `url` label and
	// thus won't generate new metrics.
	DoNotUseRequestPathFor404 bool
}

// PrometheusMiddleware specifies the metrics that is going to be generated
type prometheusMiddleware struct {
	request *prometheus.CounterVec
	latency *prometheus.HistogramVec

	// cfg is the configuration
	cfg Config

	// reg is a prometheus registry
	reg prometheus.Registerer
}

// NewPrometheusMiddleware creates a new PrometheusMiddleware instance
func NewPrometheusMiddleware(cnfg Config) *prometheusMiddleware {
	// Set default subsystem if not provided
	if cnfg.Subsystem == "" {
		cnfg.Subsystem = defaultSubsystem
	}

	// Fix bucket assignment logic
	buckets := dflBuckets
	if len(cnfg.Buckets) > 0 {
		buckets = cnfg.Buckets
	}

	m := &prometheusMiddleware{
		reg: prometheus.DefaultRegisterer,
		cfg: cnfg,
		request: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cnfg.Namespace,
				Subsystem: cnfg.Subsystem,
				Name:      httpRequestsCount,
				Help:      "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			},
			[]string{"code", "method", "path"},
		),
		latency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cnfg.Namespace,
			Subsystem: cnfg.Subsystem,
			Name:      httpRequestsDuration,
			Help:      "How long it took to process the request, partitioned by status code, method and HTTP path.",
			Buckets:   buckets,
		},
			[]string{"code", "method", "path"},
		),
	}

	// Register all the middleware metrics on prometheus registerer.
	m.registerMetrics()

	return m
}

// Middleware implements the middleware function that will be called by the Gorilla Mux router
// to handle the request.
func (m *prometheusMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()

		delegate := &responseWriterDelegator{
			ResponseWriter: w,
			status:         http.StatusOK, // Initialize with default status
		}

		next.ServeHTTP(delegate, r) // call original

		route := mux.CurrentRoute(r)

		var path string
		if route != nil {
			path, _ = route.GetPathTemplate()
		} else {
			// If no route was matched, it's a 404
			path = r.URL.Path
			if m.cfg.DoNotUseRequestPathFor404 {
				path = "404"
			}
		}

		code := sanitizeCode(delegate.status)
		method := sanitizeMethod(r.Method)

		// Execute counter and histogram operations directly instead of in goroutines
		m.request.WithLabelValues(code, method, path).Inc()
		m.latency.WithLabelValues(code, method, path).Observe(float64(time.Since(begin)) / float64(time.Second))
	})
}

// registerMetrics registers all the metrics on prometheus registerer.
func (m *prometheusMiddleware) registerMetrics() {
	m.reg.MustRegister(
		m.request,
		m.latency,
	)
}

func sanitizeMethod(m string) string {
	return strings.ToLower(m)
}

func sanitizeCode(s int) string {
	return strconv.Itoa(s)
}
