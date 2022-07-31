package internal

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gqlgen-starter/config"
	"gqlgen-starter/internal/graph/generated"
	"gqlgen-starter/internal/graph/resolvers"
	"net/http"
	"os"
)

func StartServer() {
	initZeroLogger()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Info().Msgf("connect to http://localhost:%s/ for GraphQL playground", config.Application.ServerPort)
	err := http.ListenAndServe(":"+config.Application.ServerPort, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start HTTP server")
	}
	log.Info().Msg("Started GraphQL API Server 🚀")
}

// This should be called first, so we have a proper logger
func initZeroLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel, err := zerolog.ParseLevel(config.Application.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	if config.Application.GoEnv == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
