package internal

import (
	"encoding/gob"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/rs/zerolog"
	"gqlgen-starter/cmd/build"
	"gqlgen-starter/config"
	"gqlgen-starter/db"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/graph/generated"
	"gqlgen-starter/internal/graph/resolvers"
	"gqlgen-starter/internal/middleware"
	"net/http"
	"os"
	"time"
)

var sessionManager *scs.SessionManager

func StartServer() {
	logger := initZeroLogger()

	logger.Info().Msg("******************************************")
	logger.Info().Msgf("\tBuild Version: %s", build.BuildVersion)
	logger.Info().Msgf("\tBuild Commit: %s", build.BuildCommit)
	logger.Info().Msgf("\tBuild Time: %s", build.BuildTime)
	logger.Info().Msg("******************************************")

	conn, entClient, err := db.OpenConnection(logger, config.Application.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("NO DATABASE CONNECTION")
	}

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Name = middleware.SessionCookieName
	sessionManager.Store = postgresstore.New(conn)
	gob.Register(&ent.User{})

	appCtx := &app.AppContext{DB: conn, EntClient: entClient, Logger: logger, SessionManager: sessionManager}

	rootResolver := resolvers.NewRootResolver(appCtx)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: rootResolver}))

	if config.Application.IsDevelopment() {
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		logger.Info().Msgf("connect to http://localhost:%s/ for GraphQL playground", config.Application.ServerPort)
	}

	http.Handle("/query", sessionManager.LoadAndSave(middleware.AddContextUser(appCtx, srv)))

	logger.Info().Msgf("Starting GraphQL API Server at :%s 🚀", config.Application.ServerPort)
	if err = http.ListenAndServe(fmt.Sprintf(":%s", config.Application.ServerPort), nil); err != nil {
		logger.Fatal().Err(err).Msg("💀  Could not start HTTP server")
	}
}

// This should be called first, so we have a proper logger. Inspiration taken from
// https://betterstack.com/community/guides/logging/zerolog/
func initZeroLogger() *zerolog.Logger {
	logLevel, err := zerolog.ParseLevel(config.Application.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	var logger zerolog.Logger

	// Usually we'll always want logs as JSON unless we're working on our local machine
	if config.Application.StructuredLogging {
		logger = zerolog.New(os.Stdout).
			Level(logLevel).
			With().
			Timestamp().
			Str("BuildCommit", build.BuildCommit).
			Str("BuildVersion", build.BuildVersion).
			Logger()
	} else {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			Level(logLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	return &logger
}
