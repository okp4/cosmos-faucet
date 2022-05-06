package server

import (
    "fmt"
    "net/http"

    "github.com/rs/zerolog/log"
)

func newHealthRequestHandlerFunc() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        if _, err := w.Write([]byte("")); err != nil {
            log.Error().Msg(fmt.Sprintf("Error while sending health response: %v", err))
        }
    }
}
