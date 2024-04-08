package app

import (
	"database/sql"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"gqlgen-starter/internal/app/loaders"
	"gqlgen-starter/internal/ent"
)

type AppContext struct {
	DB             *sql.DB
	EntClient      *ent.Client
	Loaders        *loaders.Loaders
	Logger         *zerolog.Logger
	Nats           *nats.Conn
	RedisPool      *redis.Pool
	SessionManager *scs.SessionManager
}
