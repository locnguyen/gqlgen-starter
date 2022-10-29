package resolvers

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartPgContainer(logger testcontainers.Logging, containerName string) (testcontainers.Container, *string, error) {
	ctx := context.Background()
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
		Logger:       logger,
		Reuse:        true,
	})

	if err != nil {
		logger.Printf("Error starting Postgres container for integration testing")
		return nil, nil, err
	}

	mappedPort, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		logger.Printf("Error geting mapped port for test DB container")
		return nil, nil, err
	}

	databaseURL := fmt.Sprintf("postgres://postgres:postgres@localhost:%s/test?sslmode=disable", mappedPort.Port())

	return postgresC, &databaseURL, nil
}
