package loaders

import (
	"context"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"errors"
	"github.com/go-faker/faker/v4"
	"github.com/graph-gophers/dataloader"
	"github.com/stretchr/testify/assert"
	"gqlgen-starter/db"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/oops"
	"gqlgen-starter/internal/testsupport"
	"net/http"
	"strconv"
	"testing"
)

func TestUserReader_GetUsersBatchFn(t *testing.T) {
	ctx := context.Background()
	_, connStr, err := db.StartPgContainer(ctx, "TestUserReader_GetUsersBatchFn")
	assert.NoError(t, err)

	dbConn, err := db.OpenPostgresConn(ctx, *connStr)
	assert.NoError(t, err)
	entC := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, dbConn)))
	if err := entC.Schema.Create(ctx); err != nil {
		t.Fatal(err)
	}

	ctx = ent.NewContext(ctx, entC)
	randomInts, err := faker.RandomInt(20)
	assert.NoError(t, err)

	users := []*ent.User{}

	for i := 0; i < randomInts[0]; i++ {
		u, _ := testsupport.CreateDummyUser(t, entC)
		assert.NotNil(t, u)
		users = append(users, u)
	}

	ur := userReader{
		EntClient: entC,
	}

	idKeys := make([]string, 0)
	for _, u := range users {
		idKeys = append(idKeys, strconv.FormatInt(u.ID, 10))
	}

	results := ur.GetUsersBatchFn(ctx, dataloader.NewKeysFromStrings(idKeys))

	assert.Equal(t, len(users), len(results))
}

func TestUserReader_GetUsersBatchFnNotFound(t *testing.T) {
	ctx := context.Background()
	_, connStr, err := db.StartPgContainer(ctx, "TestUserReader_GetUsersBatchFn")
	assert.NoError(t, err)

	dbConn, err := db.OpenPostgresConn(ctx, *connStr)
	assert.NoError(t, err)
	entC := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, dbConn)))
	ur := userReader{
		EntClient: entC,
	}

	if err := entC.Schema.Create(ctx); err != nil {
		t.Fatal(err)
	}

	ctx = ent.NewContext(ctx, entC)
	results := ur.GetUsersBatchFn(ctx, dataloader.NewKeysFromStrings([]string{"789275982", "323537"}))
	assert.NotEmpty(t, 2, len(results), "should get two items back still")
	for _, r := range results {
		assert.ErrorContains(t, r.Error, "not found")
		var codedErr *oops.CodedError
		if errors.As(r.Error, &codedErr) {
			assert.Equal(t, codedErr.HttpStatus, http.StatusNotFound)
		}
	}
}
