package resolvers

import (
	"context"
	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/gql/generated"
	"gqlgen-starter/internal/graph/extensions"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	*app.AppContext
}

func NewRootResolver(appCtx *app.AppContext) *Resolver {
	return &Resolver{appCtx}
}

// CreateGqlServer creates a server that understands GraphQL queries. It loads the entire schema from the root resolver.
func CreateGqlServer(appCtx *app.AppContext, schemaConfig *generated.Config) *handler.Server {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(*schemaConfig))

	// Enable mutations to be wrapped in a transaction, so we don't have to manually create and commit transactions
	// https://entgo.io/docs/graphql/#transactional-mutations
	srv.Use(entgql.Transactioner{TxOpener: appCtx.EntClient})

	// Catch any case when the ent client is not in the context and add it (for gql queries).  gql mutations should be
	// covered by entgql.Transactioner above
	srv.Use(extensions.EntClientContextInjector{Entc: appCtx.EntClient})

	// A centralized place where we can process raw error instances and return something useful to API clients, log the
	// errors, and send alerts when they occur.
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		presentedErr := graphql.DefaultErrorPresenter(ctx, e)
		return presentedErr
	})
	return srv
}
