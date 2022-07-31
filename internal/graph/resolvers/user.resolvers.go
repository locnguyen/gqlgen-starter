package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"gqlgen-starter/internal/graph/model"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	return &model.User{
		ID:   "1",
		Name: "Natasha Romanova",
	}, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	users := make([]*model.User, 2)
	users[0] = &model.User{
		ID:   "1",
		Name: "Natasha Romanova",
	}

	users[1] = &model.User{
		ID:   "2",
		Name: "Steve Rogers",
	}

	return users, nil
}
