// Package ssq is glue code for squirrel and sqlscan.
package ssq

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/sqlscan"
)

func QueryRows(ctx context.Context, db sqlscan.Querier, sb sq.SelectBuilder) (*sql.Rows, *sqlscan.RowScanner, error) {
	sql, args, err := sb.ToSql()
	if err != nil {
		return nil, nil, err
	}

	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, nil, err
	}

	return rows, sqlscan.NewRowScanner(rows), nil
}

func Query(ctx context.Context, db sqlscan.Querier, dst any, sb sq.SelectBuilder) error {
	sql, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	return sqlscan.Select(ctx, db, dst, sql, args...)
}

func QueryOne(ctx context.Context, db sqlscan.Querier, dst any, sb sq.SelectBuilder) error {
	sql, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return sqlscan.ScanOne(dst, rows)
}
