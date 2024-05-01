package resolvers

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SchemaResolverSuite struct {
	suite.Suite
	TestContext
}

func TestSchemaResolverSuite(t *testing.T) {
	suite.Run(t, &SchemaResolverSuite{})
}

func (s *SchemaResolverSuite) SetupSuite() {
	s.TestContext = *InitTestContext(s.T(), "SchemaResolverSuite")
}

func (s *SchemaResolverSuite) TearDownSuite() {
	defer s.PgContainer.Terminate(context.Background())
	defer s.AppCtx.DB.Close()
}

func (s *SessionResolverSuite) TestHealthCheckQuery() {
	q := `query { healthCheck }`
	var resp struct {
		HealthCheck string
	}

	s.GqlGenClient.MustPost(q, &resp, AddContextViewerForTesting(nil))
	assert.NotNil(s.T(), resp.HealthCheck)
}

func (s *SessionResolverSuite) TestHelloMutation() {
	q := `mutation { hello }`
	var resp struct {
		Hello string
	}

	s.GqlGenClient.MustPost(q, &resp, AddContextViewerForTesting(nil))
	assert.Equal(s.T(), "Hello world", resp.Hello)
}
