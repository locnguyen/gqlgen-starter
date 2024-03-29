package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"golang.org/x/crypto/bcrypt"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/graph/generated"
	"gqlgen-starter/internal/graph/model"
	"gqlgen-starter/internal/middleware"
	"strconv"
	"strings"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*ent.Session, error) {
	if input.Password != input.PasswordConfirmation {
		return nil, gqlerror.Errorf("Password confirmation does not match password")
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		r.Logger.Err(err).Msg("Error generating password")
		return nil, err
	}

	u, err := r.EntClient.User.Create().
		SetEmail(strings.ToLower(input.Email)).
		SetFirstName(input.FirstName).
		SetLastName(input.LastName).
		SetPhoneNumber(input.PhoneNumber).
		SetHashedPassword(hashedPw).
		Save(ctx)

	if err != nil {
		r.Logger.Err(err).Str("email", input.Email).Msg("Error creating user in database")
		return nil, err
	}

	r.Logger.Debug().
		Str("mutation", "createUser").
		Interface("result", u).
		Msg("Created user")

	if err := r.AppContext.SessionManager.RenewToken(ctx); err != nil {
		r.Logger.Error().Err(err).Msg("Error renewing token after creating user")
		return nil, gqlerror.Errorf("Error creating session for new user")
	}
	r.AppContext.SessionManager.Put(ctx, middleware.ContextUserKey, u)

	return &ent.Session{
		Token:  r.AppContext.SessionManager.Token(ctx),
		Expiry: r.AppContext.SessionManager.Deadline(ctx),
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*ent.User, error) {
	id64, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		r.Logger.Err(err).Msgf("Error parsing id=%s", id)
		return nil, err
	}

	u, err := r.EntClient.User.Get(ctx, id64)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, gqlerror.Errorf("User not found")
		}
		r.Logger.Err(err).Msgf("Error while querying for user=%s", id)
		return nil, err
	}

	r.Logger.Debug().
		Str("query", "user(id: ID!)").
		Str("id", id).
		Interface("result", u).
		Msg("Queried for user")
	return u, nil
}

// ID is the resolver for the id field.
func (r *userResolver) ID(ctx context.Context, obj *ent.User) (string, error) {
	id := strconv.FormatInt(obj.ID, 10)
	return id, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
