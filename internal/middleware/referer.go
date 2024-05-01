package middleware

import (
	"context"
	"fmt"
	"golang.org/x/exp/slices"
	"gqlgen-starter/config"
	"gqlgen-starter/internal/app"
	"net/http"
	"net/url"
	"strings"
)

func getAllowedReferrers() []string {
	if config.Application.GoEnv == "production" {
		return []string{}
	} else {
		return []string{
			"localhost",
		}
	}

}

// RefererServer Sometimes we need to know the referrer to build URLs with the correct hostname. For example,
// magic links and redirect URLs.
func RefererServer(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := r.Header.Get("referer")

		if len(ref) == 0 {
			r = r.WithContext(context.WithValue(r.Context(), app.RefererServerKey, config.Application.DefaultReferrerURL))
		} else {
			refUrl, err := url.Parse(ref)
			if err != nil {
				r = r.WithContext(context.WithValue(r.Context(), app.RefererServerKey, config.Application.DefaultReferrerURL))
			}
			refererHostIsAllowed := slices.IndexFunc(getAllowedReferrers(), func(s string) bool {
				parts := strings.Split(refUrl.Hostname(), ".")
				var hostOrDomain string
				if len(parts) > 1 {
					hostOrDomain = parts[len(parts)-2] + "." + parts[len(parts)-1]
				} else {
					hostOrDomain = parts[0]
				}
				return s == hostOrDomain
			}) >= 0

			if refererHostIsAllowed {
				refererServer := fmt.Sprintf("%s://%s", refUrl.Scheme, refUrl.Hostname())
				if len(refUrl.Port()) > 0 {
					refererServer = refererServer + ":" + refUrl.Port()
				}
				r = r.WithContext(context.WithValue(r.Context(), app.RefererServerKey, refererServer))
			} else {
				r = r.WithContext(context.WithValue(r.Context(), app.RefererServerKey, config.Application.DefaultReferrerURL))
			}
		}
		next.ServeHTTP(w, r)
	}
}
