package testsupport

import (
	"context"
	"github.com/go-faker/faker/v4"
	"golang.org/x/crypto/bcrypt"
	"gqlgen-starter/internal/app/models"
	"gqlgen-starter/internal/ent"
	"testing"
)

func CreateDummyUser(t *testing.T, client *ent.Client) (*ent.User, string) {
	pw := faker.Password()
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		t.Error(err)
	}
	u, err := client.User.Create().
		SetEmail(faker.Email()).
		SetFirstName(faker.FirstName()).
		SetLastName(faker.LastName()).
		SetHashedPassword(hashedPw).
		SetPhoneNumber(faker.E164PhoneNumber()).
		SetRoles([]models.Role{models.RoleGenPop}).
		Save(context.Background())

	if err != nil {
		t.Error(err)
	}
	return u, pw
}
