package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
)

func beginTx(ctx context.Context, db *sql.DB, write bool) (*sql.Tx, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	if write {
		// This prevents SQLITE_BUSY (5) and database locked (517) when doing write transactions
		// because we tell sqlite that we are going to do a write transaction through the dummy DELETE query.
		_, _ = tx.ExecContext(ctx, "DELETE FROM _ WHERE 0 = 1;")
	}

	return tx, nil
}

type Querier interface {
	Conn() *sql.DB
	BeginTx(ctx context.Context, write bool) (QuerierTx, error)
}

type QuerierTx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Commit() error
	Rollback() error
}

func New(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("invalid database path: %s", dbPath)
	}

	db, err := connect(dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

type DB struct {
	*sql.DB
}

func NewDB(db *sql.DB) DB {
	return DB{db}
}

func (db DB) Conn() *sql.DB {
	return db.DB
}

func (db DB) BeginTx(ctx context.Context, write bool) (QuerierTx, error) {
	return beginTx(ctx, db.DB, write)
}

type DebugTx struct {
	*sql.Tx
}

type DebugDB struct {
	*sql.DB
}

func NewDebugDB(db *sql.DB) DebugDB {
	return DebugDB{db}
}

func (db DebugDB) Conn() *sql.DB {
	return db.DB
}

func (db DebugDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	log.Debug().
		Str("func", "PrepareContext").
		Msg(query)
	return db.DB.PrepareContext(ctx, query)
}

func (tx DebugTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	log.Debug().
		Str("func", "PrepareContext (Tx)").
		Msg(query)
	return tx.Tx.PrepareContext(ctx, query)
}

func (db DebugDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	log.Debug().
		Str("func", "ExecContext").
		Any("args", args).
		Msg(query)
	return db.DB.ExecContext(ctx, query, args...)
}

func (tx DebugTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	log.Debug().
		Str("func", "ExecContext (Tx)").
		Any("args", args).
		Msg(query)
	return tx.Tx.ExecContext(ctx, query, args...)
}

func (db DebugDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	log.Debug().
		Str("func", "QueryContext").
		Any("args", args).
		Msg(query)
	return db.DB.QueryContext(ctx, query, args...)
}

func (tx DebugTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	log.Debug().
		Str("func", "QueryContext (Tx)").
		Any("args", args).
		Msg(query)
	return tx.Tx.QueryContext(ctx, query, args...)
}

func (db DebugDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	log.Debug().
		Str("func", "QueryRowContext").
		Any("args", args).
		Msg(query)
	return db.DB.QueryRowContext(ctx, query, args...)
}

func (tx DebugTx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	log.Debug().
		Str("func", "QueryRowContext (Tx)").
		Any("args", args).
		Msg(query)
	return tx.Tx.QueryRowContext(ctx, query, args...)
}

func (db DebugDB) BeginTx(ctx context.Context, write bool) (QuerierTx, error) {
	log.Debug().
		Msg("BeginTx (Tx)")
	tx, err := beginTx(ctx, db.DB, write)
	if err != nil {
		return DebugTx{}, err
	}
	return DebugTx{Tx: tx}, nil
}

func (tx DebugTx) Commit() error {
	log.Debug().
		Str("func", "Commit (Tx)").
		Msg("")
	return tx.Tx.Commit()
}

func (tx DebugTx) Rollback() error {
	log.Debug().
		Str("func", "Rollback (Tx)").
		Msg("")
	return tx.Tx.Rollback()
}
