package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"golang.org/x/crypto/bcrypt"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/ent/user"
	"gqlgen-starter/internal/graph/model"
	"gqlgen-starter/internal/middleware"
)

// CreateSession is the resolver for the createSession field.
func (r *mutationResolver) CreateSession(ctx context.Context, input model.CreateSessionInput) (*ent.Session, error) {
	u, err := r.EntClient.User.Query().Where(user.Email(input.Email)).Only(ctx)
	if err != nil {
		r.Logger.Warn().Str("email", input.Email).Msg("Unable to find user to create session")
		return nil, gqlerror.Errorf("Invalid credentials")
	}
	if err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(input.Password)); err != nil {
		r.Logger.Error().Err(err).Str("email", input.Email).Msg("Unable to verify password")
		return nil, gqlerror.Errorf("Invalid credentials")
	}

	err = r.AppContext.SessionManager.RenewToken(ctx)

	if err != nil {
		r.Logger.Error().Err(err).Msg("Error renewing session token")
		return nil, gqlerror.Errorf("Error creating session")
	}

	r.AppContext.SessionManager.Put(ctx, middleware.ContextUserKey, u)

	return &ent.Session{
		Token:  r.AppContext.SessionManager.Token(ctx),
		Expiry: r.AppContext.SessionManager.Deadline(ctx),
	}, nil
}

// DeleteSession is the resolver for the deleteSession field.
func (r *mutationResolver) DeleteSession(ctx context.Context) (bool, error) {
	_, err := middleware.GetContextUser(ctx)

	if err != nil {
		r.Logger.Warn().Msg("context user not in session context")
		return false, err
	}

	if err := r.SessionManager.Destroy(ctx); err != nil {
		r.Logger.Error().Err(err).Msg("Error while destroying session")
		return false, err
	}

	return true, nil
}

// Viewer is the resolver for the viewer field.
func (r *queryResolver) Viewer(ctx context.Context) (*ent.User, error) {
	u, err := middleware.GetContextUser(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("No viewer found")
	}
	return u, nil
}
