package middleware

import (
	"bytes"
	"context"
	"github.com/rs/zerolog"
	"github.com/urfave/negroni"
	"gqlgen-starter/config"
	"gqlgen-starter/internal/app"
	"io"
	"net/http"
	"time"
)

func HttpLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := negroni.NewResponseWriter(w)
		startTime := time.Now()
		ipAddress := r.RemoteAddr
		xForwardedFor := r.Header.Get("x-forwarded-for")
		if len(xForwardedFor) != 0 {
			ipAddress = xForwardedFor
		}

		log := app.GetLogger()
		ctx := context.WithValue(r.Context(), app.ContextCorrelationIDKey, ipAddress)
		r = r.WithContext(ctx)

		// To help debug the peskiest of request problems
		if config.Application.LogLevel == zerolog.LevelDebugValue {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error().
					Err(err).
					Msg("reading request body from incoming HTTP request")
			} else {
				log.Debug().
					RawJSON("body", data).
					Msgf("incoming %s request body", r.Method)
			}
			r.Body = io.NopCloser(bytes.NewReader(data))
		}

		next.ServeHTTP(w, r)
		log.Info().
			Msgf("%v %v %v %v %v %v %v", ipAddress, startTime.Format(time.UnixDate), r.Method, r.RequestURI, r.Proto, rw.Status(), time.Since(startTime))
	}
}
