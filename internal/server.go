package internal

import (
	"gqlgen-starter/config"
	"gqlgen-starter/internal/graph/generated"
	"gqlgen-starter/internal/graph/resolvers"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func StartServer() {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.RootResolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", config.Application.ServerPort)
	log.Fatal(http.ListenAndServe(":"+config.Application.ServerPort, nil))
}
