package resolvers

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

type PostResolverSuite struct {
	suite.Suite
	TestContext
}

func TestPostResolvers(t *testing.T) {
	suite.Run(t, new(PostResolverSuite))
}

func (s *PostResolverSuite) SetupSuite() {
	s.TestContext = *InitTestContext(s.T(), "PostResolverSuite")
}

func (s *PostResolverSuite) TeardownSuite() {
	defer s.PgContainer.Terminate(context.Background())
	defer s.AppCtx.DB.Close()
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

func (s *PostResolverSuite) TestGetPostQuery() {
	author, _ := CreateDummyUser(s.T(), s.AppCtx.EntClient)

	p, err := s.AppCtx.EntClient.Post.Create().
		SetAuthor(author).
		SetContent(faker.Paragraph()).
		Save(context.Background())

	if err != nil {
		s.T().Error(err)
	}

	var resp postResp

	s.GqlGenClient.MustPost(getPostQuery, &resp, client.Var("id", p.ID), AddContextViewerForTesting(author))
	assert.Equal(s.T(), resp.Post.ID, fmt.Sprint(p.ID))
	assert.Equal(s.T(), resp.Post.Content, p.Content)
	assert.Equal(s.T(), resp.Post.Author.ID, strconv.FormatInt(author.ID, 10))
	assert.Equal(s.T(), resp.Post.Author.Email, author.Email)
	assert.NotEmpty(s.T(), resp.Post.CreateTime)
	assert.NotEmpty(s.T(), resp.Post.UpdateTime)

	assert.NotEmpty(s.T(), resp.Post.Author.ID)
	assert.Equal(s.T(), resp.Post.Author.Email, author.Email)
}

func (s *PostResolverSuite) TestCreatePostMutation() {
	var q = `
		mutation ($input: CreatePostInput!) {
			createPost(input: $input) {
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

	var resp struct {
		CreatePost PostObj `json:"createPost"`
	}

	author, _ := CreateDummyUser(s.T(), s.AppCtx.EntClient)
	content := faker.Paragraph()
	type input struct {
		Content string `json:"content"`
	}
	i := &input{
		Content: content,
	}

	s.GqlGenClient.MustPost(q, &resp, client.Var("input", i), AddContextViewerForTesting(author))

	assert.NotEmpty(s.T(), resp.CreatePost.ID)
	assert.Equal(s.T(), content, resp.CreatePost.Content)
	assert.NotEmpty(s.T(), resp.CreatePost.CreateTime)
	assert.NotEmpty(s.T(), resp.CreatePost.UpdateTime)
	assert.Equal(s.T(), resp.CreatePost.Author.ID, strconv.FormatInt(author.ID, 10))
	assert.Equal(s.T(), resp.CreatePost.Author.Email, author.Email)
}
