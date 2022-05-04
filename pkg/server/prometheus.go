package server

import (
    "context"
    "net/http"

    "okp4/cosmos-faucet/pkg/client"
)

// NewHealthRequestHandlerFunc returns a REST handler func to check app status
func NewHealthRequestHandlerFunc() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte(""))
    }
}

// NewMetricsRequestHandlerFunc returns a REST handler func returning useful prometheus metrics
func NewMetricsRequestHandlerFunc(ctx context.Context, faucet *client.Faucet) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    }
}
