package middleware

import (
	"context"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"gqlgen-starter/internal/app"
	"net/http"
)

const correlationIDKey = "correlation_id"

// CorrelationID Add a correlation ID to all requests to make tracing through logs easier
func CorrelationID(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := app.GetLogger()
		id := ulid.Make().String()
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str(correlationIDKey, id)
		})
		w.Header().Add("x-correlation-id", id)
		r = r.WithContext(log.WithContext(r.Context()))
		r = r.WithContext(context.WithValue(r.Context(), app.ContextCorrelationIDKey, id))
		next.ServeHTTP(w, r)
	}
}
