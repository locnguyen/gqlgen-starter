package extensions

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"gqlgen-starter/internal/ent"
)

// EntClientContextInjector for graphql queries.
type EntClientContextInjector struct{ Entc *ent.Client }

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
} = EntClientContextInjector{}

// ExtensionName returns the extension name.
func (EntClientContextInjector) ExtensionName() string {
	return "EntClientContextInjector"
}

// Validate is called when adding an extension to the server, it allows validation against the server's schema.
func (t EntClientContextInjector) Validate(graphql.ExecutableSchema) error {
	if t.Entc == nil {
		return errors.New("EntClientContextInjector: Ent client is nil")
	}
	return nil
}

// InterceptResponse check if there is an ent client in context and add it if not there. This makes it available in all
// graphql resolvers
func (t EntClientContextInjector) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	entClient := ent.FromContext(ctx)
	if entClient == nil {
		ctx = ent.NewContext(ctx, t.Entc)
	}
	return next(ctx)
}
