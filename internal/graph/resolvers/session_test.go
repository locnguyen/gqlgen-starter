package resolvers

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SessionResolverSuite struct {
	suite.Suite
	TestContext
}

func TestSessionResolverSuite(t *testing.T) {
	suite.Run(t, &SessionResolverSuite{})
}

func (suite *SessionResolverSuite) SetupSuite() {
	suite.TestContext = *InitTestContext(suite.T(), "UserResolverSuite")
}

func (suite *SessionResolverSuite) TearDownSuite() {
	defer suite.pgContainer.Terminate(context.Background())
	defer suite.AppCtx.DB.Close()
}

func (suite *SessionResolverSuite) TestViewerQuery() {
	u, _ := CreateDummyUser(suite.T(), *suite.AppCtx.EntClient)

	q := `query { viewer { id email firstName lastName phoneNumber } }`
	var resp struct {
		Viewer UserObj
	}

	suite.GqlGenClient.MustPost(q, &resp, AddContextUserForTesting(u, nil))

	assert.Equal(suite.T(), u.Email, resp.Viewer.Email)
	assert.Equal(suite.T(), u.FirstName, resp.Viewer.FirstName)
	assert.Equal(suite.T(), u.LastName, resp.Viewer.LastName)
	assert.Equal(suite.T(), u.PhoneNumber, resp.Viewer.PhoneNumber)
}

func (suite *SessionResolverSuite) TestCreateSessionMutation() {
	u, pw := CreateDummyUser(suite.T(), *suite.AppCtx.EntClient)
	suite.AppCtx.Logger.Debug().Str("password", pw).Msg("Created dummy user with password")
	q := `
		mutation ($input: CreateSessionInput!) {
			createSession(input: $input) {
				token
				expiry
			}
		}
	`

	var resp struct {
		CreateSession struct {
			Token  string `json:"token"`
			Expiry string `json:"expiry"`
		} `json:"createSession"`
	}

	type input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	i := &input{
		Email:    u.Email,
		Password: pw,
	}

	suite.AppCtx.Logger.Debug().Interface("input", i).Msg("creds")
	suite.GqlGenClient.MustPost(q, &resp, client.Var("input", i), AddContextUserForTesting(nil, nil))

	assert.NotEmpty(suite.T(), resp.CreateSession.Token)
	assert.NotEmpty(suite.T(), resp.CreateSession.Expiry)
}

func (suite *SessionResolverSuite) TestDeleteSessionMutation() {
	u, _ := CreateDummyUser(suite.T(), *suite.AppCtx.EntClient)
	var resp struct {
		DeleteSession bool `json:"deleteSession"`
	}

	sess, err := suite.AppCtx.EntClient.Session.Create().
		SetData([]byte("hi")).
		SetToken(faker.UUIDHyphenated()).
		SetExpiry(time.Now().Add(1 * time.Hour)).
		Save(context.Background())

	if err != nil {
		suite.T().Error(err)
	}

	err = suite.GqlGenClient.Post(`mutation { deleteSession }`, &resp, func(bd *client.Request) {}, AddContextUserForTesting(u, &sess.Token))
	if err != nil {
		suite.AppCtx.Logger.Error().Err(err).Msg("Error POSTing deleteSession")
		suite.T().Error(err)
	}

	assert.True(suite.T(), resp.DeleteSession)
}
