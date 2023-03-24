package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type CoreDatabaseInterface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type RaltsDatabase struct {
	Conn *pgx.Conn
}

func NewRaltsDatabase(connString string) *RaltsDatabase {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}

	return &RaltsDatabase{
		Conn: conn,
	}
}

func (db *RaltsDatabase) Close() {
	db.Conn.Close(context.Background())
}

func (db *RaltsDatabase) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return db.Conn.Query(ctx, sql, args...)
}

func (db *RaltsDatabase) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return db.Conn.QueryRow(ctx, sql, args...)
}
