// Package middleware provides HTTP middleware components for the application.
// It includes authentication, CORS, logging, and telemetry middleware.
package middleware

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger" //nolint:staticcheck // TODO: migrate to OTLP
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryConfig struct {
	Host          string
	Port          int
	ServiceName   string
	Trace         trace.Tracer
	Propagators   propagation.TextMapPropagator
	TraceProvider *tracesdk.TracerProvider
}

// InitTracer initializes the Jaeger tracer
func (c *TelemetryConfig) InitTracer() (*tracesdk.TracerProvider, error) {
	// Validate the configuration
	if err := c.validateConfig(); err != nil {
		return nil, err
	}

	// Set the tracer
	c.Trace = otel.Tracer(c.ServiceName)

	// Create the Jaeger exporter
	exp, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(c.Host),
			jaeger.WithAgentPort(strconv.Itoa(c.Port)),
		),
	)
	if err != nil {
		return nil, err
	}
	provider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(c.ServiceName),
		)),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return provider, nil
}

// OpenTelemetryMiddleware is a middleware for tracing HTTP requests
func (c *TelemetryConfig) OpenTelemetryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		delegate := &responseWriterDelegator{
			ResponseWriter: w,
			status:         http.StatusOK, // Initialize with default status
		}

		if c.Propagators == nil {
			c.Propagators = otel.GetTextMapPropagator()
		}

		var path string

		cr := mux.CurrentRoute(r)
		if cr != nil {
			var err error
			path, err = cr.GetPathTemplate()
			if err != nil {
				// don't create traces for 404 pages in gorilla mux
				next.ServeHTTP(delegate, r)
				return
			}
		}

		ctx := c.Propagators.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		opts := []trace.SpanStartOption{
			trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
			trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
			trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(c.ServiceName, path, r)...),
			trace.WithSpanKind(trace.SpanKindServer),
		}

		ctx, span := c.Trace.Start(ctx, path, opts...)
		defer span.End()

		next.ServeHTTP(delegate, r.WithContext(ctx))

		span.SetAttributes(attribute.Int("http.status", delegate.status))
	})
}

// validateConfig validates the configuration
func (c *TelemetryConfig) validateConfig() error {
	if c.Host == "" {
		return errors.New("host is required")
	}
	if c.Port <= 0 {
		return errors.New("invalid port")
	}
	return nil
}
