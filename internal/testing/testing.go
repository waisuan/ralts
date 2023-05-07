package testing

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"ralts/internal/config"
)

var ctx = context.Background()

func TestHelper(cfg *config.Config) func() {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisConn,
		DB:   0, // use default DB
	})

	rdb.FlushAll(ctx)
	ClearDB(cfg)
	InitDB(cfg)

	return func() {
		rdb.FlushAll(ctx)
		ClearDB(cfg)
	}
}

func InitDB(cfg *config.Config) {
	m, err := migrate.New(
		"file:../../db/migrations",
		cfg.DatabaseConn)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Warn(err)
	}
}

func ClearDB(cfg *config.Config) {
	m, err := migrate.New(
		"file:../../db/migrations",
		cfg.DatabaseConn)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Down(); err != nil {
		log.Warn(err)
	}
}
