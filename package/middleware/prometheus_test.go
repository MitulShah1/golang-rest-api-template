package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPrometheusMiddleware(t *testing.T) {

	t.Skip("Skipping testNewPrometheusMiddleware")

	tests := []struct {
		name           string
		config         Config
		expectedConfig Config
	}{
		{
			name:   "Default configuration",
			config: Config{},
			expectedConfig: Config{
				Subsystem: defaultSubsystem,
			},
		},
		{
			name: "Custom namespace and subsystem",
			config: Config{
				Namespace: "test",
				Subsystem: "api",
			},
			expectedConfig: Config{
				Namespace: "test",
				Subsystem: "api",
			},
		},
		{
			name: "Custom buckets",
			config: Config{
				Buckets: []float64{1.0, 2.0, 3.0},
			},
			expectedConfig: Config{
				Subsystem: defaultSubsystem,
				Buckets:   []float64{1.0, 2.0, 3.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new registry for testing
			reg := prometheus.NewRegistry()

			// Store the default registerer
			defaultReg := prometheus.DefaultRegisterer
			// Replace default registerer with our test registry
			prometheus.DefaultRegisterer = reg
			// Restore the default registerer after the test
			defer func() { prometheus.DefaultRegisterer = defaultReg }()

			middleware := NewPrometheusMiddleware(tt.config)

			assert.Equal(t, tt.expectedConfig.Namespace, middleware.cfg.Namespace)
			assert.Equal(t, tt.expectedConfig.Subsystem, middleware.cfg.Subsystem)

			// Check if buckets are set correctly
			if len(tt.config.Buckets) > 0 {
				assert.Equal(t, tt.config.Buckets, middleware.cfg.Buckets)
			}

			// Check if metrics are registered
			metrics, err := reg.Gather()
			require.NoError(t, err)

			// There should be two metrics (counter and histogram)
			assert.Equal(t, 2, len(metrics))
		})
	}
}

// TestMiddlewareRequestCounting tests the request counting functionality of the middleware.
func TestMiddlewareRequestCounting(t *testing.T) {
	// Create a new registry for testing
	reg := prometheus.NewRegistry()

	// Store the default registerer
	defaultReg := prometheus.DefaultRegisterer
	// Replace default registerer with our test registry
	prometheus.DefaultRegisterer = reg
	// Restore the default registerer after the test
	defer func() { prometheus.DefaultRegisterer = defaultReg }()

	// Create middleware with test config
	config := Config{
		Namespace: "test",
		Subsystem: "api",
	}
	middleware := NewPrometheusMiddleware(config)

	// Create a test router
	router := mux.NewRouter()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}).Methods("GET")

	router.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	}).Methods("POST")

	// Wrap the router with our middleware
	routerWithMiddleware := middleware.Middleware(router)

	// Create test server
	server := httptest.NewServer(routerWithMiddleware)
	defer server.Close()

	// Make a request to /test
	resp, err := http.Get(server.URL + "/test")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Make a request to /error
	req, err := http.NewRequest("POST", server.URL+"/error", nil)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// Make a request to non-existent route
	resp, err = http.Get(server.URL + "/not-found")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	// Check metrics
	metricFamilies, err := reg.Gather()
	require.NoError(t, err)

	// Find our counter metric
	var counterFound bool
	var histogramFound bool

	expectedCounterName := "test_api_http_requests_total"
	expectedHistogramName := "test_api_http_request_duration_seconds"

	for _, mf := range metricFamilies {
		if mf.GetName() == expectedCounterName {
			counterFound = true

			// Should have 3 metrics (one for each request)
			assert.Equal(t, 3, len(mf.GetMetric()))

			// Verify labels for each metric
			labelSets := make(map[string]bool)
			for _, m := range mf.GetMetric() {
				labelSet := make(map[string]string)
				for _, l := range m.GetLabel() {
					labelSet[l.GetName()] = l.GetValue()
				}

				key := labelSet["code"] + ":" + labelSet["method"] + ":" + labelSet["path"]
				labelSets[key] = true
			}

			// Check that we have metrics for all our requests
			assert.True(t, labelSets["200:get:/test"])
			assert.True(t, labelSets["400:post:/error"])
			assert.True(t, labelSets["404:get:/not-found"])
		}

		if mf.GetName() == expectedHistogramName {
			histogramFound = true

			// Should have 3 histogram metrics (one for each request)
			assert.Equal(t, 3, len(mf.GetMetric()))
		}
	}

	assert.True(t, counterFound, "Counter metric not found")
	assert.True(t, histogramFound, "Histogram metric not found")
}

func TestDoNotUseRequestPathFor404(t *testing.T) {
	// Create a new registry for testing
	reg := prometheus.NewRegistry()

	// Store the default registerer
	defaultReg := prometheus.DefaultRegisterer
	// Replace default registerer with our test registry
	prometheus.DefaultRegisterer = reg
	// Restore the default registerer after the test
	defer func() { prometheus.DefaultRegisterer = defaultReg }()

	// Create middleware with DoNotUseRequestPathFor404 enabled
	config := Config{
		Namespace:                 "test",
		Subsystem:                 "api",
		DoNotUseRequestPathFor404: true,
	}
	middleware := NewPrometheusMiddleware(config)

	// Create a test router
	router := mux.NewRouter()

	// Wrap the router with our middleware
	routerWithMiddleware := middleware.Middleware(router)

	// Create test server
	server := httptest.NewServer(routerWithMiddleware)
	defer server.Close()

	// Make multiple requests to different non-existent routes
	paths := []string{"/not-found-1", "/not-found-2", "/some/other/path"}
	for _, path := range paths {
		resp, err := http.Get(server.URL + path)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
		resp.Body.Close()
	}

	// Check metrics
	metricFamilies, err := reg.Gather()
	require.NoError(t, err)

	// Find our counter metric
	var counterMetric *dto.MetricFamily

	expectedCounterName := "test_api_http_requests_total"

	for _, mf := range metricFamilies {
		if mf.GetName() == expectedCounterName {
			counterMetric = mf
			break
		}
	}

	require.NotNil(t, counterMetric, "Counter metric not found")

	// Should have only one metric entry for 404s
	pathCounts := make(map[string]int)
	for _, m := range counterMetric.GetMetric() {
		labelSet := make(map[string]string)
		for _, l := range m.GetLabel() {
			labelSet[l.GetName()] = l.GetValue()
		}

		if labelSet["code"] == "404" {
			pathCounts[labelSet["path"]]++
		}
	}

	// Since DoNotUseRequestPathFor404 is true, we should only have one entry with path "404"
	assert.Equal(t, 1, len(pathCounts))
	assert.Contains(t, pathCounts, "404")
	assert.Equal(t, 1, pathCounts["404"], "Counter should increase for each 404 request")
}

// TestResponseWriterDelegator tests the ResponseWriterDelegator implementation.
func TestResponseWriterDelegator(t *testing.T) {
	tests := []struct {
		name           string
		executeFunc    func(w http.ResponseWriter)
		expectedStatus int
		expectedOutput string
	}{
		{
			name: "Explicit write header",
			executeFunc: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("created"))
			},
			expectedStatus: http.StatusCreated,
			expectedOutput: "created",
		},
		{
			name: "Implicit header with write",
			executeFunc: func(w http.ResponseWriter) {
				w.Write([]byte("success"))
			},
			expectedStatus: http.StatusOK,
			expectedOutput: "success",
		},
		{
			name: "Multiple writes",
			executeFunc: func(w http.ResponseWriter) {
				w.Write([]byte("part1"))
				w.Write([]byte("part2"))
			},
			expectedStatus: http.StatusOK,
			expectedOutput: "part1part2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test response recorder
			recorder := httptest.NewRecorder()

			// Create our delegator
			delegator := &responseWriterDelegator{
				ResponseWriter: recorder,
			}

			// Execute the test function
			tt.executeFunc(delegator)

			// Check status
			assert.Equal(t, tt.expectedStatus, delegator.status)

			// Check output
			assert.Equal(t, tt.expectedOutput, recorder.Body.String())

			// Check written bytes
			assert.Equal(t, int64(len(tt.expectedOutput)), delegator.written)
		})
	}
}
