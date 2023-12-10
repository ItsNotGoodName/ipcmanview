// ssq is glue code for squirrel and sqlscan.
package ssq

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/sqlscan"
)

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

	return sqlscan.ScanOne(dst, rows)
}
