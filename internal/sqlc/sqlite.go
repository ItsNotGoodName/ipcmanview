package sqlc

import (
	"context"
)

type SQLite interface {
	DBTX
	BeginTx(ctx context.Context, write bool) (SQLiteTx, error)
}

type SQLiteTx interface {
	DBTX
	Commit() error
	Rollback() error
}

type DB struct {
	SQLite
	*Queries
}

func NewDB(sqlite SQLite) DB {
	return DB{
		SQLite:  sqlite,
		Queries: New(sqlite),
	}
}
