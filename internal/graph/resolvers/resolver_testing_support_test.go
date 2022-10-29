package resolvers

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gqlgen-starter/db"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/graph/generated"
	"os"
	"testing"
)

type TestContext struct {
	GqlGenClient *client.Client
	AppCtx       *app.AppContext
	pgContainer  testcontainers.Container
}

func InitTestContext(t *testing.T) *TestContext {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	ctx := context.Background()

	postgresC, databaseURL, err := StartPgContainer(&logger, "user_test_db")
	if err != nil {
		t.Error(err)
	}

	dbConn, entClient, err := db.OpenConnection(&logger, *databaseURL)

	if err := entClient.Schema.Create(ctx); err != nil {
		t.Fatal(err)
	}

	appCtx := &app.AppContext{
		DB:        dbConn,
		EntClient: entClient,
		Logger:    &logger,
	}

	gqlGenClient := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: NewRootResolver(appCtx)})))

	return &TestContext{
		GqlGenClient: gqlGenClient,
		AppCtx:       appCtx,
		pgContainer:  postgresC,
	}
}

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
