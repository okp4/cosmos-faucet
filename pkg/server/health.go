package server

import "net/http"

// NewHealthRequestHandlerFunc returns a REST handler func to check app status
func NewHealthRequestHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(""))
	}
}
