package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type MockDB struct {
	hasErrors bool
}

func NewMockDB(hasErrors bool) *MockDB {
	return &MockDB{
		hasErrors: hasErrors,
	}
}

func (db *MockDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if db.hasErrors {
		return nil, errors.New("can't query for some reason")
	} else {
		return nil, nil
	}
}

func (db *MockDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if db.hasErrors {
		return &BadRow{}
	} else {
		return nil
	}
}

type BadRow struct{}

func (r *BadRow) Scan(dest ...any) error {
	return errors.New("can't save for some reason")
}
