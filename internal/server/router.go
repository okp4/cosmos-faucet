package server

import (
	"okp4/cosmos-faucet/internal/server/handlers"
	"okp4/cosmos-faucet/pkg/server"
)

func (s *httpServer) createRoutes(config Config) {
	s.router.Use(handlers.PrometheusMiddleware)
	s.router.Path("/").
		Queries("address", "{address}").
		HandlerFunc(server.NewSendRequestHandlerFn(config.Faucet)).
		Methods("GET")
	if config.EnableHealth {
		s.router.Path("/health").
			HandlerFunc(handlers.NewHealthRequestHandlerFunc()).
			Methods("GET")
	}
	if config.EnableMetrics {
		s.router.Path("/metrics").
			Handler(handlers.NewMetricsRequestHandler()).
			Methods("GET")
	}
}
