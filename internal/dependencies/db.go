package dependencies

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"ralts/internal/config"
)

type CoreStorageInterface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}

type Database struct {
	Conn *pgxpool.Pool
}

func NewDB(cfg *config.Config) *Database {
	conn, err := pgxpool.Connect(context.Background(), cfg.DatabaseConn)
	if err != nil {
		log.Fatalf("unable to start DB instance: %e", err)
	}

	return &Database{
		Conn: conn,
	}
}

func (db *Database) Close() {
	db.Conn.Close()
}

func (db *Database) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return db.Conn.Query(ctx, sql, args...)
}

func (db *Database) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return db.Conn.QueryRow(ctx, sql, args...)
}