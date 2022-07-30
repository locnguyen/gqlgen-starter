package resolvers

import (
	"context"
	"gqlgen-starter/internal/graph/model"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	return &model.Todo{
		ID:   "1",
		Text: "A Todo For You",
		Done: false,
		User: nil,
	}, nil
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	todos := make([]*model.Todo, 2)
	todos[0] = &model.Todo{
		ID:   "1",
		Text: "A 1st Todo For You",
		Done: false,
		User: nil,
	}

	todos[1] = &model.Todo{
		ID:   "2",
		Text: "A 2nd Todo For You",
		Done: false,
		User: nil,
	}

	return todos, nil
}
