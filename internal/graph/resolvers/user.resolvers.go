package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"gqlgen-starter/internal/graph/model"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	r.Logger.Debug().Str("mutation", "createUser").Msgf("%v", input)
	return &model.User{
		ID:             "123",
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		PhoneNumber:    input.PhoneNumber,
		GraduationDate: nil,
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	r.Logger.Debug().Msg("query: user(id: ID!)")
	return &model.User{
		ID:             "123",
		FirstName:      "Natasha",
		LastName:       "Romanova",
		Email:          "blackwidow@avengers.com",
		PhoneNumber:    "+17142358092",
		GraduationDate: nil,
	}, nil
}
