package resolvers

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-faker/faker/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"gqlgen-starter/db"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/ent/session"
	"gqlgen-starter/internal/ent/user"
	"gqlgen-starter/internal/graph/generated"
	"os"
	"testing"
)

type UserResolverSuite struct {
	suite.Suite
	GqlGenClient *client.Client
	AppCtx       *app.AppContext
	pgContainer  testcontainers.Container
}

func TestUserResolverSuite(t *testing.T) {
	suite.Run(t, new(UserResolverSuite))
}

func (suite *UserResolverSuite) SetupSuite() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	ctx := context.Background()

	postgresC, databaseURL, err := StartPgContainer(&logger, "user_test_db")
	if err != nil {
		suite.T().Error(err)
	}
	suite.pgContainer = postgresC

	dbConn, entClient, err := db.OpenConnection(&logger, *databaseURL)

	if err := entClient.Schema.Create(ctx); err != nil {
		suite.T().Fatal(err)
	}

	suite.AppCtx = &app.AppContext{
		DB:        dbConn,
		EntClient: entClient,
		Logger:    &logger,
	}

	suite.GqlGenClient = client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: NewRootResolver(suite.AppCtx)})))
}

func (suite *UserResolverSuite) TearDownSuite() {
	defer suite.pgContainer.Terminate(context.Background())
	defer suite.AppCtx.DB.Close()
}

func (suite *UserResolverSuite) TestUserQuery() {
	u, err := suite.AppCtx.EntClient.User.Create().
		SetEmail(faker.Email()).
		SetFirstName(faker.FirstName()).
		SetLastName(faker.LastName()).
		SetHashedPassword([]byte(faker.Password())).
		SetPhoneNumber(faker.Phonenumber()).
		Save(context.Background())
	if err != nil {
		suite.T().Error(err)
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

	suite.GqlGenClient.MustPost("query GetUser($id: ID!) { user(id: $id) { id firstName lastName phoneNumber email } }", &resp, client.Var("id", u.ID))
	assert.Equal(suite.T(), resp.User.ID, fmt.Sprint(u.ID))
	assert.Equal(suite.T(), resp.User.FirstName, u.FirstName)
	assert.Equal(suite.T(), resp.User.LastName, u.LastName)
	assert.Equal(suite.T(), resp.User.PhoneNumber, u.PhoneNumber)
	assert.Equal(suite.T(), resp.User.Email, u.Email)
}

func (suite *UserResolverSuite) TestCreateUserMutation() {
	var resp struct {
		CreateUser struct {
			Sid    string `json:"sid"`
			Expiry string `json:"expiry"`
		}
	}

	suite.GqlGenClient.MustPost(`mutation { createUser(input: { firstName: "Natasha" lastName: "Romanova" email: "blackwidow@avengers.com" phoneNumber: "+8888888888" password: "P@ssw0rd!" passwordConfirmation: "P@ssw0rd!" }) { sid expiry } }`, &resp)
	assert.NotEmpty(suite.T(), resp.CreateUser.Sid)
	assert.NotEmpty(suite.T(), resp.CreateUser.Expiry)

	sess, err := suite.AppCtx.EntClient.Session.Query().
		Where(session.Sid(resp.CreateUser.Sid)).
		Only(context.Background())

	if err != nil {
		suite.T().Error(err)
	}

	assert.Equal(suite.T(), resp.CreateUser.Sid, sess.Sid)

	subject, err := suite.AppCtx.EntClient.User.Query().
		Where(user.Email("blackwidow@avengers.com")).
		Only(context.Background())

	if err != nil {
		suite.T().Error(err)
	}

	assert.Equal(suite.T(), "blackwidow@avengers.com", subject.Email)
}
