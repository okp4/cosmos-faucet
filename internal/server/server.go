package server

import (
	"net/http"
	"okp4/cosmos-faucet/pkg/client"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Config holds config of the http server.
type Config struct {
	EnableMetrics bool `mapstructure:"metrics"`
	EnableHealth  bool `mapstructure:"health"`
	Faucet        client.Faucet
	CaptchaSecret string `mapstructure:"captcha-secret"`
}

// HTTPServer exposes server methods.
type HTTPServer interface {
	Start(string)
}

type httpServer struct {
	router *mux.Router
}

// NewServer creates a new httpServer containing router.
func NewServer(config Config) HTTPServer {
	server := &httpServer{
		router: mux.NewRouter().StrictSlash(true),
	}
	server.createRoutes(config)
	return server
}

// Start starts the http server on specified address.
func (s httpServer) Start(address string) {
	log.Info().Msgf("Server listening at %s", address)
	log.Fatal().Err(http.ListenAndServe(address, s.router)).Msg("Server listening stopped")
}
