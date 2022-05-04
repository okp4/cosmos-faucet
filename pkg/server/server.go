package server

import (
	"github.com/rs/zerolog/log"
	"net/http"

	"okp4/cosmos-faucet/pkg/client"

	"github.com/gorilla/mux"
)

type HttpServer struct {
	router *mux.Router
}

// NewServer creates a new httpServer containing router
func NewServer(faucet *client.Faucet) HttpServer {
	server := HttpServer{}

	server.router = mux.NewRouter().StrictSlash(true)
	server.createRoutes(faucet)

	return server
}

// Start starts the http server on specified address
func (s HttpServer) Start(address string) {
	log.Info().Msgf("Server listening at %s", address)
	log.Fatal().Err(http.ListenAndServe(address, s.router)).Msg("Server listening stopped")
}
