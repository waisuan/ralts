package testing

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
)

func InitDB() {
	m, err := migrate.New(
		"file:../../db/migrations",
		"postgres://xxxx:xxxx@localhost:5433/ralts_test?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Warn(err)
	}
}

func ClearDB() {
	m, err := migrate.New(
		"file:../../db/migrations",
		"postgres://xxxx:xxxx@localhost:5433/ralts_test?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Down(); err != nil {
		log.Warn(err)
	}
}
