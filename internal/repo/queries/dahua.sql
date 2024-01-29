-- name: DahuaCreateDevice :one
INSERT INTO
  dahua_devices (
    name,
    url,
    ip,
    username,
    password,
    location,
    feature,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id;

-- name: DahuaCheckDevice :one
SELECT
  COUNT(*) = 1
FROM
  dahua_devices
WHERE
  id = ?;

-- name: DahuaUpdateDevice :one
UPDATE dahua_devices
SET
  name = ?,
  url = ?,
  ip = ?,
  username = ?,
  password = ?,
  location = ?,
  feature = ?,
  updated_at = ?
WHERE
  id = ? RETURNING id;

-- name: DahuaGetDeviceName :one
SELECT
  name
FROM
  dahua_devices
WHERE
  id = ?;

-- name: DahuaListDevicePermissions :many
SELECT
  dahua_devices.id,
  coalesce(p.level, 2)
FROM
  dahua_devices
  LEFT JOIN dahua_permissions AS p ON p.device_id = dahua_devices.id
WHERE
  -- Allow if user is admin
  EXISTS (SELECT user_id FROM admins WHERE admins.user_id = sqlc.arg ('user_id'))
  -- Allow if user is a part of the group the owns the permission
  OR p.group_id IN (
    SELECT
      group_id
    FROM
      group_users
    WHERE
      group_users.user_id = sqlc.arg ('user_id')
  )
  -- Allow if user owns the permission
  OR p.user_id = sqlc.arg ('user_id')
GROUP BY
  -- Remove duplicate devices with different permissions
  dahua_devices.id
ORDER BY
  -- Get the highest permission level
  p.level DESC;

-- name: DahuaDeleteDevice :exec
DELETE FROM dahua_devices
WHERE
  id = ?;

-- name: DahuaAllocateSeed :exec
UPDATE dahua_seeds
SET
  device_id = sqlc.arg('device_id')
WHERE
  seed = (
    SELECT
      seed
    FROM
      dahua_seeds
    WHERE
      device_id = sqlc.arg('device_id')
      OR device_id IS NULL
    ORDER BY
      device_id ASC
    LIMIT
      1
  );

-- name: DahuaNormalizeFileCursors :exec
INSERT OR IGNORE INTO
  dahua_file_cursors (
    device_id,
    quick_cursor,
    full_cursor,
    full_epoch,
    scan,
    scan_percent,
    scan_type
  )
SELECT id, ?, ?, ?, ?, ?, ?
FROM dahua_devices;

-- name: DahuaUpdateFileCursorScanPercent :one
UPDATE dahua_file_cursors
SET
  scan_percent = ?
WHERE
  device_id = ? RETURNING *;

-- name: DahuaListFileCursors :many
SELECT
  c.*,
  count(f.device_id) AS files
FROM
  dahua_file_cursors AS c
  LEFT JOIN dahua_files AS f ON f.device_id = c.device_id
GROUP BY
  c.device_id;

-- name: DahuaUpdateFileCursor :one
UPDATE dahua_file_cursors
SET
  quick_cursor = ?,
  full_cursor = ?,
  full_epoch = ?,
  scan = ?,
  scan_percent = ?,
  scan_type = ?
WHERE
  device_id = ? RETURNING *;

-- name: DahuaCreateFileCursor :exec
INSERT INTO
  dahua_file_cursors (
    device_id,
    quick_cursor,
    full_cursor,
    full_epoch,
    scan,
    scan_percent,
    scan_type
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?);

-- name: DahuaListFileTypes :many
SELECT DISTINCT
  type
FROM
  dahua_files;

-- name: DahuaCreateFile :one
INSERT INTO
  dahua_files (
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
    storage
  )
VALUES 
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (start_time) DO
UPDATE
SET
  id = id RETURNING id;

-- name: DahuaGetFile :one
SELECT
  *
FROM
  dahua_files
WHERE
  id = ?;

-- name: DahuaGetFileForThumbnail :one
SELECT
  dahua_files.id,
  device_id,
  type,
  file_path,
  name,
  ready
FROM
  dahua_files
  LEFT JOIN dahua_thumbnails ON dahua_thumbnails.file_id = dahua_files.id
  LEFT JOIN dahua_afero_files ON dahua_afero_files.thumbnail_id = dahua_thumbnails.id
WHERE
  dahua_files.id = ?;

-- name: DahuaGetFileByFilePath :one
SELECT
  *
FROM
  dahua_files
WHERE
  device_id = ?
  and file_path = ?;

-- name: DahuaGetOldestFileStartTime :one
SELECT
  start_time
FROM
  dahua_files
WHERE
  device_id = ?
ORDER BY
  start_time ASC
LIMIT
  1;

-- name: DahuaUpdateFile :one
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
  storage = ?
WHERE
  device_id = ?
  AND file_path = ? RETURNING id;

-- name: DahuaDeleteFile :exec
DELETE FROM dahua_files
WHERE
  updated_at < sqlc.arg ('updated_at')
  AND device_id = sqlc.arg ('device_id')
  AND start_time <= sqlc.arg ('end')
  AND sqlc.arg ('start') < start_time;

-- name: DahuaCreateThumbnail :one
INSERT INTO
  dahua_thumbnails (file_id, email_attachment_id, width, height)
VALUES
  (?, ?, ?, ?) RETURNING *;

-- name: DahuaOrphanDeleteThumbnail :exec
DELETE FROM dahua_thumbnails
WHERE
  id IN (
    SELECT
      thumbnail_id
    FROM
      dahua_afero_files
    WHERE
      ready = false
      AND created_at < ?
  );

-- name: DahuaCreateEvent :one
INSERT INTO
  dahua_events (
    device_id,
    code,
    action,
    `index`,
    data,
    created_at
  )
VALUES
  (?, ?, ?, ?, ?, ?) RETURNING id;

-- name: DahuaListEventCodes :many
SELECT DISTINCT
  code
FROM
  dahua_events;

-- name: DahuaListEventActions :many
SELECT DISTINCT
  action
FROM
  dahua_events;

-- name: DahuaGetEventData :one
SELECT
  data
FROM
  dahua_events
WHERE
  id = ?;

-- name: DahuaDeleteEvent :exec
DELETE FROM dahua_events;

-- name: DahuaGetEventRuleByEvent :many
SELECT
  ignore_db,
  ignore_live,
  ignore_mqtt,
  code
FROM
  dahua_event_device_rules
WHERE
  device_id = sqlc.arg ('device_id')
  AND (
    dahua_event_device_rules.code = sqlc.arg ('code')
    OR dahua_event_device_rules.code = ''
  )
UNION ALL
SELECT
  ignore_db,
  ignore_live,
  ignore_mqtt,
  code
FROM
  dahua_event_rules
WHERE
  dahua_event_rules.code = sqlc.arg ('code')
  OR dahua_event_rules.code = ''
ORDER BY
  code DESC;

-- name: DahuaGetEventRule :one
SELECT
  *
FROM
  dahua_event_rules
WHERE
  id = ?;

-- name: DahuaListEventRules :many
SELECT
  *
FROM
  dahua_event_rules;

-- name: DahuaUpdateEventRule :exec
UPDATE dahua_event_rules
SET
  code = ?,
  ignore_db = ?,
  ignore_live = ?,
  ignore_mqtt = ?
WHERE
  id = ?;

-- name: DahuaCreateEventRule :exec
INSERT INTO
  dahua_event_rules (code, ignore_db, ignore_live, ignore_mqtt)
VALUES
  (?, ?, ?, ?);

-- name: DahuaDeleteEventRule :exec
DELETE FROM dahua_event_rules
WHERE
  id = ?;

-- name: DahuaCreateEventWorkerState :exec
INSERT INTO
  dahua_event_worker_states (device_id, state, error, created_at)
VALUES
  (?, ?, ?, ?);

-- name: DahuaListEventWorkerStates :many
SELECT
  *,
  max(created_at)
FROM
  dahua_event_worker_states
GROUP BY
  device_id;

-- name: DahuaGetStorageDestination :one
SELECT
  *
FROM
  dahua_storage_destinations
WHERE
  id = ?;

-- name: DahuaGetStorageDestinationByServerAddressAndStorage :one
SELECT
  *
FROM
  dahua_storage_destinations
WHERE
  server_address = ?
  AND storage = ?;

-- name: DahuaListStorageDestinations :many
SELECT
  *
FROM
  dahua_storage_destinations;

-- name: DahuaCreateStorageDestination :one
INSERT INTO
  dahua_storage_destinations (
    name,
    storage,
    server_address,
    port,
    username,
    password,
    remote_directory
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?) RETURNING id;

-- name: DahuaUpdateStorageDestination :one
UPDATE dahua_storage_destinations
SET
  name = ?,
  storage = ?,
  server_address = ?,
  port = ?,
  username = ?,
  password = ?,
  remote_directory = ?
WHERE
  id = ? RETURNING id;

-- name: DahuaDeleteStorageDestination :exec
DELETE FROM dahua_storage_destinations
WHERE
  id = ?;

-- name: DahuaCreateStreamForInternal :one
INSERT INTO
  dahua_streams (
    device_id,
    channel,
    subtype,
    name,
    mediamtx_path,
    internal
  )
VALUES
  (?, ?, ?, ?, ?, true)
ON CONFLICT DO
UPDATE
SET
  internal = true RETURNING ID;

-- name: DahuaUpdateStreamForInternal :exec
UPDATE dahua_streams
SET
  internal = false
WHERE
  device_id = ?;

-- name: DahuaDeleteStream :exec
DELETE FROM dahua_streams
WHERE
  id = ?;

-- name: DahuaListStreamsByDevice :many
SELECT
  *
FROM
  dahua_streams
WHERE
  device_id = ?;

-- name: DahuaListStreams :many
SELECT
  *
FROM
  dahua_streams
ORDER BY
  device_id;

-- name: DahuaGetStream :one
SELECT
  *
FROM
  dahua_streams
WHERE
  id = ?;

-- name: DahuaUpdateStream :one
UPDATE dahua_streams
SET
  name = ?,
  mediamtx_path = ?
WHERE
  id = ? RETURNING *;

-- name: DahuaCreateEmailMessage :one
INSERT INTO
  dahua_email_messages (
    device_id,
    date,
    'from',
    `to`,
    subject,
    `text`,
    alarm_event,
    alarm_input_channel,
    alarm_name,
    created_at
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: DahuaCreateEmailAttachment :one
INSERT INTO
  dahua_email_attachments (message_id, file_name)
VALUES
  (?, ?) RETURNING *;

-- name: DahuaCreateAferoFile :one
INSERT INTO
  dahua_afero_files (
    file_id,
    thumbnail_id,
    email_attachment_id,
    name,
    created_at
  )
VALUES
  (?, ?, ?, ?, ?) RETURNING id;

-- name: DahuaGetAferoFileByFileID :one
SELECT
  *
FROM
  dahua_afero_files
WHERE
  file_id = ?;

-- name: DahuaReadyAferoFile :one
UPDATE dahua_afero_files
SET
  ready = true,
  size = ?,
  created_at = ?
WHERE
  id = ? RETURNING id;

-- name: DahuaDeleteAferoFile :exec
DELETE FROM dahua_afero_files
WHERE
  id = ?;

-- name: DahuaOrphanListAferoFiles :many
SELECT
  *
FROM
  dahua_afero_files
WHERE
  file_id IS NULL
  AND thumbnail_id IS NULL
  AND email_attachment_id IS NULL
  AND ready = true
LIMIT
  ?;
