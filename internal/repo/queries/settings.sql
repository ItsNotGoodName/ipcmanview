-- name: GetSettings :one
SELECT
  *
FROM
  settings
LIMIT
  1;

-- name: UpdateSettings :one
UPDATE settings
SET
  location = coalesce(sqlc.narg ('location'), location),
  site_name = coalesce(sqlc.narg ('site_name'), site_name)
WHERE
  1 = 1 RETURNING *;
