package internal

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alexedwards/scs/v2"
	"gqlgen-starter/config"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/gql/generated"
	"gqlgen-starter/internal/graph/directives"
	"gqlgen-starter/internal/graph/resolvers"
	"gqlgen-starter/internal/middleware"
	"net/http"
)

var sessionManager *scs.SessionManager

func StartServer() {
	log := app.GetLogger()
	appCtx, err := Initialize()
	if err != nil {
		log.Error().
			Err(err).
			Msg("🚨  could not initialize application context  🚨")
		return
	}

	schemaConfig := generated.Config{
		Resolvers:  resolvers.NewRootResolver(appCtx),
		Directives: *directives.NewDirectiveRoot(),
	}
	srv := resolvers.CreateGqlServer(appCtx, &schemaConfig)

	if config.Application.IsDevelopment() {
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		log.Info().
			Msgf("connect to http://localhost:%s/ for GraphQL playground", config.Application.ServerPort)
	}

	http.Handle("/query", middleware.RefererServer(middleware.CorrelationID(middleware.HttpLogger(sessionManager.LoadAndSave(middleware.AddContextViewer(appCtx, srv))))))

	log.Info().
		Msgf("🚀  starting GraphQL Server at :%s", config.Application.ServerPort)

	if err = http.ListenAndServe(fmt.Sprintf(":%s", config.Application.ServerPort), nil); err != nil {
		log.Fatal().
			Err(err).
			Msg("💀  could not start HTTP server")
	}
}
