package server

import (
	"context"

	"okp4/cosmos-faucet/pkg/client"
)

func (s HttpServer) createRoutes(faucet *client.Faucet) {
	s.router.Path("/").
		Queries("address", "{address}").
		HandlerFunc(NewSendRequestHandlerFn(context.Background(), faucet)).
		Methods("GET")
}
