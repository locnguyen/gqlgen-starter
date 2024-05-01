package middleware

import (
	"github.com/stretchr/testify/assert"
	"gqlgen-starter/internal/app"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorrelationID(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(app.ContextCorrelationIDKey)
		assert.NotEmpty(t, id, "correlation ID should not be empty")
	})

	handlerToTest := CorrelationID(nextHandler)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}
