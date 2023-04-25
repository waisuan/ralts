package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"ralts/internal/config"
)

type CoreDatabaseInterface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type RaltsDatabase struct {
	Conn *pgxpool.Pool
}

func NewRaltsDatabase(cfg *config.Config) *RaltsDatabase {
	conn, err := pgxpool.Connect(context.Background(), cfg.DatabaseConn)
	if err != nil {
		log.Fatal(err)
	}

	return &RaltsDatabase{
		Conn: conn,
	}
}

func (db *RaltsDatabase) Close() {
	db.Conn.Close()
}

func (db *RaltsDatabase) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return db.Conn.Query(ctx, sql, args...)
}

func (db *RaltsDatabase) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return db.Conn.QueryRow(ctx, sql, args...)
}
