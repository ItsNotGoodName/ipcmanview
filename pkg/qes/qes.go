package qes

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type QueryExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func QueryMany(ctx context.Context, qe QueryExecutor, stmt postgres.Statement) (pgx.Rows, error) {
	sql, args := stmt.Sql()
	return qe.Query(ctx, sql, args...)
}

func QueryOne(ctx context.Context, qe QueryExecutor, stmt postgres.Statement) pgx.Row {
	sql, args := stmt.Sql()
	return qe.QueryRow(ctx, sql, args...)
}

func Exec(ctx context.Context, qe QueryExecutor, stmt postgres.Statement) (pgconn.CommandTag, error) {
	sql, args := stmt.Sql()
	return qe.Exec(ctx, sql, args...)
}

func ExecOne(ctx context.Context, qe QueryExecutor, stmt postgres.Statement) (pgconn.CommandTag, error) {
	sql, args := stmt.Sql()
	res, err := qe.Exec(ctx, sql, args...)
	if err != nil {
		return res, err
	}

	if res.RowsAffected() == 0 {
		return res, pgx.ErrNoRows
	}

	return res, nil
}

func ScanMany(ctx context.Context, qe QueryExecutor, dst any, stmt postgres.Statement) error {
	sql, args := stmt.Sql()
	return pgxscan.Select(ctx, qe, dst, sql, args...)
}

func ScanOne(ctx context.Context, qe QueryExecutor, dst any, stmt postgres.Statement) error {
	sql, args := stmt.Sql()
	return pgxscan.Get(ctx, qe, dst, sql, args...)
}
