package middleware

import (
	"context"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/ent"
	"net/http"
)

// AddContextViewer Check for a session cookie and if there, validate it and find a user to add to the context
func AddContextViewer(appCtx *app.AppContext, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if appCtx.SessionManager.Exists(r.Context(), ent.ContextViewerKey.Name) {
			u := appCtx.SessionManager.Get(r.Context(), ent.ContextViewerKey.Name).(*ent.User)
			r = r.WithContext(context.WithValue(r.Context(), ent.ContextViewerKey, u))
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
			return
		}
	}
}
