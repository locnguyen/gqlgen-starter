package resolvers

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gqlgen-starter/internal/ent/session"
	"gqlgen-starter/internal/ent/user"
	"testing"
)

type UserResolverSuite struct {
	suite.Suite
	TestContext
}

func TestUserResolverSuite(t *testing.T) {
	suite.Run(t, new(UserResolverSuite))
}

func (suite *UserResolverSuite) SetupSuite() {
	suite.TestContext = *InitTestContext(suite.T(), "UserResolverSuite")
}

func (suite *UserResolverSuite) TearDownSuite() {
	defer suite.pgContainer.Terminate(context.Background())
	defer suite.AppCtx.DB.Close()
}

type userResp struct {
	User *struct {
		ID          string `json:"id"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`
	}
}

var getUserQuery = `
	query GetUser($id: ID!) {
		user(id: $id) {
			id
			firstName
			lastName
			phoneNumber
			email
		}
	}
`

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

	var resp userResp
	suite.GqlGenClient.MustPost(getUserQuery, &resp, client.Var("id", u.ID))

	assert.Equal(suite.T(), resp.User.ID, fmt.Sprint(u.ID))
	assert.Equal(suite.T(), resp.User.FirstName, u.FirstName)
	assert.Equal(suite.T(), resp.User.LastName, u.LastName)
	assert.Equal(suite.T(), resp.User.PhoneNumber, u.PhoneNumber)
	assert.Equal(suite.T(), resp.User.Email, u.Email)
}

func (suite *UserResolverSuite) TestUserQueryNotFound() {
	rands, err := faker.RandomInt(100, 1000)
	if err != nil {
		suite.T().Error(err)
	}
	var resp userResp
	err = suite.GqlGenClient.Post(getUserQuery, &resp, client.Var("id", rands[0]))
	assert.ErrorContains(suite.T(), err, "User not found")
	assert.Nil(suite.T(), resp.User)
}

func (suite *UserResolverSuite) TestCreateUserMutation() {
	var resp struct {
		CreateUser struct {
			Sid    string `json:"sid"`
			Expiry string `json:"expiry"`
		}
	}
	suite.GqlGenClient.MustPost(`mutation { createUser(input: { firstName: "Natasha" lastName: "Romanova" email: "blackwidow@avengers.com" phoneNumber: "+8888888888" password: "P@ssw0rd!" passwordConfirmation: "P@ssw0rd!" }) { sid expiry } }`, &resp, AddContextCookieForTesting(nil))
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
