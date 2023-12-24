-- name: createDahuaDevice :one
INSERT INTO dahua_devices (
  name, address, username, password, location, feature, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING id;

-- name: dahuaDeviceExists :one
SELECT COUNT(id) FROM dahua_devices WHERE id = ?;

-- name: UpdateDahuaDevice :one
UPDATE dahua_devices 
SET name = ?, address = ?, username = ?, password = ?, location = ?, feature = ?, updated_at = ?
WHERE id = ?
RETURNING id;

-- name: GetDahuaDevice :one
SELECT dahua_devices.*, coalesce(seed, id) FROM dahua_devices 
LEFT JOIN dahua_seeds ON dahua_seeds.device_id = dahua_devices.id
WHERE id = ? LIMIT 1;

-- name: ListDahuaDevice :many
SELECT dahua_devices.*, coalesce(seed, id) FROM dahua_devices 
LEFT JOIN dahua_seeds ON dahua_seeds.device_id = dahua_devices.id;

-- name: ListDahuaDeviceByIDs :many
SELECT dahua_devices.*, coalesce(seed, id) FROM dahua_devices 
LEFT JOIN dahua_seeds ON dahua_seeds.device_id = dahua_devices.id
WHERE id IN (sqlc.slice('ids'));

-- name: listDahuaDeviceByFeature :many
SELECT dahua_devices.*, coalesce(seed, id) FROM dahua_devices 
LEFT JOIN dahua_seeds ON dahua_seeds.device_id = dahua_devices.id
WHERE feature & sqlc.arg('feature') = sqlc.arg('feature');

-- name: DeleteDahuaDevice :exec
DELETE FROM dahua_devices WHERE id = ?;

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

-- name: allocateDahuaSeed :exec
UPDATE dahua_seeds 
SET device_id = ?1
WHERE seed = (SELECT seed FROM dahua_seeds WHERE device_id = ?1 OR device_id IS NULL ORDER BY device_id asc LIMIT 1);

-- name: CreateDahuaFileScanLock :one
INSERT INTO dahua_file_scan_locks (
  device_id, touched_at
) VALUES (
  ?, ?
) RETURNING *;

-- name: DeleteDahuaFileScanLock :exec
DELETE FROM dahua_file_scan_locks WHERE device_id = ?;

-- name: DeleteDahuaFileScanLockByAge :exec
DELETE FROM dahua_file_scan_locks WHERE touched_at < ?;

-- name: TouchDahuaFileScanLock :exec
UPDATE dahua_file_scan_locks
SET touched_at = ?
WHERE device_id = ?;

-- name: UpdateDahuaFileCursorPercent :one
UPDATE dahua_file_cursors 
SET
  percent = ?
WHERE device_id = ?
RETURNING *;

-- name: ListDahuaFileCursor :many
SELECT 
  c.*,
  count(f.device_id) AS files,
  coalesce(l.touched_at > ?, false) AS locked
FROM dahua_file_cursors AS c
LEFT JOIN dahua_files AS f ON f.device_id = c.device_id
LEFT JOIN dahua_file_scan_locks AS l ON l.device_id = c.device_id
GROUP BY c.device_id;

-- name: UpdateDahuaFileCursor :one
UPDATE dahua_file_cursors
SET 
  quick_cursor = ?,
  full_cursor = ?,
  full_epoch = ?,
  percent = ?
WHERE device_id = ?
RETURNING *;

-- name: createDahuaFileCursor :exec
INSERT INTO dahua_file_cursors (
  device_id,
  quick_cursor,
  full_cursor,
  full_epoch,
  percent
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: ListDahuaFileTypes :many
SELECT DISTINCT type
FROM dahua_files;

-- name: CreateDahuaFile :one
INSERT INTO dahua_files (
  device_id,
  channel,
  start_time,
  end_time,
  length,
  type,
  file_path,
  duration,
  disk,
  video_stream,
  flags,
  events,
  cluster,
  partition,
  pic_index,
  repeat,
  work_dir,
  work_dir_sn,
  updated_at,
  local
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING id;

-- name: GetDahuaFileByFilePath :one
SELECT *
FROM dahua_files
WHERE device_id = ? and file_path = ?;

-- name: UpdateDahuaFile :one
UPDATE dahua_files 
SET 
  channel = ?,
  start_time = ?,
  end_time = ?,
  length = ?,
  type = ?,
  duration = ?,
  disk = ?,
  video_stream = ?,
  flags = ?,
  events = ?,
  cluster = ?,
  partition = ?,
  pic_index = ?,
  repeat = ?,
  work_dir = ?,
  work_dir_sn = ?,
  updated_at = ?,
  local = ?
WHERE device_id = ? AND file_path = ?
RETURNING id;

-- name: DeleteDahuaFile :exec
DELETE FROM dahua_files
WHERE
  updated_at < sqlc.arg('updated_at') AND
  device_id = sqlc.arg('device_id') AND
  start_time <= sqlc.arg('end') AND
  sqlc.arg('start') < start_time;

-- name: CreateDahuaEvent :one
INSERT INTO dahua_events (
  device_id,
  code,
  action,
  `index`,
  data,
  created_at
) VALUES (
  ?, ?, ?, ?, ?, ?
) RETURNING id;

-- name: ListDahuaEventCodes :many
SELECT DISTINCT code FROM dahua_events;

-- name: ListDahuaEventActions :many
SELECT DISTINCT action FROM dahua_events;

-- name: GetDahuaEventData :one
SELECT data FROM dahua_events WHERE id = ?;

-- name: getDahuaEventRuleByEvent :many
SELECT 
  ignore_db,
  ignore_live,
  ignore_mqtt,
  code
FROM dahua_event_device_rules 
WHERE device_id = sqlc.arg('device_id') AND (dahua_event_device_rules.code = sqlc.arg('code') OR dahua_event_device_rules.code = '')
UNION ALL
SELECT 
  ignore_db,
  ignore_live,
  ignore_mqtt,
  code
FROM dahua_event_rules
WHERE dahua_event_rules.code = sqlc.arg('code') OR dahua_event_rules.code = ''
ORDER BY code DESC;

-- name: GetDahuaEventRule :one
SELECT * FROM dahua_event_rules
WHERE id = ?;

-- name: ListDahuaEventRule :many
SELECT * FROM dahua_event_rules;

-- name: UpdateDahuaEventRule :exec
UPDATE dahua_event_rules 
SET 
  code = ?,
  ignore_db = ?,
  ignore_live = ?,
  ignore_mqtt = ?
WHERE id = ?;

-- name: CreateDahuaEventRule :exec
INSERT INTO dahua_event_rules(
  code,
  ignore_db,
  ignore_live,
  ignore_mqtt
) VALUES(
  ?,
  ?,
  ?,
  ?
);

-- name: DeleteDahuaEventRule :exec
DELETE FROM dahua_event_rules WHERE id = ?;

-- name: CreateDahuaEventWorkerState :exec
INSERT INTO dahua_event_worker_states(
  device_id,
  state,
  error,
  created_at
) VALUES(
  ?,
  ?,
  ?,
  ?
);

-- name: ListDahuaEventWorkerState :many
SELECT *,max(created_at) FROM dahua_event_worker_states GROUP BY device_id;
