package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/middleware"
)

// Viewer is the resolver for the viewer field.
func (r *queryResolver) Viewer(ctx context.Context) (*ent.User, error) {
	u, err := middleware.GetContextUser(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("No viewer found")
	}
	return u, nil
}
