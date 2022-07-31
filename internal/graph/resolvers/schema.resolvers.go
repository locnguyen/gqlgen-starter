package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"gqlgen-starter/internal/graph/generated"
)

// Hello is the resolver for the hello field.
func (r *mutationResolver) Hello(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Test is the resolver for the test field.
func (r *queryResolver) Test(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
