-- name: createEvent :one
INSERT INTO
  events (action, data, user_id, actor, created_at)
VALUES
  (?, ?, ?, ?, ?) RETURNING id;

-- name: NextEventByCursor :one
SELECT
  *
from
  events
WHERE
  id > ?
LIMIT
  1;

-- name: GetEventCursor :one
SELECT
  id
from
  events
ORDER BY
  id DESC
LIMIT
  1;
