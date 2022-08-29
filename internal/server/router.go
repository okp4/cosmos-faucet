package server

import (
	"net/http"
	"okp4/cosmos-faucet/graph"
	"okp4/cosmos-faucet/graph/generated"
	"okp4/cosmos-faucet/internal/server/handlers"
	"time"

	"github.com/99designs/gqlgen/graphql"
	graphqlserver "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
)

func (s *httpServer) createRoutes(graphqlResolver *graph.Resolver, health, metrics bool) {
	s.router.Path("/").
		HandlerFunc(playground.Handler("GraphQL playground", "/graphql")).
		Methods("GET")
	s.router.Path("/graphql").
		Handler(
			newGraphQLServer(
				generated.NewExecutableSchema(generated.Config{
					Resolvers: graphqlResolver,
				}),
			),
		).
		Methods("GET", "POST", "OPTIONS")
	if health {
		s.router.Path("/health").
			HandlerFunc(handlers.NewHealthRequestHandlerFunc()).
			Methods("GET")
	}
	if metrics {
		s.router.Path("/metrics").
			Handler(handlers.NewMetricsRequestHandler()).
			Methods("GET")
	}
}

func newGraphQLServer(schema graphql.ExecutableSchema) http.Handler {
	srv := graphqlserver.New(schema)

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return srv
}
