package resolvers

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/alexedwards/scs/v2"
	"github.com/go-faker/faker/v4"
	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
	"gqlgen-starter/db"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/graph/generated"
	"gqlgen-starter/internal/middleware"
	"os"
	"testing"
)

type TestContext struct {
	GqlGenClient *client.Client
	AppCtx       *app.AppContext
	pgContainer  testcontainers.Container
}

// In the real server a ContextCookie is always created in the CookieAuth middleware
// In unit tests this does not happen upstream, so we need to provide it
func AddContextUserForTesting(user *ent.User, sid *string) client.Option {
	return func(bd *client.Request) {
		//ctxCookie := &middleware.ContextCookie{
		//	User:   user,
		//	Writer: httptest.NewRecorder(),
		//	Sid:    sid,
		//}

		ctx := context.WithValue(context.Background(), middleware.ContextUserKey, user)
		bd.HTTP = bd.HTTP.WithContext(ctx)
	}
}

func InitTestContext(t *testing.T, testName string) *TestContext {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	ctx := context.Background()

	postgresC, databaseURL, err := StartPgContainer(&logger, fmt.Sprintf("%s_db", testName))
	if err != nil {
		t.Error(err)
	}

	dbConn, entClient, err := db.OpenConnection(&logger, *databaseURL)

	if err := entClient.Schema.Create(ctx); err != nil {
		t.Fatal(err)
	}

	sessionManager := scs.New()
	gob.Register(&ent.User{})
	appCtx := &app.AppContext{
		DB:             dbConn,
		EntClient:      entClient,
		Logger:         &logger,
		SessionManager: sessionManager,
	}

	gqlGenClient := client.New(sessionManager.LoadAndSave(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: NewRootResolver(appCtx)}))))

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

func CreateDummyUser(t *testing.T, client ent.Client) (*ent.User, string) {
	pw := faker.Password()
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		t.Error(err)
	}
	u, err := client.User.Create().
		SetEmail(faker.Email()).
		SetFirstName(faker.FirstName()).
		SetLastName(faker.LastName()).
		SetHashedPassword(hashedPw).
		SetPhoneNumber(faker.Phonenumber()).
		Save(context.Background())

	if err != nil {
		t.Error(err)
	}
	return u, pw
}
