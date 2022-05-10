package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// PrometheusMiddleware adds prometheus metrics on every request to a router.
func PrometheusMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))

		rw := newResponseWriter(w)
		handler.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		totalRequests.WithLabelValues(methods[0], path, strconv.Itoa(statusCode)).Inc()
		timer.ObserveDuration()
	})
}

// NewMetricsRequestHandler exposes prometheus metrics.
func NewMetricsRequestHandler() http.Handler {
	return promhttp.Handler()
}

var totalRequests = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "HTTP requests counter.",
	}, []string{"method", "path", "code"},
)

var httpDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of HTTP requests.",
	},
	[]string{"path"},
)
