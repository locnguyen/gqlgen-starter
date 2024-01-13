package db

import (
	"context"
	"database/sql"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gqlgen-starter/config"
)

func OpenPostgresConn(ctx context.Context, databaseURL string) (*sql.DB, error) {
	log := zerolog.Ctx(ctx)
	log.Info().
		Msg("opening database connection...")

	conn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("could not connect to Postgres")
		return nil, err
	}

	log.Info().
		Msg("pinging Postgres...")
	if err := conn.Ping(); err != nil {
		log.Error().
			Err(err).
			Msg("🚨️   could not ping Postgres after opening connection")
		return nil, err
	}
	log.Info().Msg("☎️️   Postgres ping succeeded!")

	return conn, nil
}

func NewRedisPool(ctx context.Context, cfg config.ApplicationConfig) (*redis.Pool, error) {
	log := zerolog.Ctx(ctx)
	if cfg.RedisURL == "" {
		log.Error().
			Msg("🚨️   environment variable REDIS_URL not set")
		return nil, errors.New("REDIS_URL not set")
	}

	log.Info().
		Msg("🏗️  creating Redis connection pool")
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.RedisURL, redis.DialPassword(cfg.RedisAuth))
			if err != nil {
				log.Error().
					Err(err).
					Msgf("🚨️   error connecting to Redis from pool: %s", cfg.RedisURL)
				panic(err.Error())
			}
			log.Info().
				Msgf("☎️   connected to Redis from pool: %s", cfg.RedisURL)
			return c, err
		},
	}, nil
}
