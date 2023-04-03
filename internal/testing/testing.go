package testing

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"ralts/internal/config"
)

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
