package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestInitTracer(t *testing.T) {
	tests := []struct {
		name        string
		config      *TelemetryConfig
		expectError bool
	}{
		{
			name: "Valid configuration",
			config: &TelemetryConfig{
				Host:        "localhost",
				Port:        6831,
				ServiceName: "test-service",
			},
			expectError: false,
		},
		{
			name: "Empty host",
			config: &TelemetryConfig{
				Host:        "",
				Port:        6831,
				ServiceName: "test-service",
			},
			expectError: true,
		},
		{
			name: "Invalid port",
			config: &TelemetryConfig{
				Host:        "localhost",
				Port:        0,
				ServiceName: "test-service",
			},
			expectError: true,
		},
		{
			name: "Empty service name",
			config: &TelemetryConfig{
				Host:        "localhost",
				Port:        6831,
				ServiceName: "",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := tt.config.InitTracer()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
				assert.IsType(t, &trace.TracerProvider{}, provider)
			}
		})
	}
}

// TestOpenTelemetryMiddleware tests the OpenTelemetry middleware
func TestOpenTelemetryMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		config         *TelemetryConfig
		setupRequest   func(*http.Request)
		expectedPath   string
		expectedStatus int
	}{
		{
			name: "Valid request with route path",
			config: &TelemetryConfig{
				ServiceName: "test-service",
				Host:        "localhost",
				Port:        6831,
			},
			setupRequest: func(r *http.Request) {
				r.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
			},
			expectedPath:   "/test",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Request with custom propagator",
			config: &TelemetryConfig{
				ServiceName: "test-service",
				Host:        "localhost",
				Port:        6831,
				Propagators: propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
			},
			setupRequest: func(r *http.Request) {
				r.Header.Set("baggage", "user_id=123")
			},
			expectedPath:   "/test",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Request without route path",
			config: &TelemetryConfig{
				ServiceName: "test-service",
				Host:        "localhost",
				Port:        6831,
			},
			setupRequest:   func(r *http.Request) {},
			expectedPath:   "/test",
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "Request with nil propagator",
			config: &TelemetryConfig{
				ServiceName: "test-service",
				Host:        "localhost",
				Port:        6831,
				Propagators: nil,
			},
			setupRequest:   func(r *http.Request) {},
			expectedPath:   "/test",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := mux.NewRouter()
			router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.expectedStatus)
			})

			provider, err := tt.config.InitTracer()
			if err != nil {
				t.Fatalf("Failed to initialize tracer: %v", err)
			}
			defer func() {
				if provider != nil {
					_ = provider.Shutdown(context.Background())
				}
			}()

			tt.config.Trace = provider.Tracer("test")

			handler := tt.config.OpenTelemetryMiddleware(router)
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupRequest(req)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedPath != "" {
				assert.Equal(t, tt.expectedPath, "/test")
			}
		})
	}
}
