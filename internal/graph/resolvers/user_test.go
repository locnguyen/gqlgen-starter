package resolvers

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gqlgen-starter/internal/app/models"
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

type userObj struct {
	ID          string `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

type userResp struct {
	User *userObj
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
		SetRoles([]models.Role{}).
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
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp.User)
}

func (suite *UserResolverSuite) TestCreateUserMutation() {
	var resp struct {
		CreateUser struct {
			Token  string `json:"token"`
			Expiry string `json:"expiry"`
		}
	}
	suite.GqlGenClient.MustPost(`mutation { createUser(input: { firstName: "Natasha" lastName: "Romanova" email: "blackwidow@avengers.com" phoneNumber: "+8888888888" password: "P@ssw0rd!" passwordConfirmation: "P@ssw0rd!" }) { token expiry } }`, &resp, AddContextViewerForTesting(nil))
	assert.NotEmpty(suite.T(), resp.CreateUser.Token)
	assert.NotEmpty(suite.T(), resp.CreateUser.Expiry)

	subject, err := suite.AppCtx.EntClient.User.Query().
		Where(user.Email("blackwidow@avengers.com")).
		Only(context.Background())

	if err != nil {
		suite.T().Error(err)
	}

	assert.Equal(suite.T(), "blackwidow@avengers.com", subject.Email)
}

func (suite *UserResolverSuite) TestCreateUserMutationPasswordMismatch() {
	var resp struct {
		CreateUser struct {
			Token  string `json:"token"`
			Expiry string `json:"expiry"`
		}
	}
	err := suite.GqlGenClient.Post(`mutation { createUser(input: { firstName: "Natasha" lastName: "Romanova" email: "blackwidow@avengers.com" phoneNumber: "+8888888888" password: "P@ssw0rd!" passwordConfirmation: "xxx" }) { token expiry } }`, &resp, AddContextViewerForTesting(nil))
	assert.Error(suite.T(), err)
}
