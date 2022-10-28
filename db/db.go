package db

import (
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"gqlgen-starter/config"
	"gqlgen-starter/internal/ent"
)

func OpenConnection(logger *zerolog.Logger) (*sql.DB, *ent.Client, error) {
	logger.Info().Msg("Trying to open database connection...")
	conn, err := sql.Open("pgx", config.Application.DatabaseURL)
	if err != nil {
		logger.Error().Err(err).Msgf("Could not connect to Postgres")
		return nil, nil, err
	}

	logger.Info().Msg("Trying to ping DB...")
	if err = conn.Ping(); err != nil {
		logger.Error().Err(err).Msg("Could not ping DB after opening connection")
		return nil, nil, err
	}
	logger.Info().Msg("Database ping succeeded")

	entClient := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, conn)))
	return conn, entClient, nil
}
