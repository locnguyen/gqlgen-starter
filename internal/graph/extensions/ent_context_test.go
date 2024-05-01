package extensions

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gqlgen-starter/internal/ent"
	"testing"
)

func TestEntClientContextInjector_InterceptResponse(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		assert.NoError(t, err)
	}
	defer db.Close()
	driver := sql.OpenDB("postgres", db)
	driverOption := ent.Driver(driver)
	entClient := ent.NewClient(driverOption)

	subject := EntClientContextInjector{
		Entc: entClient,
	}
	ctx := context.Background()
	subject.InterceptResponse(ctx, func(ctx context.Context) *graphql.Response {
		assert.NotNil(t, ent.FromContext(ctx), "interceptor should put ent client in context")
		return &graphql.Response{Data: []byte(`{"name":"test"}`)}
	})
}
