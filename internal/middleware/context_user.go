package middleware

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/ent"
	"net/http"
)

const SessionCookieName = "session-cookie"
const ContextUserKey = "ContextUser"

// AddContextUser Check for a session cookie and if there, validate it and find a user to add to the context
func AddContextUser(appCtx *app.AppContext, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if appCtx.SessionManager.Exists(r.Context(), ContextUserKey) {
			u := appCtx.SessionManager.Get(r.Context(), ContextUserKey).(*ent.User)
			r = r.WithContext(context.WithValue(r.Context(), ContextUserKey, u))
			next.ServeHTTP(w, r)
		} else {
			// Allow the request to pass through
			next.ServeHTTP(w, r)
			return
		}
	}
}

func GetContextUser(ctx context.Context) (*ent.User, error) {
	u := ctx.Value(ContextUserKey)

	if u == nil {
		return nil, &gqlerror.Error{
			Path:    graphql.GetPath(ctx),
			Message: "Session expired",
		}
	}
	return u.(*ent.User), nil
}
