package resolvers

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"gqlgen-starter/internal/graph/generated"
	"testing"
)

func TestUserResolvers(t *testing.T) {
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: NewRootResolver()})))

	t.Run("Users query", func(t *testing.T) {
		var resp struct {
			Users []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
		}

		c.MustPost(`{ users { id name } }`, &resp)
		assert.NotEmpty(t, resp.Users)
	})

	t.Run("createUser mutation", func(t *testing.T) {
		var resp struct {
			CreateUser struct {
				ID   string
				Name string
			}
		}

		c.MustPost(`mutation { createUser(input: { name: "Natasha Romanova" }) { id name } }`, &resp)
		assert.Equal(t, "Natasha Romanova", resp.CreateUser.Name)
	})
}
