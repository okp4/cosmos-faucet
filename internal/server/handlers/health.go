package handlers

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// NewHealthRequestHandlerFunc only returns 200 OK to check app status.
func NewHealthRequestHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("")); err != nil {
			log.Error().Msg(fmt.Sprintf("Error while sending health response: %v", err))
		}
	}
}
