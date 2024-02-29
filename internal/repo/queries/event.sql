-- name: CreateEvent :one
INSERT INTO
  events (action, data, user_id, actor, created_at)
VALUES
  (?, ?, ?, ?, ?) RETURNING id;
