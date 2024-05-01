package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenPostgresConn(t *testing.T) {
	ctx := context.Background()

	_, databaseURL, err := StartPgContainer(ctx, fmt.Sprintf("%s_db", "TestOpenPostgresConn"))
	if err != nil {
		t.Error(err)
	}

	dbConn, err := OpenPostgresConn(ctx, *databaseURL)
	assert.NoError(t, err, "should connect to database container")
	assert.NotNil(t, dbConn, "should have a database connection")
}
