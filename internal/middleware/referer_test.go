package middleware

import (
	"github.com/stretchr/testify/assert"
	"gqlgen-starter/config"
	"gqlgen-starter/internal/app"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRefererServer_NoReferer(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ref := r.Context().Value(app.RefererServerKey)
		assert.NotEmpty(t, ref, "referer header should not be empty")
		assert.Equal(t, config.Application.DefaultReferrerURL, ref.(string))
	})

	handlerToTest := RefererServer(nextHandler)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}

func TestRefererServer_AllowListed(t *testing.T) {
	nextHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ref := r.Context().Value(app.RefererServerKey)
			assert.NotEmpty(t, ref, "referer header should not be empty")
			assert.Equal(t, "http://localhost", ref.(string))
		}
	}

	handlerToTest := RefererServer(nextHandler())
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add("referer", "http://localhost")
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}

func TestRefererServer_NowAllowListed(t *testing.T) {
	nextHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ref := r.Context().Value(app.RefererServerKey)
			assert.NotEmpty(t, ref, "referer header should not be empty")
			assert.Equal(t, "http://localhost:9000", ref.(string))
		}
	}

	handlerToTest := RefererServer(nextHandler())
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add("referer", "https://unittest.com")
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}
