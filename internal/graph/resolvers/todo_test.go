package resolvers

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"gqlgen-starter/internal/graph/generated"
	"testing"
)

func TestTodoResolver(t *testing.T) {
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: NewRootResolver()})))

	t.Run("Todos query", func(t *testing.T) {
		var resp struct {
			Todos []struct {
				ID   string `json:"id"`
				Text string `json:"text"`
				Done bool   `json:"done"`
			}
		}

		c.MustPost(`query { todos { id text done } }`, &resp)
		assert.NotEmpty(t, resp.Todos)
	})

	t.Run("createTodo mutation", func(t *testing.T) {
		var resp struct {
			CreateTodo struct {
				ID   string
				Text string
			}
		}

		c.MustPost(`mutation { createTodo(input: { text: "asdf" userId: "1" }) { id } }`, &resp)
		assert.Equal(t, "1", resp.CreateTodo.ID)
	})
}
