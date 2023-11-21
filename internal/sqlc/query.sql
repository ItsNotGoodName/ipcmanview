-- name: CreateDahuaCamera :one
INSERT INTO dahua_cameras (
  name, address, username, password, location, created_at
) VALUES (
  ?, ?, ?, ?, ?, ?
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

