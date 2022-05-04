package server

import (
	"github.com/rs/zerolog/log"
	"net/http"

	"okp4/cosmos-faucet/pkg/client"

	"github.com/gorilla/mux"
)

// HttpServer exposes server methods
type HttpServer interface {
	Start(string)
}

type httpServer struct {
	router *mux.Router
}

// NewServer creates a new httpServer containing router
func NewServer(faucet *client.Faucet) HttpServer {
	server := httpServer{
		router: mux.NewRouter().StrictSlash(true),
	}
	server.createRoutes(faucet)
	return server
}

// Start starts the http server on specified address
func (s httpServer) Start(address string) {
	log.Info().Msgf("Server listening at %s", address)
	log.Fatal().Err(http.ListenAndServe(address, s.router)).Msg("Server listening stopped")
}
