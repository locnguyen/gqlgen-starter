package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"gqlgen-starter/config"
)

func OpenConnection(logger *zerolog.Logger) (*sql.DB, error) {
	conn, err := sql.Open("pgx", config.Application.DatabaseURL)
	if err != nil {
		logger.Error().Err(err).Msgf("Could not connect to Postgres")
		return nil, err
	}

	logger.Info().Msg("Trying to ping DB...")
	if err = conn.Ping(); err != nil {
		logger.Error().Err(err).Msg("Could not ping DB after opening connection")
		return nil, err
	}
	logger.Info().Msg("Database ping succeeded")
	return conn, nil
}
