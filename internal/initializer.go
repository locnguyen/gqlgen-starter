package internal

import (
	"context"
	"encoding/gob"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/alexedwards/scs/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"gqlgen-starter/cmd/build"
	"gqlgen-starter/config"
	"gqlgen-starter/db"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/app/loaders"
	"gqlgen-starter/internal/app/sessionstore"
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

	nc, err := nats.Connect(config.Application.NatsURL)
	if err != nil {
		log.Error().
			Err(err).
			Msg("💀  could not connect to NATS server  💀")
		return nil, err
	} else {
		log.Info().
			Msgf("📣  NATS status: %s", nc.Status().String())
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Error().
			Err(err).
			Msg("💀  could not create JetStream context from NATS connection  💀")
		return nil, err
	}

	sessionManager = scs.New()
	sessionManager.Lifetime = 48 * time.Hour
	sessionManager.Cookie.Name = "session-cookie"
	sessionManager.Store, err = sessionstore.New(js)
	if err != nil {
		return nil, err
	}
	log.Info().
		Msg("🍪  created session manager linked to NATS bucket")
	gob.Register(&ent.User{})

	appCtx := &app.AppContext{
		DB:             pgConn,
		EntClient:      entClient,
		Loaders:        loaders.NewLoaders(entClient),
		Logger:         &log,
		Nats:           nc,
		SessionManager: sessionManager,
	}

	return appCtx, nil
}
