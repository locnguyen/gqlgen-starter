package directives

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/rs/zerolog"
	"gqlgen-starter/internal/app/models"
	"gqlgen-starter/internal/ent"
	"gqlgen-starter/internal/gql/generated"
	"gqlgen-starter/internal/oops"
	"net/http"
)

func NewDirectiveRoot() *generated.DirectiveRoot {
	hasRole := func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []models.Role, userID string) (res interface{}, err error) {
		log := zerolog.Ctx(ctx)

		v, err := ent.GetContextViewer(ctx)
		if err != nil {
			return nil, err
		}
		u, ok := v.GetUser()
		if !ok {
			return nil, &oops.CodedError{
				HumanMessage: "Access denied to unauthenticated viewer",
				Context:      "nil user when getting from viewer context to check for role",
				HttpStatus:   http.StatusUnauthorized,
				Err:          nil,
			}
		}

		m := make(map[models.Role]bool)
		for _, r := range roles {
			m[r] = true
		}
		userHasRole := false
		for _, r := range u.Roles {
			if m[r] {
				userHasRole = true
				break
			}
		}

		if !userHasRole {
			log.Warn().
				Interface("required", roles).
				Interface("has", u.Roles).
				Str("graphqlPath", graphql.GetPath(ctx).String()).
				Msg("user missing required roles")

			return nil, &oops.CodedError{
				HumanMessage: "Required role is missing",
				Context:      fmt.Sprintf("user does not have role %v", roles),
				HttpStatus:   http.StatusUnauthorized,
				Err:          nil,
			}
		}

		return next(ctx)
	}
	return &generated.DirectiveRoot{
		HasRole: hasRole,
	}
}
