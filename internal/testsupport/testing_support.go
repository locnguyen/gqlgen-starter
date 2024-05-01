package testsupport

import (
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
	"gqlgen-starter/internal/app/models"
	"gqlgen-starter/internal/ent"
	"testing"
)

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

	databaseURL := fmt.Sprintf("postgres://postgres:postgres@localhost:%s/test?sslmode=disable", mappedPort.Port())

	return postgresC, &databaseURL, nil
}

func CreateDummyUser(t *testing.T, client *ent.Client) (*ent.User, string) {
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
		SetPhoneNumber(faker.E164PhoneNumber()).
		SetRoles([]models.Role{models.RoleGenPop}).
		Save(context.Background())

	if err != nil {
		t.Error(err)
	}
	return u, pw
}
