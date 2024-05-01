package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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

// StartPgContainer will create Postgres container for testing
func StartPgContainer(ctx context.Context, containerName string) (testcontainers.Container, *string, error) {
	log := zerolog.Ctx(ctx)
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:13-alpine",
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForExposedPort(),
			Name:         containerName,
			Env: map[string]string{
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_DB":       "test",
			},
		},
		Started:      true,
		ProviderType: 0,
		Logger:       log,
		Reuse:        true,
	})

	if err != nil {
		log.Error().
			Err(err).
			Msg("starting Postgres container for testing")
		return nil, nil, err
	}

	mappedPort, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		log.Error().
			Err(err).
			Msg("getting mapped port for test DB container")
		return nil, nil, err
	}

	connStr := fmt.Sprintf("postgres://postgres:postgres@localhost:%s/test?sslmode=disable", mappedPort.Port())

	return postgresC, &connStr, nil
}
