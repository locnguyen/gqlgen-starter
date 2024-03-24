package internal

import (
	"context"
	"encoding/gob"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"gqlgen-starter/cmd/build"
	"gqlgen-starter/config"
	"gqlgen-starter/db"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/app/loaders"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/graph/directives"
	"gqlgen-starter/internal/graph/generated"
	"gqlgen-starter/internal/graph/resolvers"
	"gqlgen-starter/internal/middleware"
	"net/http"
	"time"
)

var sessionManager *scs.SessionManager

func StartServer() {
	log := app.GetLogger()

	log.Info().Msg("******************************************")
	log.Info().Msgf("\tBuild Version: %s", build.BuildVersion)
	log.Info().Msgf("\tBuild Commit: %s", build.BuildCommit)
	log.Info().Msgf("\tBuild Time: %s", build.BuildTime)
	log.Info().Msg("******************************************")

	ctx := log.WithContext(context.Background())
	pgConn, err := db.OpenPostgresConn(ctx, config.Application.DatabaseURL)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("☠️  NO DATABASE CONNECTION  ☠️")
		return
	}

	entClient := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, pgConn)))

	redisPool, err := db.NewRedisPool(ctx, config.Application)
	if err != nil {
		log.Error().
			Err(err).
			Msg("💀  Could not create Redis connection pool  💀")
		return
	}

	sessionManager = scs.New()
	sessionManager.Lifetime = 48 * time.Hour
	sessionManager.Cookie.Name = "session-viewer"
	sessionManager.Store = redisstore.New(redisPool)
	gob.Register(&ent.User{})

	appCtx := &app.AppContext{
		DB:             pgConn,
		EntClient:      entClient,
		Loaders:        loaders.NewLoaders(entClient),
		Logger:         &log,
		SessionManager: sessionManager,
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
