package resolvers

import (
	"context"
	"gqlgen-starter/internal/graph/model"
)

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
