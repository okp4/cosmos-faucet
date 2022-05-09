package server

import (
	"okp4/cosmos-faucet/internal/server/handlers"
	"okp4/cosmos-faucet/pkg/client"
	"okp4/cosmos-faucet/pkg/server"
)

func (s *httpServer) createRoutes(faucet *client.Faucet) {
	s.router.Use(handlers.PrometheusMiddleware)
	s.router.Path("/").
		Queries("address", "{address}").
		HandlerFunc(server.NewSendRequestHandlerFn(faucet)).
		Methods("GET")
	s.router.Path("/health").
		HandlerFunc(handlers.NewHealthRequestHandlerFunc()).
		Methods("GET")
	s.router.Path("/metrics").
		Handler(handlers.NewMetricsRequestHandler()).
		Methods("GET")
}
