-- name: CreateDahuaCamera :one
INSERT INTO dahua_cameras (
  name, address, username, password, location, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateDahuaCamera :one
UPDATE dahua_cameras 
SET name = ?, address = ?, username = ?, password = ?, location = ?
WHERE id = ? RETURNING *;

-- name: GetDahuaCamera :one
SELECT * FROM dahua_cameras
WHERE id = ? LIMIT 1;

-- name: ListDahuaCamera :many
SELECT * FROM dahua_cameras;

-- name: DeleteDahuaCamera :exec
DELETE FROM dahua_cameras WHERE id = ?;

-- name: CreateDahuaEvent :one
INSERT INTO dahua_events (
  camera_id,
  content_type,
  content_length,
  code,
  action,
  `index`,
  data,
  created_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: ListDahuaEvent :many
SELECT * FROM dahua_events
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
