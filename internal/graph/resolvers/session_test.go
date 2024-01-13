package resolvers

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SessionResolverSuite struct {
	suite.Suite
	TestContext
}

func TestSessionResolverSuite(t *testing.T) {
	suite.Run(t, &SessionResolverSuite{})
}

func (s *SessionResolverSuite) SetupSuite() {
	s.TestContext = *InitTestContext(s.T(), "SessionResolverSuite")
}

func (s *SessionResolverSuite) TearDownSuite() {
	defer s.pgContainer.Terminate(context.Background())
	defer s.AppCtx.DB.Close()
}

func (s *SessionResolverSuite) TestViewerQuery() {
	u, _ := CreateDummyUser(s.T(), s.AppCtx.EntClient)

	q := `query { viewer { user { id email firstName lastName phoneNumber } } }`
	var resp struct {
		Viewer struct {
			User *userObj `json:"user"`
		} `json:"viewer"`
	}

	s.GqlGenClient.MustPost(q, &resp, AddContextViewerForTesting(u))
	assert.NotNil(s.T(), resp.Viewer.User)
	assert.Equal(s.T(), u.Email, resp.Viewer.User.Email)
	assert.Equal(s.T(), u.FirstName, resp.Viewer.User.FirstName)
	assert.Equal(s.T(), u.LastName, resp.Viewer.User.LastName)
	assert.Equal(s.T(), u.PhoneNumber, resp.Viewer.User.PhoneNumber)
}

func (s *SessionResolverSuite) TestCreateSessionMutation() {
	u, pw := CreateDummyUser(s.T(), s.AppCtx.EntClient)
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

	s.GqlGenClient.MustPost(q, &resp, client.Var("input", i), AddContextViewerForTesting(nil))

	assert.NotEmpty(s.T(), resp.CreateSession.Token)
	assert.NotEmpty(s.T(), resp.CreateSession.Expiry)
}

func (s *SessionResolverSuite) TestDeleteSessionMutation() {
	u, _ := CreateDummyUser(s.T(), s.AppCtx.EntClient)
	var resp struct {
		DeleteSession bool `json:"deleteSession"`
	}

	err := s.GqlGenClient.Post(`mutation { deleteSession }`, &resp, func(bd *client.Request) {}, AddContextViewerForTesting(u))
	if err != nil {
		s.T().Error(err)
	}

	assert.True(s.T(), resp.DeleteSession)
}
