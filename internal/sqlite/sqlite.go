package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/rs/zerolog/log"
)

const (
	CONSTRAINT_CHECK      = 275
	CONSTRAINT_COMMITHOOK = 531
	CONSTRAINT_DATATYPE   = 3091
	CONSTRAINT_FOREIGNKEY = 787
	CONSTRAINT_FUNCTION   = 1043
	CONSTRAINT_NOTNULL    = 1299
	CONSTRAINT_PINNED     = 2835
	CONSTRAINT_PRIMARYKEY = 1555
	CONSTRAINT_ROWID      = 2579
	CONSTRAINT_TRIGGER    = 1811
	CONSTRAINT_UNIQUE     = 2067
	CONSTRAINT_VTAB       = 2323
)

type Error struct {
	Code int
	Msg  string
}

func AsConstraintError(err error, code int, codes ...int) (ConstraintError, bool) {
	e, ok := AsError(err)
	if !(ok || e.Code == code || slices.Contains(codes, e.Code)) {
		return ConstraintError{}, false
	}
	return ConstraintError(e), true
}

type ConstraintError Error

func (e ConstraintError) IsField(field string) bool {
	return strings.Contains(string(e.Msg), field)
}

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

func (db DB) BeginTx(ctx context.Context, write bool) (repo.SQLiteTx, error) {
	return beginTx(ctx, db.DB, write)
}

type DebugTx struct {
	tx *sql.Tx
}

type DebugDB struct {
	db *sql.DB
}

func NewDebugDB(db *sql.DB) DebugDB {
	return DebugDB{db}
}

func (db DebugDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	log.Debug().
		Str("func", "PrepareContext").
		Msg(query)
	return db.db.PrepareContext(ctx, query)
}

func (tx DebugTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	log.Debug().
		Str("func", "PrepareContext (Tx)").
		Msg(query)
	return tx.tx.PrepareContext(ctx, query)
}

func (db DebugDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	log.Debug().
		Str("func", "ExecContext").
		Any("args", args).
		Msg(query)
	return db.db.ExecContext(ctx, query, args...)
}

func (tx DebugTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	log.Debug().
		Str("func", "ExecContext (Tx)").
		Any("args", args).
		Msg(query)
	return tx.tx.ExecContext(ctx, query, args...)
}

func (db DebugDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	log.Debug().
		Str("func", "QueryContext").
		Any("args", args).
		Msg(query)
	return db.db.QueryContext(ctx, query, args...)
}

func (tx DebugTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	log.Debug().
		Str("func", "QueryContext (Tx)").
		Any("args", args).
		Msg(query)
	return tx.tx.QueryContext(ctx, query, args...)
}

func (db DebugDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	log.Debug().
		Str("func", "QueryRowContext").
		Any("args", args).
		Msg(query)
	return db.db.QueryRowContext(ctx, query, args...)
}

func (tx DebugTx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	log.Debug().
		Str("func", "QueryRowContext (Tx)").
		Any("args", args).
		Msg(query)
	return tx.tx.QueryRowContext(ctx, query, args...)
}

func (db DebugDB) BeginTx(ctx context.Context, write bool) (repo.SQLiteTx, error) {
	log.Debug().
		Msg("BeginTx (Tx)")
	tx, err := beginTx(ctx, db.db, write)
	if err != nil {
		return DebugTx{}, err
	}
	return DebugTx{tx: tx}, nil
}

func (tx DebugTx) Commit() error {
	log.Debug().
		Str("func", "Commit (Tx)").
		Msg("")
	return tx.tx.Commit()
}

func (tx DebugTx) Rollback() error {
	log.Debug().
		Str("func", "Rollback (Tx)").
		Msg("")
	return tx.tx.Rollback()
}
