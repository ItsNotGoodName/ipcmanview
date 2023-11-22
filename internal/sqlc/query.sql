-- name: createDahuaCamera :one
INSERT INTO dahua_cameras (
  name, address, username, password, location, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
) RETURNING id;

-- name: UpdateDahuaCamera :one
UPDATE dahua_cameras 
SET name = ?, address = ?, username = ?, password = ?, location = ?
WHERE id = ?
RETURNING id;

-- name: GetDahuaCamera :one
SELECT id, name, address, username, password, location, created_at, coalesce(seed, id) FROM dahua_cameras 
LEFT JOIN dahua_seeds ON dahua_seeds.camera_id = dahua_cameras.id
WHERE id = ? LIMIT 1;

-- name: ListDahuaCamera :many
SELECT id, name, address, username, password, location, created_at, coalesce(seed, id) FROM dahua_cameras 
LEFT JOIN dahua_seeds ON dahua_seeds.camera_id = dahua_cameras.id;

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
) RETURNING id;

-- name: ListDahuaEvent :many
SELECT * FROM dahua_events
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetSettings :one
SELECT * FROM settings
LIMIT 1;

-- name: UpdateSettings :one
UPDATE settings
SET
  default_location = coalesce(sqlc.narg('default_location'), default_location),
  site_name = coalesce(sqlc.narg('site_name'), site_name)
WHERE 1 = 1
RETURNING *;

-- name: setDahuaSeed :exec
UPDATE dahua_seeds 
SET camera_id = ?1
WHERE seed = (SELECT seed FROM dahua_seeds WHERE camera_id = ?1 OR camera_id IS NULL ORDER BY camera_id asc LIMIT 1);

-- name: CreateDahuaFileScanLock :one
INSERT INTO dahua_file_scan_locks (
  camera_id, created_at
) VALUES (
  ?, ?
) RETURNING *;

-- name: DeleteDahuaFileScanLock :exec
DELETE FROM dahua_file_scan_locks WHERE camera_id = ?;

-- name: GetDahuaFileCursor :one
SELECT * FROM dahua_file_cursors 
WHERE camera_id = ?;
