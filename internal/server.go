package internal

import (
	"gqlgen-starter/config"
	"gqlgen-starter/internal/graph"
	"gqlgen-starter/internal/graph/generated"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func StartServer() {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", config.Application.ServerPort)
	log.Fatal(http.ListenAndServe(":"+config.Application.ServerPort, nil))
}
