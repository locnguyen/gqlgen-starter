package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/gql/generated"
	"gqlgen-starter/internal/gql/model"
	"gqlgen-starter/internal/oops"
	"gqlgen-starter/internal/utils"
	"net/http"
	"strconv"

	"github.com/graph-gophers/dataloader"
	"github.com/rs/zerolog"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*ent.Post, error) {
	entC := ent.FromContext(ctx)

	v, err := ent.GetContextViewer(ctx)
	if err != nil {
		return nil, err
	}
	u, _ := v.GetUser()

	p, err := entC.Post.Create().
		SetAuthor(u).
		SetContent(input.Content).
		Save(ctx)

	if err != nil {
		return nil, &oops.CodedError{
			HumanMessage: "Error saving post",
			Context:      "inserting new post",
			HttpStatus:   http.StatusInternalServerError,
			Err:          err,
		}
	}
	return p, nil
}

// ID is the resolver for the id field.
func (r *postResolver) ID(ctx context.Context, obj *ent.Post) (string, error) {
	id := strconv.FormatInt(obj.ID, 10)
	return id, nil
}

// Author is the resolver for the author field.
func (r *postResolver) Author(ctx context.Context, obj *ent.Post) (*ent.User, error) {
	log := zerolog.Ctx(ctx)
	_, err := ent.GetContextViewer(ctx)
	if err != nil {
		return nil, err
	}

	thunk := r.AppContext.Loaders.UserLoader.Load(ctx, dataloader.StringKey(strconv.FormatInt(obj.AuthorID, 10)))
	result, err := thunk()
	if err != nil {
		log.Error().
			Err(err).
			Int64("AuthorID", obj.AuthorID).
			Msg("getting Author for Post via thunk")
		return nil, &oops.CodedError{
			HumanMessage: "Error getting post author",
			Context:      "executing thunk to get case participant",
			HttpStatus:   http.StatusInternalServerError,
			Err:          err,
		}
	}

	if result != nil {
		return result.(*ent.User), nil
	}
	return nil, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*ent.Post, error) {
	entC := ent.FromContext(ctx)

	id64, err := utils.ID64(id)
	if err != nil {
		return nil, err
	}

	p, err := entC.Post.Get(ctx, id64)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &oops.CodedError{
				HumanMessage: "Post not found",
				Context:      fmt.Sprintf("querying for Post %s", id),
				HttpStatus:   http.StatusNotFound,
				Err:          err,
			}
		}
		return nil, &oops.CodedError{
			HumanMessage: "Error loading post data",
			Context:      fmt.Sprintf("querying for post %s", id),
			HttpStatus:   http.StatusInternalServerError,
			Err:          err,
		}
	}

	return p, nil
}

// Post returns generated1.PostResolver implementation.
func (r *Resolver) Post() generated.PostResolver { return &postResolver{r} }

type postResolver struct{ *Resolver }
