package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	dflBuckets = []float64{0.3, 1.0, 2.5, 5.0}
)

const (
	httpRequestsCount    = "http_requests_total"
	httpRequestsDuration = "http_request_duration_seconds"
)

// Config is the configuration for the middleware
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
}

// NewPrometheusMiddleware creates a new PrometheusMiddleware instance
func NewPrometheusMiddleware(cnfg Config) *prometheusMiddleware {
	var prometheusMiddleware prometheusMiddleware

	prometheusMiddleware.request = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: cnfg.Namespace,
			Subsystem: cnfg.Subsystem,
			Name:      httpRequestsCount,
			Help:      "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
		},
		[]string{"code", "method", "path"},
	)

	if err := prometheus.Register(prometheusMiddleware.request); err != nil {
		log.Println("prometheusMiddleware.request was not registered:", err)
	}

	buckets := dflBuckets
	if len(cnfg.Buckets) == 0 {
		buckets = cnfg.Buckets
	}

	prometheusMiddleware.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: cnfg.Namespace,
		Subsystem: cnfg.Subsystem,
		Name:      httpRequestsDuration,
		Help:      "How long it took to process the request, partitioned by status code, method and HTTP path.",
		Buckets:   buckets,
	},
		[]string{"code", "method", "path"},
	)

	if err := prometheus.Register(prometheusMiddleware.latency); err != nil {
		log.Println("prometheusMiddleware.latency was not registered:", err)
	}

	// Register all the middleware metrics on prometheus registerer.
	//prometheusMiddleware.registerMetrics()

	return &prometheusMiddleware
}

// registerMetrics registers all the metrics on prometheus registerer.
// func (m *prometheusMiddleware) registerMetrics() {
// 	m.reg.MustRegister(
// 		m.request,
// 		m.latency,
// 	)
// }

// Middleware implements the middleware function that will be called by the Gorilla Mux router
// to handle the request.
func (p *prometheusMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()

		delegate := &responseWriterDelegator{ResponseWriter: w}
		rw := delegate

		next.ServeHTTP(rw, r) // call original

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		code := sanitizeCode(delegate.status)
		method := sanitizeMethod(r.Method)

		/// If the status code is 404, and we don't want to use the request path for 404s,
		if p.cfg.DoNotUseRequestPathFor404 && code == "404" {
			path = "404"
		}

		go p.request.WithLabelValues(code, method, path).Inc()

		go p.latency.WithLabelValues(code, method, path).Observe(float64(time.Since(begin)) / float64(time.Second))
	})
}

type responseWriterDelegator struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (r *responseWriterDelegator) WriteHeader(code int) {
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseWriterDelegator) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
}

func sanitizeMethod(m string) string {
	return strings.ToLower(m)
}

func sanitizeCode(s int) string {
	return strconv.Itoa(s)
}
