package server

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

type prometheusHandler struct {
	totalRequests  *prometheus.CounterVec
	responseStatus *prometheus.CounterVec
	httpDuration   *prometheus.HistogramVec
}

func (p prometheusHandler) init() {
	_ = prometheus.Register(p.totalRequests)
	_ = prometheus.Register(p.responseStatus)
	_ = prometheus.Register(p.httpDuration)
}

func newPrometheusHandler() prometheusHandler {
	handler := prometheusHandler{
		totalRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Number of get requests.",
			}, []string{"path"},
		),
		responseStatus: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "response_status",
				Help: "Status of HTTP response",
			},
			[]string{"status"},
		),
		httpDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_response_time_seconds",
			Help: "Duration of HTTP requests.",
		}, []string{"path"}),
	}
	handler.init()
	return handler
}

func prometheusMiddleware(handler http.Handler) http.Handler {
	promHandler := newPrometheusHandler()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(promHandler.httpDuration.WithLabelValues(path))

		rw := newResponseWriter(w)
		handler.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		if !(path == "/metrics" || path == "/health") {
			promHandler.responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
			promHandler.totalRequests.WithLabelValues(path).Inc()
			timer.ObserveDuration()
		}
	})
}

// NewMetricsRequestHandlerFunc returns a REST handler func returning useful prometheus metrics
func NewMetricsRequestHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler()
	}
}
