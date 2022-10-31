package middleware

import (
	"context"
	"crypto/sha256"
	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/ent/session"
	"gqlgen-starter/internal/ent/user"
	"net/http"
	"time"
)

const authCookieName = "auth-cookie"

type contextKey struct {
	name string
}

var CookieCtxKey = &contextKey{authCookieName}

type CtxAuthTokenIDKey string

const CtxAuthTokenID CtxAuthTokenIDKey = "requestAuthTokenID"

type ContextCookie struct {
	Writer http.ResponseWriter
	User   *ent.User
	Token  [32]byte
}

func (this *ContextCookie) SetSession(ses *ent.Session) {
	http.SetCookie(this.Writer, &http.Cookie{
		Name:     authCookieName,
		Value:    ses.Sid,
		HttpOnly: true,
		Path:     "/",
		Expires:  ses.Expiry,
	})
}

func (this *ContextCookie) RemoveSession() {
	http.SetCookie(this.Writer, &http.Cookie{
		Name:    authCookieName,
		MaxAge:  -1,
		Path:    "/",
		Expires: time.Now().Add(24 * -7 * time.Hour),
	})
}

// Middleware https://gqlgen.com/recipes/authentication/
func AuthCookie(env *app.AppContext, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Give all resolvers access to the cookie
		ctxCookie := &ContextCookie{
			Writer: w,
		}
		r = r.WithContext(context.WithValue(r.Context(), CookieCtxKey, ctxCookie))

		c, err := r.Cookie(authCookieName)

		// Allow unauthenticated users in because we are just looking for a cookie to find a user
		if err != nil && c == nil {
			next.ServeHTTP(w, r)
			return
		}

		u, s, err := validateAndGetUser(c, env.EntClient)

		if err != nil {
			log.Error().Err(err).Msg("Error with validateAndGetUser")
			http.SetCookie(w, &http.Cookie{
				Name:    authCookieName,
				MaxAge:  -1,
				Path:    "/",
				Expires: time.Now().Add(24 * -7 * time.Hour),
			})
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctxCookie.User = u
		ctxCookie.Token = sha256.Sum256([]byte(c.Value))

		r = r.WithContext(context.WithValue(r.Context(), CtxAuthTokenID, s.ID))
		next.ServeHTTP(w, r)
	}
}

func GetContextUser(ctx context.Context) (*ent.User, error) {
	ctxCookie := ctx.Value(CookieCtxKey).(*ContextCookie)

	if ctxCookie.User == nil {
		log.Warn().
			Msg("ctxCookie.User is nil")

		return nil, &gqlerror.Error{
			Path:    graphql.GetPath(ctx),
			Message: "Session expired",
		}
	}
	return ctxCookie.User, nil
}

func validateAndGetUser(cookie *http.Cookie, client *ent.Client) (*ent.User, *ent.Session, error) {
	t, err := client.Session.Query().
		Where(
			session.And(
				session.SidEQ(cookie.Value),
				session.TypeEQ(session.TypeGeneral),
			),
		).
		Only(context.Background())

	if err != nil {
		log.Error().
			Err(err).
			Str("sid", cookie.Value).
			Msg("Error finding session while trying to validate and get user")
		return nil, nil, err
	}

	now := time.Now()
	if t.Expiry.Before(now) {
		log.Info().
			Str("sid", cookie.Value).
			Str("sessionExpiry", t.Expiry.String()).
			Str("now", now.String()).
			Msg("Session has expired")

		return nil, nil, errors.New("Invalid session")
	}

	u, err := client.User.Query().
		Where(user.ID(t.UserID)).
		Only(context.Background())

	if err != nil {
		log.Error().
			Err(err).
			Msg("Error getting user after validating session")
		return nil, nil, err
	}

	return u, t, nil
}
