package loaders

import (
	"entgo.io/ent/dialect/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"gqlgen-starter/internal/ent"
	"testing"
)

func TestHandleErrors(t *testing.T) {
	ints, _ := faker.RandomInt(1)
	loaderResults := handleError(ints[0], errors.New("handleErrors test"))
	assert.Equal(t, ints[0], len(loaderResults), "there should be same number of errors as items")
}

func TestNewLoader(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		assert.NoError(t, err)
	}
	defer db.Close()
	driver := sql.OpenDB("postgres", db)
	driverOption := ent.Driver(driver)
	entClient := ent.NewClient(driverOption)

	loaders := NewLoaders(entClient)
	assert.NotNil(t, loaders, "should return a new loader map")
}
