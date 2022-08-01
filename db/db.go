package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"
	"gqlgen-starter/config"
)

func OpenConnection() (*sql.DB, error) {
	conn, err := sql.Open("pgx", config.Application.DatabaseURL)
	if err != nil {
		log.Error().Err(err).Msgf("Could not connect to Postgres")
		return nil, err
	}

	log.Info().Msg("Trying to ping DB...")
	if err = conn.Ping(); err != nil {
		log.Error().Err(err).Msg("Could not ping DB after opening connection")
		return nil, err
	}

	return conn, nil
}
