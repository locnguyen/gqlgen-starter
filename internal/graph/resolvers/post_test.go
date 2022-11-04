package resolvers

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PostResolverSuite struct {
	suite.Suite
	TestContext
}

func TestPostResolvers(t *testing.T) {
	suite.Run(t, new(PostResolverSuite))
}

func (suite *PostResolverSuite) SetupSuite() {
	suite.TestContext = *InitTestContext(suite.T(), "PostResolverSuite")
}

func (suite *PostResolverSuite) TeardownSuite() {
	defer suite.pgContainer.Terminate(context.Background())
	defer suite.AppCtx.DB.Close()
}

var getPostQuery = `
		query ($id: ID!) {
			post(id: $id) {
				id
				content
				createTime
				updateTime
				author {
					id
					email
				}
			}
		}
	`

type PostObj struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
	Author     struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"author"`
}

type postResp struct {
	Post   *PostObj `json:"post"`
	Errors *[]struct {
		Message string `json:"message"`
		Path    string `json:"path"`
	} `json:"errors"`
}

func (suite *PostResolverSuite) TestGetPostQuery() {
	u, _ := CreateDummyUser(suite.T(), *suite.AppCtx.EntClient)

	p, err := suite.AppCtx.EntClient.Post.Create().
		SetAuthor(u).
		SetContent(faker.Paragraph()).
		Save(context.Background())

	if err != nil {
		suite.T().Error(err)
	}

	var resp postResp

	suite.GqlGenClient.MustPost(getPostQuery, &resp, client.Var("id", p.ID))
	assert.Equal(suite.T(), resp.Post.ID, fmt.Sprint(p.ID))
	assert.Equal(suite.T(), resp.Post.Content, p.Content)
	assert.NotEmpty(suite.T(), resp.Post.CreateTime)
	assert.NotEmpty(suite.T(), resp.Post.UpdateTime)

	assert.NotEmpty(suite.T(), resp.Post.Author.ID)
	assert.Equal(suite.T(), resp.Post.Author.Email, u.Email)
}

func (suite *PostResolverSuite) TestGetPostNotFoundQuery() {
	rands, err := faker.RandomInt(100, 1000)
	if err != nil {
		suite.T().Error(err)
	}
	var resp postResp
	err = suite.GqlGenClient.Post(getPostQuery, &resp, client.Var("id", rands[0]))
	assert.ErrorContains(suite.T(), err, "Post not found")
	assert.Nil(suite.T(), resp.Post)
}

func (suite *PostResolverSuite) TestCreatePostMutation() {
	var q = `
		mutation ($input: CreatePostInput!) {
			createPost(input: $input) {
				id
				content
				createTime
				updateTime
			}
		}
	`

	var resp struct {
		CreatePost PostObj `json:"createPost"`
	}

	author, _ := CreateDummyUser(suite.T(), *suite.AppCtx.EntClient)
	content := faker.Paragraph()
	type input struct {
		Content string `json:"content"`
	}
	i := &input{
		Content: content,
	}

	suite.GqlGenClient.MustPost(q, &resp, client.Var("input", i), AddContextCookieForTesting(author, nil))

	assert.NotEmpty(suite.T(), resp.CreatePost.ID)
	assert.Equal(suite.T(), content, resp.CreatePost.Content)
	assert.NotEmpty(suite.T(), resp.CreatePost.CreateTime)
	assert.NotEmpty(suite.T(), resp.CreatePost.UpdateTime)
}
