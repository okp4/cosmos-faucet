package server

import (
	"context"

	"okp4/cosmos-faucet/pkg/client"
)

func (s httpServer) createRoutes(faucet *client.Faucet) {
	s.router.Path("/").
		Queries("address", "{address}").
		HandlerFunc(NewSendRequestHandlerFn(context.Background(), faucet)).
		Methods("GET")
	s.router.Path("/health").
		HandlerFunc(NewHealthRequestHandlerFunc()).
		Methods("GET")
	s.router.Path("/metrics").
		HandlerFunc(NewMetricsRequestHandlerFunc(context.Background(), faucet)).
		Methods("GET")
}
