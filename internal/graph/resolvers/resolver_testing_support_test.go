package resolvers

import (
	"context"
	"encoding/gob"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/alexedwards/scs/v2"
	"github.com/go-faker/faker/v4"
	"github.com/testcontainers/testcontainers-go"
	"golang.org/x/crypto/bcrypt"
	"gqlgen-starter/db"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/app/loaders"
	"gqlgen-starter/internal/app/models"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/gql/generated"
	"gqlgen-starter/internal/graph/directives"
	"gqlgen-starter/internal/testsupport"
	"testing"
)

type TestContext struct {
	GqlGenClient *client.Client
	AppCtx       *app.AppContext
	pgContainer  testcontainers.Container
}

func AddContextViewerForTesting(u *ent.User) client.Option {
	return func(bd *client.Request) {
		ctx := context.WithValue(context.Background(), ent.ContextViewerKey, &ent.Viewer{User: u})
		bd.HTTP = bd.HTTP.WithContext(ctx)
	}
}

func InitTestContext(t *testing.T, testName string) *TestContext {
	log := app.GetLogger()
	ctx := context.Background()
	postgresC, databaseURL, err := testsupport.StartPgContainer(ctx, fmt.Sprintf("%s_db", testName))
	if err != nil {
		t.Error(err)
	}

	dbConn, err := db.OpenPostgresConn(ctx, *databaseURL)
	if err != nil {
		log.Error().
			Err(err).
			Str("databaseURL", *databaseURL).
			Msg("starting PG container for unit testing")
		t.Error(err)
	}
	entClient := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, dbConn)))
	if err := entClient.Schema.Create(ctx); err != nil {
		t.Fatal(err)
	}

	sessionManager := scs.New()
	gob.Register(&ent.User{})
	appCtx := &app.AppContext{
		DB:             dbConn,
		EntClient:      entClient,
		Loaders:        loaders.NewLoaders(entClient),
		Logger:         &log,
		SessionManager: sessionManager,
	}

	schemaConfig := generated.Config{
		Resolvers:  NewRootResolver(appCtx),
		Directives: *directives.NewDirectiveRoot(),
	}

	gqlServer := CreateGqlServer(appCtx, &schemaConfig)
	gqlGenClient := client.New(appCtx.SessionManager.LoadAndSave(gqlServer))

	return &TestContext{
		GqlGenClient: gqlGenClient,
		AppCtx:       appCtx,
		pgContainer:  postgresC,
	}
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
		SetPhoneNumber(faker.Phonenumber()).
		SetRoles([]models.Role{models.RoleGenPop}).
		Save(context.Background())

	if err != nil {
		t.Error(err)
	}
	return u, pw
}
