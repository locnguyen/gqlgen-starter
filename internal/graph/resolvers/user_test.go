package resolvers

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gqlgen-starter/internal/app"
	"gqlgen-starter/internal/graph/generated"
	"os"
	"testing"
)

func TestUserResolvers(t *testing.T) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	appCtx := &app.AppContext{
		DB:     nil,
		Logger: &logger,
	}
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: NewRootResolver(appCtx)})))

	t.Run("user(id: ID!) query", func(t *testing.T) {
		var resp struct {
			User struct {
				ID        string `json:"id"`
				FirstName string `json:"firstName"`
			}
		}

		c.MustPost(`{ user(id: "123") { id firstName } }`, &resp)
		assert.NotEmpty(t, resp.User)
	})

	t.Run("createUser mutation", func(t *testing.T) {
		var resp struct {
			CreateUser struct {
				ID          string `json:"id"`
				FirstName   string `json:"firstName"`
				LastName    string `json:"lastName"`
				Email       string `json:"email"`
				PhoneNumber string `json:"phoneNumber"`
			}
		}

		c.MustPost(`mutation { createUser(input: { firstName: "Natasha" lastName: "Romanova" email: "blackwidow@avengers.com" phoneNumber: "+17142358092" password: "P@ssw0rd!" passwordConfirmation: "P@ssw0rd!" }) { id firstName } }`, &resp)
		assert.Equal(t, "Natasha", resp.CreateUser.FirstName)
	})
}
