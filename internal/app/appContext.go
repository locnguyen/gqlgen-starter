package app

import (
	"database/sql"
	"github.com/rs/zerolog"
	"gqlgen-starter/internal/ent"
)

type AppContext struct {
	DB        *sql.DB
	EntClient *ent.Client
	Logger    *zerolog.Logger
}
