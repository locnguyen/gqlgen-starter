package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/ent/session"
	"gqlgen-starter/internal/ent/user"
	"gqlgen-starter/internal/graph/model"
	"gqlgen-starter/internal/middleware"
	"time"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"golang.org/x/crypto/bcrypt"
)

// CreateSession is the resolver for the createSession field.
func (r *mutationResolver) CreateSession(ctx context.Context, input model.CreateSessionInput) (*ent.Session, error) {
	u, err := r.EntClient.User.Query().Where(user.Email(input.Email)).Only(ctx)
	if err != nil {
		r.Logger.Warn().Str("email", input.Email).Msg("Unable to find user to create session")
		return nil, gqlerror.Errorf("Invalid credentials")
	}
	r.Logger.Debug().Interface("HashedPassword", u.HashedPassword).Msg("user's hashed password")
	if err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(input.Password)); err != nil {
		r.Logger.Error().Err(err).Str("email", input.Email).Msg("Unable to verify password")
		return nil, gqlerror.Errorf("Invalid credentials")
	}

	randomBytes := make([]byte, 16)
	_, err = rand.Read(randomBytes)
	if err != nil {
		r.Logger.Err(err).Msg("Error creating session")
		return nil, err
	}

	sid := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	sess, err := r.EntClient.Session.Create().
		SetSid(sid).
		SetExpiry(time.Now().Add(24 * time.Hour)).
		SetType(session.TypeGeneral).
		SetUser(u).
		Save(ctx)

	r.Logger.Debug().Interface("session", sess).Msg("Created session")
	ctxCookie := ctx.Value(middleware.CookieCtxKey).(*middleware.ContextCookie)
	ctxCookie.SetSession(sess)

	return sess, nil
}

// DeleteSession is the resolver for the deleteSession field.
func (r *mutationResolver) DeleteSession(ctx context.Context) (*ent.Session, error) {
	_, err := middleware.GetContextUser(ctx)

	if err != nil {
		r.Logger.Warn().Msg("context user not in session context")
		return nil, err
	}

	ctxCookie := ctx.Value(middleware.CookieCtxKey).(*middleware.ContextCookie)
	ctxCookie.RemoveSession()

	sess, err := r.EntClient.Session.Query().
		Where(session.And(session.Sid(*ctxCookie.Sid), session.Deleted(false))).
		Only(ctx)

	if err != nil {
		r.Logger.Error().Err(err).Str("sid", *ctxCookie.Sid).Msg("Error querying for session")
		return nil, gqlerror.Errorf("Error finding session to delete")
	}

	sess, err = r.EntClient.Session.UpdateOne(sess).SetDeleted(true).Save(ctx)

	if err != nil {
		r.Logger.Error().Err(err).Str("sid", *ctxCookie.Sid).Msg("Error updating deleted=true for session")
		return nil, gqlerror.Errorf("Error deleting session")
	}
	r.Logger.Debug().Int("sessionId", sess.ID).Msg("Updated session to be deleted")

	return sess, nil
}

// Viewer is the resolver for the viewer field.
func (r *queryResolver) Viewer(ctx context.Context) (*ent.User, error) {
	u, err := middleware.GetContextUser(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("No viewer found")
	}
	return u, nil
}
