package server

import (
	"net/http"
	"okp4/cosmos-faucet/graph"
	"okp4/cosmos-faucet/graph/generated"
	"okp4/cosmos-faucet/internal/server/handlers"
	"okp4/cosmos-faucet/pkg/captcha"
	"time"

	"github.com/99designs/gqlgen/graphql"
	graphqlserver "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
)

func (s *httpServer) createRoutes(config Config) {
	s.router.Path("/").
		HandlerFunc(playground.Handler("GraphQL playground", "/graphql")).
		Methods("GET")
	s.router.Path("/graphql").
		Handler(
			newGraphQLServer(
				generated.NewExecutableSchema(generated.Config{
					Resolvers: &graph.Resolver{
						Faucet:          config.Faucet,
						CaptchaResolver: captcha.NewCaptchaResolver(config.CaptchaConf),
					},
				}),
			),
		).
		Methods("GET", "POST", "OPTIONS")
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
