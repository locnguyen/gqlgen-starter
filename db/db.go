package db

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
)

func OpenPostgresConn(ctx context.Context, databaseURL string) (*sql.DB, error) {
	log := zerolog.Ctx(ctx)
	log.Info().
		Msg("opening database connection...")

	conn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("could not connect to Postgres")
		return nil, err
	}

	log.Info().
		Msg("pinging Postgres...")
	if err := conn.Ping(); err != nil {
		log.Error().
			Err(err).
			Msg("🚨️   could not ping Postgres after opening connection")
		return nil, err
	}
	log.Info().Msg("☎️️   Postgres ping succeeded!")

	return conn, nil
}
