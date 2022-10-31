package middleware

import (
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gqlgen-starter/internal/ent"
	"testing"
	"time"
)

type CookieAuthSuite struct {
	suite.Suite
}

func TestAuthCookie(t *testing.T) {
	suite.Run(t, new(CookieAuthSuite))
}

func (suite *CookieAuthSuite) TestGetContextUserExpired() {
	ctx := context.WithValue(context.Background(), CookieCtxKey, &ContextCookie{User: nil})
	_, err := GetContextUser(ctx)
	assert.ErrorContains(suite.T(), err, "Session expired")
}

func (suite *CookieAuthSuite) TestGetContextUser() {
	randInts, err := faker.RandomInt(1, 1000)
	if err != nil {
		suite.T().Error(err)
	}

	subject := &ent.User{
		ID:             int64(randInts[0]),
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		Email:          faker.Email(),
		HashedPassword: []byte(faker.Password()),
		FirstName:      faker.FirstName(),
		LastName:       faker.LastName(),
		PhoneNumber:    faker.Phonenumber(),
	}

	ctx := context.WithValue(context.Background(), CookieCtxKey, &ContextCookie{User: subject})

	u, err := GetContextUser(ctx)

	if err != nil {
		suite.T().Error(err)
	}

	assert.EqualValues(suite.T(), subject, u)
}
