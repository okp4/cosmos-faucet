package server

import (
	"okp4/cosmos-faucet/graph"
	"okp4/cosmos-faucet/graph/generated"
	"okp4/cosmos-faucet/internal/server/handlers"
	"okp4/cosmos-faucet/pkg/server"

	graphql "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func (s *httpServer) createRoutes(config Config) {
	s.router.Use(handlers.PrometheusMiddleware)
	s.router.Path("/").
		HandlerFunc(playground.Handler("GraphQL playground", "/query"))
	s.router.Path("/query").
		Handler(graphql.NewDefaultServer(generated.NewExecutableSchema(generated.
			Config{Resolvers: &graph.Resolver{Faucet: faucet}})))
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
