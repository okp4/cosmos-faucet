package server

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

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func initPrometheus() {
	_ = prometheus.Register(totalRequests)
	_ = prometheus.Register(responseStatus)
	_ = prometheus.Register(httpDuration)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	}, []string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests.",
	},
	[]string{"path"},
)

func prometheusMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))

		rw := newResponseWriter(w)
		handler.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		if !(path == "/health" || path == "/metrics") {
			responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
			totalRequests.WithLabelValues(path).Inc()
			timer.ObserveDuration()
		}
	})
}

// NewMetricsRequestHandler returns a REST handler returning useful prometheus metrics
func NewMetricsRequestHandler() http.Handler {
	return promhttp.Handler()
}
