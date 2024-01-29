// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: settings.sql

package repo

import (
	"context"
	"database/sql"
)

const getSettings = `-- name: GetSettings :one
SELECT
  setup, site_name, location, coordinates
FROM
  settings
LIMIT
  1
`

func (q *Queries) GetSettings(ctx context.Context) (Setting, error) {
	row := q.db.QueryRowContext(ctx, getSettings)
	var i Setting
	err := row.Scan(
		&i.Setup,
		&i.SiteName,
		&i.Location,
		&i.Coordinates,
	)
	return i, err
}

const updateSettings = `-- name: UpdateSettings :one
UPDATE settings
SET
  location = coalesce(?1, location),
  site_name = coalesce(?2, site_name)
WHERE
  1 = 1 RETURNING setup, site_name, location, coordinates
`

type UpdateSettingsParams struct {
	Location sql.NullString
	SiteName sql.NullString
}

func (q *Queries) UpdateSettings(ctx context.Context, arg UpdateSettingsParams) (Setting, error) {
	row := q.db.QueryRowContext(ctx, updateSettings, arg.Location, arg.SiteName)
	var i Setting
	err := row.Scan(
		&i.Setup,
		&i.SiteName,
		&i.Location,
		&i.Coordinates,
	)
	return i, err
}