package ent

import (
	"context"
	"gqlgen-starter/internal/oops"
	"net/http"
)

type contextViewerKey struct {
	Name string
}

var ContextViewerKey = &contextViewerKey{"contextViewerKey"}

type ContextViewer interface {
	GetUser() (*User, bool)
}

type Viewer struct {
	User *User
}

func (cv Viewer) GetUser() (*User, bool) {
	if cv.User != nil {
		return cv.User, true
	}
	return nil, false
}

func GetContextViewer(ctx context.Context) (ContextViewer, error) {
	v := ctx.Value(ContextViewerKey)

	if v == nil {
		return nil, &oops.CodedError{
			HumanMessage: "Your session has expired, please sign-in again",
			Context:      "getting context viewer but it was nil",
			HttpStatus:   http.StatusUnauthorized,
		}
	}
	return v.(ContextViewer), nil
}
