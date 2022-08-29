package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewMetricsRequestHandler exposes prometheus metrics.
func NewMetricsRequestHandler() http.Handler {
	return promhttp.Handler()
}
