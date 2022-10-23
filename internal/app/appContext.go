package app

import (
	"database/sql"
	"github.com/rs/zerolog"
)

type AppContext struct {
	DB     *sql.DB
	Logger *zerolog.Logger
}
