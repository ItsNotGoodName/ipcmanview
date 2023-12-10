package repo

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

func (db DB) BeginTx(ctx context.Context, write bool) (DBTx, error) {
	tx, err := db.SQLite.BeginTx(ctx, write)
	if err != nil {
		return DBTx{}, err
	}
	return DBTx{
		SQLiteTx: tx,
		Queries:  New(tx),
	}, nil
}

type DBTx struct {
	SQLiteTx
	*Queries
}
