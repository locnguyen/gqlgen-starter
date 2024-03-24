package internal

import (
	"context"
	"encoding/gob"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"gqlgen-starter/cmd/build"
	"gqlgen-starter/config"
	"gqlgen-starter/db"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/app/loaders"
	"gqlgen-starter/internal/ent"
	"time"
)

func Initialize() (*app.AppContext, error) {
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
		return nil, err
	}

	entClient := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, pgConn)))

	redisPool, err := db.NewRedisPool(ctx, config.Application)
	if err != nil {
		log.Error().
			Err(err).
			Msg("💀  could not create Redis connection pool  💀")
		return nil, err
	}

	sessionManager = scs.New()
	sessionManager.Lifetime = 48 * time.Hour
	sessionManager.Cookie.Name = "session-cookie"
	sessionManager.Store = redisstore.New(redisPool)
	gob.Register(&ent.User{})

	appCtx := &app.AppContext{
		DB:             pgConn,
		EntClient:      entClient,
		Loaders:        loaders.NewLoaders(entClient),
		Logger:         &log,
		SessionManager: sessionManager,
	}

	return appCtx, nil
}
