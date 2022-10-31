package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/graph/generated"
	"gqlgen-starter/internal/graph/model"
	"gqlgen-starter/internal/middleware"
	"strconv"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*ent.Post, error) {
	u, err := middleware.GetContextUser(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("Invalid session")
	}

	p, err := r.EntClient.Post.Create().SetAuthor(u).SetContent(input.Content).Save(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("Error saving post")
	}
	return p, nil
}

// ID is the resolver for the id field.
func (r *postResolver) ID(ctx context.Context, obj *ent.Post) (string, error) {
	id := strconv.FormatInt(obj.ID, 10)
	return id, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*ent.Post, error) {
	id64, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		r.Logger.Err(err).Msgf("Error parsing id=%s", id)
		return nil, err
	}
	r.Logger.Debug().Msgf("Try to query for post=%s", id)
	p, err := r.EntClient.Post.Get(ctx, id64)

	if err != nil {
		if ent.IsNotFound(err) {
			r.Logger.Debug().Msgf("Post %s not found", id)
			return nil, gqlerror.Errorf("Post not found")
		}
		r.Logger.Err(err).Msgf("Error while querying for post=%s", id)
		return nil, gqlerror.Errorf("Error querying for post")
	}

	return p, nil
}

// Post returns generated.PostResolver implementation.
func (r *Resolver) Post() generated.PostResolver { return &postResolver{r} }

type postResolver struct{ *Resolver }
