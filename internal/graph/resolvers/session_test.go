package resolvers

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gqlgen-starter/internal/ent/session"
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

	suite.GqlGenClient.MustPost(q, &resp, AddContextCookieForTesting(u, nil))

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
				sid
				expiry
			}
		}
	`

	var resp struct {
		CreateSession struct {
			Sid    string `json:"sid"`
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
	suite.GqlGenClient.MustPost(q, &resp, client.Var("input", i), AddContextCookieForTesting(nil, nil))

	assert.NotEmpty(suite.T(), resp.CreateSession.Sid)
	assert.NotEmpty(suite.T(), resp.CreateSession.Expiry)
}

func (suite *SessionResolverSuite) TestDeleteSessionMutation() {
	u, _ := CreateDummyUser(suite.T(), *suite.AppCtx.EntClient)
	var resp struct {
		DeleteSession struct {
			Sid    string `json:"sid"`
			Expiry string `json:"expiry"`
		}
	}

	sess, err := suite.AppCtx.EntClient.Session.Create().
		SetUser(u).
		SetType(session.TypeGeneral).
		SetSid(faker.UUIDHyphenated()).
		SetDeleted(false).
		SetExpiry(time.Now().Add(1 * time.Hour)).
		Save(context.Background())

	if err != nil {
		suite.T().Error(err)
	}

	err = suite.GqlGenClient.Post(`mutation { deleteSession { sid expiry } }`, &resp, func(bd *client.Request) {}, AddContextCookieForTesting(u, &sess.Sid))

	if err != nil {
		suite.AppCtx.Logger.Error().Err(err).Msg("Error POSTing deleteSession")
		suite.T().Error(err)
	}

	sess, err = suite.AppCtx.EntClient.Session.Get(context.Background(), sess.ID)
	assert.True(suite.T(), sess.Deleted)
}
