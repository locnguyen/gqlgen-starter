package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"gqlgen-starter/internal/graph/generated"
)

// Mutation returns generated.MutationResolver implementation.
func (r *RootResolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *RootResolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *RootResolver }

type queryResolver struct{ *RootResolver }
