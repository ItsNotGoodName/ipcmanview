-- name: createDahuaDevice :one
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

-- name: dahuaDeviceExists :one
SELECT
  COUNT(id)
FROM
  dahua_devices
WHERE
  id = ?;

-- name: UpdateDahuaDevice :one
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

-- name: GetDahuaDeviceName :one
SELECT
  name
FROM
  dahua_devices
WHERE
  id = ?;

-- name: GetDahuaDevice :one
SELECT
  dahua_devices.*,
  coalesce(seed, id)
FROM
  dahua_devices
  LEFT JOIN dahua_seeds ON dahua_seeds.device_id = dahua_devices.id
WHERE
  id = ?
LIMIT
  1;

-- name: GetDahuaDeviceByIP :one
SELECT
  dahua_devices.*,
  coalesce(seed, id)
FROM
  dahua_devices
  LEFT JOIN dahua_seeds ON dahua_seeds.device_id = dahua_devices.id
WHERE
  ip = ?
LIMIT
  1;

-- name: ListDahuaDevice :many
SELECT
  dahua_devices.*,
  coalesce(seed, id)
FROM
  dahua_devices
  LEFT JOIN dahua_seeds ON dahua_seeds.device_id = dahua_devices.id;

-- name: ListDahuaDeviceByIDs :many
SELECT
  dahua_devices.*,
  coalesce(seed, id)
FROM
  dahua_devices
  LEFT JOIN dahua_seeds ON dahua_seeds.device_id = dahua_devices.id
WHERE
  id IN (sqlc.slice ('ids'));

-- name: listDahuaDeviceByFeature :many
SELECT
  dahua_devices.*,
  coalesce(seed, id)
FROM
  dahua_devices
  LEFT JOIN dahua_seeds ON dahua_seeds.device_id = dahua_devices.id
WHERE
  feature & sqlc.arg ('feature') = sqlc.arg ('feature');

-- name: DeleteDahuaDevice :exec
DELETE FROM dahua_devices
WHERE
  id = ?;

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
  default_location = coalesce(sqlc.narg ('default_location'), default_location),
  site_name = coalesce(sqlc.narg ('site_name'), site_name)
WHERE
  1 = 1 RETURNING *;

-- name: allocateDahuaSeed :exec
UPDATE dahua_seeds
SET
  device_id = ?1
WHERE
  seed = (
    SELECT
      seed
    FROM
      dahua_seeds
    WHERE
      device_id = ?1
      OR device_id IS NULL
    ORDER BY
      device_id asc
    LIMIT
      1
  );

-- name: NormalizeDahuaFileCursor :exec
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
SELECT
  id,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?
FROM
  dahua_devices;

-- name: UpdateDahuaFileCursorScanPercent :one
UPDATE dahua_file_cursors
SET
  scan_percent = ?
WHERE
  device_id = ? RETURNING *;

-- name: ListDahuaFileCursor :many
SELECT
  c.*,
  count(f.device_id) AS files
FROM
  dahua_file_cursors AS c
  LEFT JOIN dahua_files AS f ON f.device_id = c.device_id
GROUP BY
  c.device_id;

-- name: UpdateDahuaFileCursor :one
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

-- name: createDahuaFileCursor :exec
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

-- name: ListDahuaFileTypes :many
SELECT DISTINCT
  type
FROM
  dahua_files;

-- name: CreateDahuaFile :one
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
  (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
  )
ON CONFLICT (start_time) DO
UPDATE
SET
  id = id RETURNING id;

-- name: GetDahuaFile :one
SELECT
  *
FROM
  dahua_files
WHERE
  id = ?;

-- name: GetDahuaFileForThumbnail :one
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

-- name: GetDahuaFileByFilePath :one
SELECT
  *
FROM
  dahua_files
WHERE
  device_id = ?
  and file_path = ?;

-- name: GetOldestDahuaFileStartTime :one
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
  storage = ?
WHERE
  device_id = ?
  AND file_path = ? RETURNING id;

-- name: DeleteDahuaFile :exec
DELETE FROM dahua_files
WHERE
  updated_at < sqlc.arg ('updated_at')
  AND device_id = sqlc.arg ('device_id')
  AND start_time <= sqlc.arg ('end')
  AND sqlc.arg ('start') < start_time;

-- name: CreateDahuaThumbnail :one
INSERT INTO
  dahua_thumbnails (file_id, email_attachment_id, width, height)
VALUES
  (?, ?, ?, ?) RETURNING *;

-- name: OrphanDeleteDahuaThumbnail :exec
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

-- name: CreateDahuaEvent :one
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

-- name: ListDahuaEventCodes :many
SELECT DISTINCT
  code
FROM
  dahua_events;

-- name: ListDahuaEventActions :many
SELECT DISTINCT
  action
FROM
  dahua_events;

-- name: GetDahuaEventData :one
SELECT
  data
FROM
  dahua_events
WHERE
  id = ?;

-- name: DeleteDahuaEvent :exec
DELETE FROM dahua_events;

-- name: getDahuaEventRuleByEvent :many
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

-- name: GetDahuaEventRule :one
SELECT
  *
FROM
  dahua_event_rules
WHERE
  id = ?;

-- name: ListDahuaEventRule :many
SELECT
  *
FROM
  dahua_event_rules;

-- name: UpdateDahuaEventRule :exec
UPDATE dahua_event_rules
SET
  code = ?,
  ignore_db = ?,
  ignore_live = ?,
  ignore_mqtt = ?
WHERE
  id = ?;

-- name: CreateDahuaEventRule :exec
INSERT INTO
  dahua_event_rules (code, ignore_db, ignore_live, ignore_mqtt)
VALUES
  (?, ?, ?, ?);

-- name: DeleteDahuaEventRule :exec
DELETE FROM dahua_event_rules
WHERE
  id = ?;

-- name: CreateDahuaEventWorkerState :exec
INSERT INTO
  dahua_event_worker_states (device_id, state, error, created_at)
VALUES
  (?, ?, ?, ?);

-- name: ListDahuaEventWorkerState :many
SELECT
  *,
  max(created_at)
FROM
  dahua_event_worker_states
GROUP BY
  device_id;

-- name: GetDahuaStorageDestination :one
SELECT
  *
FROM
  dahua_storage_destinations
WHERE
  id = ?;

-- name: GetDahuaStorageDestinationByServerAddressAndStorage :one
SELECT
  *
FROM
  dahua_storage_destinations
WHERE
  server_address = ?
  AND storage = ?;

-- name: ListDahuaStorageDestination :many
SELECT
  *
FROM
  dahua_storage_destinations;

-- name: CreateDahuaStorageDestination :one
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

-- name: UpdateDahuaStorageDestination :one
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
  id = ? RETURNING *;

-- name: DeleteDahuaStorageDestination :exec
DELETE FROM dahua_storage_destinations
WHERE
  id = ?;

-- name: createDahuaStreamDefault :one
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

-- name: updateDahuaStreamDefault :exec
UPDATE dahua_streams
SET
  internal = false
WHERE
  device_id = ?;

-- name: DeleteDahuaStream :exec
DELETE FROM dahua_streams
WHERE
  id = ?;

-- name: ListDahuaStreamByDevice :many
SELECT
  *
FROM
  dahua_streams
WHERE
  device_id = ?;

-- name: ListDahuaStream :many
SELECT
  *
FROM
  dahua_streams
ORDER BY
  device_id;

-- name: GetDahuaStream :one
SELECT
  *
FROM
  dahua_streams
WHERE
  id = ?;

-- name: UpdateDahuaStream :one
UPDATE dahua_streams
SET
  name = ?,
  mediamtx_path = ?
WHERE
  id = ? RETURNING *;

-- name: createDahuaEmailMessage :one
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

-- name: createDahuaEmailAttachment :one
INSERT INTO
  dahua_email_attachments (message_id, file_name)
VALUES
  (?, ?) RETURNING *;

-- name: CreateDahuaAferoFile :one
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

-- name: GetDahuaAferoFileByFileID :one
SELECT
  *
FROM
  dahua_afero_files
WHERE
  file_id = ?;

-- name: ReadyDahuaAferoFile :one
UPDATE dahua_afero_files
SET
  ready = true,
  size = ?,
  created_at = ?
WHERE
  id = ? RETURNING id;

-- name: DeleteDahuaAferoFile :exec
DELETE FROM dahua_afero_files
WHERE
  id = ?;

-- name: OrphanListDahuaAferoFile :many
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
