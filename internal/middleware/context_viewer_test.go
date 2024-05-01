package middleware

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/stretchr/testify/assert"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/ent"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAddContextViewer(t *testing.T) {
	sm := scs.New()
	sm.Lifetime = 24 * time.Hour
	sm.Store = memstore.New()
	gob.Register(&ent.User{})
	appCtx := &app.AppContext{SessionManager: sm}
	userID := int64(9999)

	// need to add a fake user before AddContextViewer is executed
	addUserToSession := func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			sm.Put(r.Context(), ent.ContextViewerKey.Name, &ent.User{ID: userID})
			next.ServeHTTP(w, r)
		}
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.Context().Value(ent.ContextViewerKey)
		assert.NotNil(t, v, "viewer from context should not be nil")
		u := v.(*ent.User)
		assert.Equalf(t, userID, u.ID, "userID should be equal (%d)", u.ID)
	})

	handlerToTest := sm.LoadAndSave(addUserToSession(AddContextViewer(appCtx, nextHandler)))
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
}
