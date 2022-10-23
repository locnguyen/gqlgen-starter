package resolvers

import "gqlgen-starter/internal/app"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	*app.AppContext
}

func NewRootResolver(appCtx *app.AppContext) *Resolver {
	return &Resolver{appCtx}
}
