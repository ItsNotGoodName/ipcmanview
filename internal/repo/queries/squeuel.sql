-- name: SqueuelEnqueue :exec
INSERT INTO
  squeuel (
    id,
    task_id,
    queue,
    payload,
    timeout,
    received,
    max_received,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: SqueuelDequeue :one
UPDATE squeuel
SET
  timeout = sqlc.arg ('timeout'),
  received = received + 1
WHERE
  id = (
    SELECT
      id
    FROM
      squeuel
    WHERE
      squeuel.queue = sqlc.arg ('queue')
      AND sqlc.arg ('now') >= squeuel.timeout
      AND squeuel.received < squeuel.max_received
    ORDER BY
      created_at
    LIMIT
      1
  ) RETURNING id,
  payload,
  task_id,
  max_received;

-- name: SqueuelExtend :exec
UPDATE squeuel
SET
  timeout = ?
WHERE
  queue = ?
  AND id = ?;

-- name: SqueuelDelete :exec
DELETE FROM squeuel
where
  queue = ?
  and id = ?;

-- name: SqueuelDeleteExpired :exec
DELETE FROM squeuel
WHERE
  received >= max_received
  AND sqlc.arg ('now') >= timeout;
