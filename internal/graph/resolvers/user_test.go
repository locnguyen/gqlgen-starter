package resolvers

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-faker/faker/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gqlgen-starter/config"
	"gqlgen-starter/db"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/ent/user"
	"gqlgen-starter/internal/graph/generated"
	"os"
	"testing"
)

func TestUserResolvers(t *testing.T) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})

	logger.Debug().Msg("Starting Postgres container for integration testing....")
	ctx := context.Background()
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:13-alpine",
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForExposedPort(),
			Name:         "integrationtestingdb",
			Env: map[string]string{
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_DB":       "test",
			},
		},
		Started:      true,
		ProviderType: 0,
		Logger:       &logger,
		Reuse:        true,
	})

	if err != nil {
		t.Error(err)
	}

	defer postgresC.Terminate(ctx)

	mappedPort, err := postgresC.MappedPort(ctx, "5432")
	config.Application.DatabaseURL = fmt.Sprintf("postgres://postgres:postgres@localhost:%s/test?sslmode=disable", mappedPort.Port())
	dbConn, entClient, err := db.OpenConnection(&logger)

	if err := entClient.Schema.Create(ctx); err != nil {
		t.Fatal(err)
	}

	appCtx := &app.AppContext{
		DB:        dbConn,
		EntClient: entClient,
		Logger:    &logger,
	}

	gqlGenClient := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: NewRootResolver(appCtx)})))

	t.Run("user(id: ID!) query", func(t *testing.T) {
		u, err := entClient.User.Create().
			SetEmail(faker.Email()).
			SetFirstName(faker.FirstName()).
			SetLastName(faker.LastName()).
			SetHashedPassword([]byte(faker.Password())).
			SetPhoneNumber(faker.Phonenumber()).
			Save(ctx)
		if err != nil {
			t.Error(err)
		}

		var resp struct {
			User struct {
				ID          string `json:"id"`
				FirstName   string `json:"firstName"`
				LastName    string `json:"lastName"`
				PhoneNumber string `json:"phoneNumber"`
				Email       string `json:"email"`
			}
		}

		gqlGenClient.MustPost("query GetUser($id: ID!) { user(id: $id) { id firstName lastName phoneNumber email } }", &resp, client.Var("id", u.ID))
		assert.Equal(t, resp.User.ID, fmt.Sprint(u.ID))
		assert.Equal(t, resp.User.FirstName, u.FirstName)
		assert.Equal(t, resp.User.LastName, u.LastName)
		assert.Equal(t, resp.User.PhoneNumber, u.PhoneNumber)
		assert.Equal(t, resp.User.Email, u.Email)
	})

	t.Run("createUser mutation", func(t *testing.T) {
		var resp struct {
			CreateUser struct {
				Sid    string `json:"sid"`
				Expiry string `json:"expiry"`
			}
		}

		gqlGenClient.MustPost(`mutation { createUser(input: { firstName: "Natasha" lastName: "Romanova" email: "blackwidow@avengers.com" phoneNumber: "+8888888888" password: "P@ssw0rd!" passwordConfirmation: "P@ssw0rd!" }) { sid expiry } }`, &resp)
		assert.NotEmpty(t, resp.CreateUser.Sid)
		assert.NotEmpty(t, resp.CreateUser.Expiry)

		subject, err := entClient.User.Query().
			Where(user.Email("blackwidow@avengers.com")).
			Only(ctx)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, "blackwidow@avengers.com", subject.Email)
	})
}
