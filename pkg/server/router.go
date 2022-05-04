package server

import (
	"context"

	"okp4/cosmos-faucet/pkg/client"
)

func (s *httpServer) createRoutes(faucet *client.Faucet) {
	s.router.Use(prometheusMiddleware)
	s.router.Path("/").
		Queries("address", "{address}").
		HandlerFunc(newSendRequestHandlerFn(context.Background(), faucet)).
		Methods("GET")
	s.router.Path("/health").
		HandlerFunc(newHealthRequestHandlerFunc()).
		Methods("GET")
	s.router.Path("/metrics").
		Handler(newMetricsRequestHandler()).
		Methods("GET")
}
