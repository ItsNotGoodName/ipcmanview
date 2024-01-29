// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: dahua.sql

package repo

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

const dahuaAllocateSeed = `-- name: DahuaAllocateSeed :exec
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
      device_id ASC
    LIMIT
      1
  )
`

func (q *Queries) DahuaAllocateSeed(ctx context.Context, deviceID sql.NullInt64) error {
	_, err := q.db.ExecContext(ctx, dahuaAllocateSeed, deviceID)
	return err
}

const dahuaCheckDevice = `-- name: DahuaCheckDevice :one
SELECT
  COUNT(*) = 1
FROM
  dahua_devices
WHERE
  id = ?
`

func (q *Queries) DahuaCheckDevice(ctx context.Context, id int64) (bool, error) {
	row := q.db.QueryRowContext(ctx, dahuaCheckDevice, id)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}

const dahuaCreateAferoFile = `-- name: DahuaCreateAferoFile :one
INSERT INTO
  dahua_afero_files (
    file_id,
    thumbnail_id,
    email_attachment_id,
    name,
    created_at
  )
VALUES
  (?, ?, ?, ?, ?) RETURNING id
`

type DahuaCreateAferoFileParams struct {
	FileID            sql.NullInt64
	ThumbnailID       sql.NullInt64
	EmailAttachmentID sql.NullInt64
	Name              string
	CreatedAt         types.Time
}

func (q *Queries) DahuaCreateAferoFile(ctx context.Context, arg DahuaCreateAferoFileParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, dahuaCreateAferoFile,
		arg.FileID,
		arg.ThumbnailID,
		arg.EmailAttachmentID,
		arg.Name,
		arg.CreatedAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const dahuaCreateDevice = `-- name: DahuaCreateDevice :one
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
  (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id
`

type DahuaCreateDeviceParams struct {
	Name      string
	Url       types.URL
	Ip        string
	Username  string
	Password  string
	Location  types.Location
	Feature   models.DahuaFeature
	CreatedAt types.Time
	UpdatedAt types.Time
}

func (q *Queries) DahuaCreateDevice(ctx context.Context, arg DahuaCreateDeviceParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, dahuaCreateDevice,
		arg.Name,
		arg.Url,
		arg.Ip,
		arg.Username,
		arg.Password,
		arg.Location,
		arg.Feature,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const dahuaCreateEmailAttachment = `-- name: DahuaCreateEmailAttachment :one
INSERT INTO
  dahua_email_attachments (message_id, file_name)
VALUES
  (?, ?) RETURNING id, message_id, file_name
`

type DahuaCreateEmailAttachmentParams struct {
	MessageID int64
	FileName  string
}

func (q *Queries) DahuaCreateEmailAttachment(ctx context.Context, arg DahuaCreateEmailAttachmentParams) (DahuaEmailAttachment, error) {
	row := q.db.QueryRowContext(ctx, dahuaCreateEmailAttachment, arg.MessageID, arg.FileName)
	var i DahuaEmailAttachment
	err := row.Scan(&i.ID, &i.MessageID, &i.FileName)
	return i, err
}

const dahuaCreateEmailMessage = `-- name: DahuaCreateEmailMessage :one
INSERT INTO
  dahua_email_messages (
    device_id,
    date,
    'from',
    ` + "`" + `to` + "`" + `,
    subject,
    ` + "`" + `text` + "`" + `,
    alarm_event,
    alarm_input_channel,
    alarm_name,
    created_at
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id, device_id, date, 'from', ` + "`" + `to` + "`" + `, subject, ` + "`" + `text` + "`" + `, alarm_event, alarm_input_channel, alarm_name, created_at
`

type DahuaCreateEmailMessageParams struct {
	DeviceID          int64
	Date              types.Time
	From              string
	To                types.StringSlice
	Subject           string
	Text              string
	AlarmEvent        string
	AlarmInputChannel int64
	AlarmName         string
	CreatedAt         types.Time
}

func (q *Queries) DahuaCreateEmailMessage(ctx context.Context, arg DahuaCreateEmailMessageParams) (DahuaEmailMessage, error) {
	row := q.db.QueryRowContext(ctx, dahuaCreateEmailMessage,
		arg.DeviceID,
		arg.Date,
		arg.From,
		arg.To,
		arg.Subject,
		arg.Text,
		arg.AlarmEvent,
		arg.AlarmInputChannel,
		arg.AlarmName,
		arg.CreatedAt,
	)
	var i DahuaEmailMessage
	err := row.Scan(
		&i.ID,
		&i.DeviceID,
		&i.Date,
		&i.From,
		&i.To,
		&i.Subject,
		&i.Text,
		&i.AlarmEvent,
		&i.AlarmInputChannel,
		&i.AlarmName,
		&i.CreatedAt,
	)
	return i, err
}

const dahuaCreateEvent = `-- name: DahuaCreateEvent :one
INSERT INTO
  dahua_events (
    device_id,
    code,
    action,
    ` + "`" + `index` + "`" + `,
    data,
    created_at
  )
VALUES
  (?, ?, ?, ?, ?, ?) RETURNING id
`

type DahuaCreateEventParams struct {
	DeviceID  int64
	Code      string
	Action    string
	Index     int64
	Data      json.RawMessage
	CreatedAt types.Time
}

func (q *Queries) DahuaCreateEvent(ctx context.Context, arg DahuaCreateEventParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, dahuaCreateEvent,
		arg.DeviceID,
		arg.Code,
		arg.Action,
		arg.Index,
		arg.Data,
		arg.CreatedAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const dahuaCreateEventRule = `-- name: DahuaCreateEventRule :exec
INSERT INTO
  dahua_event_rules (code, ignore_db, ignore_live, ignore_mqtt)
VALUES
  (?, ?, ?, ?)
`

type DahuaCreateEventRuleParams struct {
	Code       string
	IgnoreDb   bool
	IgnoreLive bool
	IgnoreMqtt bool
}

func (q *Queries) DahuaCreateEventRule(ctx context.Context, arg DahuaCreateEventRuleParams) error {
	_, err := q.db.ExecContext(ctx, dahuaCreateEventRule,
		arg.Code,
		arg.IgnoreDb,
		arg.IgnoreLive,
		arg.IgnoreMqtt,
	)
	return err
}

const dahuaCreateEventWorkerState = `-- name: DahuaCreateEventWorkerState :exec
INSERT INTO
  dahua_event_worker_states (device_id, state, error, created_at)
VALUES
  (?, ?, ?, ?)
`

type DahuaCreateEventWorkerStateParams struct {
	DeviceID  int64
	State     models.DahuaEventWorkerState
	Error     sql.NullString
	CreatedAt types.Time
}

func (q *Queries) DahuaCreateEventWorkerState(ctx context.Context, arg DahuaCreateEventWorkerStateParams) error {
	_, err := q.db.ExecContext(ctx, dahuaCreateEventWorkerState,
		arg.DeviceID,
		arg.State,
		arg.Error,
		arg.CreatedAt,
	)
	return err
}

const dahuaCreateFile = `-- name: DahuaCreateFile :one
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
  id = id RETURNING id
`

type DahuaCreateFileParams struct {
	DeviceID    int64
	Channel     int64
	StartTime   types.Time
	EndTime     types.Time
	Length      int64
	Type        string
	FilePath    string
	Duration    int64
	Disk        int64
	VideoStream string
	Flags       types.StringSlice
	Events      types.StringSlice
	Cluster     int64
	Partition   int64
	PicIndex    int64
	Repeat      int64
	WorkDir     string
	WorkDirSn   bool
	UpdatedAt   types.Time
	Storage     models.Storage
}

func (q *Queries) DahuaCreateFile(ctx context.Context, arg DahuaCreateFileParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, dahuaCreateFile,
		arg.DeviceID,
		arg.Channel,
		arg.StartTime,
		arg.EndTime,
		arg.Length,
		arg.Type,
		arg.FilePath,
		arg.Duration,
		arg.Disk,
		arg.VideoStream,
		arg.Flags,
		arg.Events,
		arg.Cluster,
		arg.Partition,
		arg.PicIndex,
		arg.Repeat,
		arg.WorkDir,
		arg.WorkDirSn,
		arg.UpdatedAt,
		arg.Storage,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const dahuaCreateFileCursor = `-- name: DahuaCreateFileCursor :exec
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
  (?, ?, ?, ?, ?, ?, ?)
`

type DahuaCreateFileCursorParams struct {
	DeviceID    int64
	QuickCursor types.Time
	FullCursor  types.Time
	FullEpoch   types.Time
	Scan        bool
	ScanPercent float64
	ScanType    models.DahuaScanType
}

func (q *Queries) DahuaCreateFileCursor(ctx context.Context, arg DahuaCreateFileCursorParams) error {
	_, err := q.db.ExecContext(ctx, dahuaCreateFileCursor,
		arg.DeviceID,
		arg.QuickCursor,
		arg.FullCursor,
		arg.FullEpoch,
		arg.Scan,
		arg.ScanPercent,
		arg.ScanType,
	)
	return err
}

const dahuaCreateStorageDestination = `-- name: DahuaCreateStorageDestination :one
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
  (?, ?, ?, ?, ?, ?, ?) RETURNING id
`

type DahuaCreateStorageDestinationParams struct {
	Name            string
	Storage         models.Storage
	ServerAddress   string
	Port            int64
	Username        string
	Password        string
	RemoteDirectory string
}

func (q *Queries) DahuaCreateStorageDestination(ctx context.Context, arg DahuaCreateStorageDestinationParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, dahuaCreateStorageDestination,
		arg.Name,
		arg.Storage,
		arg.ServerAddress,
		arg.Port,
		arg.Username,
		arg.Password,
		arg.RemoteDirectory,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const dahuaCreateStreamForInternal = `-- name: DahuaCreateStreamForInternal :one
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
  internal = true RETURNING ID
`

type DahuaCreateStreamForInternalParams struct {
	DeviceID     int64
	Channel      int64
	Subtype      int64
	Name         string
	MediamtxPath string
}

func (q *Queries) DahuaCreateStreamForInternal(ctx context.Context, arg DahuaCreateStreamForInternalParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, dahuaCreateStreamForInternal,
		arg.DeviceID,
		arg.Channel,
		arg.Subtype,
		arg.Name,
		arg.MediamtxPath,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const dahuaCreateThumbnail = `-- name: DahuaCreateThumbnail :one
INSERT INTO
  dahua_thumbnails (file_id, email_attachment_id, width, height)
VALUES
  (?, ?, ?, ?) RETURNING id, file_id, email_attachment_id, width, height
`

type DahuaCreateThumbnailParams struct {
	FileID            sql.NullInt64
	EmailAttachmentID sql.NullInt64
	Width             int64
	Height            int64
}

func (q *Queries) DahuaCreateThumbnail(ctx context.Context, arg DahuaCreateThumbnailParams) (DahuaThumbnail, error) {
	row := q.db.QueryRowContext(ctx, dahuaCreateThumbnail,
		arg.FileID,
		arg.EmailAttachmentID,
		arg.Width,
		arg.Height,
	)
	var i DahuaThumbnail
	err := row.Scan(
		&i.ID,
		&i.FileID,
		&i.EmailAttachmentID,
		&i.Width,
		&i.Height,
	)
	return i, err
}

const dahuaDeleteAferoFile = `-- name: DahuaDeleteAferoFile :exec
DELETE FROM dahua_afero_files
WHERE
  id = ?
`

func (q *Queries) DahuaDeleteAferoFile(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, dahuaDeleteAferoFile, id)
	return err
}

const dahuaDeleteDevice = `-- name: DahuaDeleteDevice :exec
DELETE FROM dahua_devices
WHERE
  id = ?
`

func (q *Queries) DahuaDeleteDevice(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, dahuaDeleteDevice, id)
	return err
}

const dahuaDeleteEvent = `-- name: DahuaDeleteEvent :exec
DELETE FROM dahua_events
`

func (q *Queries) DahuaDeleteEvent(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, dahuaDeleteEvent)
	return err
}

const dahuaDeleteEventRule = `-- name: DahuaDeleteEventRule :exec
DELETE FROM dahua_event_rules
WHERE
  id = ?
`

func (q *Queries) DahuaDeleteEventRule(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, dahuaDeleteEventRule, id)
	return err
}

const dahuaDeleteFile = `-- name: DahuaDeleteFile :exec
DELETE FROM dahua_files
WHERE
  updated_at < ?1
  AND device_id = ?2
  AND start_time <= ?3
  AND ?4 < start_time
`

type DahuaDeleteFileParams struct {
	UpdatedAt types.Time
	DeviceID  int64
	End       types.Time
	Start     types.Time
}

func (q *Queries) DahuaDeleteFile(ctx context.Context, arg DahuaDeleteFileParams) error {
	_, err := q.db.ExecContext(ctx, dahuaDeleteFile,
		arg.UpdatedAt,
		arg.DeviceID,
		arg.End,
		arg.Start,
	)
	return err
}

const dahuaDeleteStorageDestination = `-- name: DahuaDeleteStorageDestination :exec
DELETE FROM dahua_storage_destinations
WHERE
  id = ?
`

func (q *Queries) DahuaDeleteStorageDestination(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, dahuaDeleteStorageDestination, id)
	return err
}

const dahuaDeleteStream = `-- name: DahuaDeleteStream :exec
DELETE FROM dahua_streams
WHERE
  id = ?
`

func (q *Queries) DahuaDeleteStream(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, dahuaDeleteStream, id)
	return err
}

const dahuaGetAferoFileByFileID = `-- name: DahuaGetAferoFileByFileID :one
SELECT
  id, file_id, thumbnail_id, email_attachment_id, name, ready, size, created_at
FROM
  dahua_afero_files
WHERE
  file_id = ?
`

func (q *Queries) DahuaGetAferoFileByFileID(ctx context.Context, fileID sql.NullInt64) (DahuaAferoFile, error) {
	row := q.db.QueryRowContext(ctx, dahuaGetAferoFileByFileID, fileID)
	var i DahuaAferoFile
	err := row.Scan(
		&i.ID,
		&i.FileID,
		&i.ThumbnailID,
		&i.EmailAttachmentID,
		&i.Name,
		&i.Ready,
		&i.Size,
		&i.CreatedAt,
	)
	return i, err
}

const dahuaGetDeviceName = `-- name: DahuaGetDeviceName :one
SELECT
  name
FROM
  dahua_devices
WHERE
  id = ?
`

func (q *Queries) DahuaGetDeviceName(ctx context.Context, id int64) (string, error) {
	row := q.db.QueryRowContext(ctx, dahuaGetDeviceName, id)
	var name string
	err := row.Scan(&name)
	return name, err
}

const dahuaGetEventData = `-- name: DahuaGetEventData :one
SELECT
  data
FROM
  dahua_events
WHERE
  id = ?
`

func (q *Queries) DahuaGetEventData(ctx context.Context, id int64) (json.RawMessage, error) {
	row := q.db.QueryRowContext(ctx, dahuaGetEventData, id)
	var data json.RawMessage
	err := row.Scan(&data)
	return data, err
}

const dahuaGetEventRule = `-- name: DahuaGetEventRule :one
SELECT
  id, code, ignore_db, ignore_live, ignore_mqtt
FROM
  dahua_event_rules
WHERE
  id = ?
`

func (q *Queries) DahuaGetEventRule(ctx context.Context, id int64) (DahuaEventRule, error) {
	row := q.db.QueryRowContext(ctx, dahuaGetEventRule, id)
	var i DahuaEventRule
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.IgnoreDb,
		&i.IgnoreLive,
		&i.IgnoreMqtt,
	)
	return i, err
}

const dahuaGetEventRuleByEvent = `-- name: DahuaGetEventRuleByEvent :many
SELECT
  ignore_db,
  ignore_live,
  ignore_mqtt,
  code
FROM
  dahua_event_device_rules
WHERE
  device_id = ?1
  AND (
    dahua_event_device_rules.code = ?2
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
  dahua_event_rules.code = ?2
  OR dahua_event_rules.code = ''
ORDER BY
  code DESC
`

type DahuaGetEventRuleByEventParams struct {
	DeviceID int64
	Code     string
}

type DahuaGetEventRuleByEventRow struct {
	IgnoreDb   bool
	IgnoreLive bool
	IgnoreMqtt bool
	Code       string
}

func (q *Queries) DahuaGetEventRuleByEvent(ctx context.Context, arg DahuaGetEventRuleByEventParams) ([]DahuaGetEventRuleByEventRow, error) {
	rows, err := q.db.QueryContext(ctx, dahuaGetEventRuleByEvent, arg.DeviceID, arg.Code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DahuaGetEventRuleByEventRow
	for rows.Next() {
		var i DahuaGetEventRuleByEventRow
		if err := rows.Scan(
			&i.IgnoreDb,
			&i.IgnoreLive,
			&i.IgnoreMqtt,
			&i.Code,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dahuaGetFile = `-- name: DahuaGetFile :one
SELECT
  id, device_id, channel, start_time, end_time, length, type, file_path, duration, disk, video_stream, flags, events, cluster, "partition", pic_index, repeat, work_dir, work_dir_sn, updated_at, storage
FROM
  dahua_files
WHERE
  id = ?
`

func (q *Queries) DahuaGetFile(ctx context.Context, id int64) (DahuaFile, error) {
	row := q.db.QueryRowContext(ctx, dahuaGetFile, id)
	var i DahuaFile
	err := row.Scan(
		&i.ID,
		&i.DeviceID,
		&i.Channel,
		&i.StartTime,
		&i.EndTime,
		&i.Length,
		&i.Type,
		&i.FilePath,
		&i.Duration,
		&i.Disk,
		&i.VideoStream,
		&i.Flags,
		&i.Events,
		&i.Cluster,
		&i.Partition,
		&i.PicIndex,
		&i.Repeat,
		&i.WorkDir,
		&i.WorkDirSn,
		&i.UpdatedAt,
		&i.Storage,
	)
	return i, err
}

const dahuaGetFileByFilePath = `-- name: DahuaGetFileByFilePath :one
SELECT
  id, device_id, channel, start_time, end_time, length, type, file_path, duration, disk, video_stream, flags, events, cluster, "partition", pic_index, repeat, work_dir, work_dir_sn, updated_at, storage
FROM
  dahua_files
WHERE
  device_id = ?
  and file_path = ?
`

type DahuaGetFileByFilePathParams struct {
	DeviceID int64
	FilePath string
}

func (q *Queries) DahuaGetFileByFilePath(ctx context.Context, arg DahuaGetFileByFilePathParams) (DahuaFile, error) {
	row := q.db.QueryRowContext(ctx, dahuaGetFileByFilePath, arg.DeviceID, arg.FilePath)
	var i DahuaFile
	err := row.Scan(
		&i.ID,
		&i.DeviceID,
		&i.Channel,
		&i.StartTime,
		&i.EndTime,
		&i.Length,
		&i.Type,
		&i.FilePath,
		&i.Duration,
		&i.Disk,
		&i.VideoStream,
		&i.Flags,
		&i.Events,
		&i.Cluster,
		&i.Partition,
		&i.PicIndex,
		&i.Repeat,
		&i.WorkDir,
		&i.WorkDirSn,
		&i.UpdatedAt,
		&i.Storage,
	)
	return i, err
}

const dahuaGetFileForThumbnail = `-- name: DahuaGetFileForThumbnail :one
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
  dahua_files.id = ?
`

type DahuaGetFileForThumbnailRow struct {
	ID       int64
	DeviceID int64
	Type     string
	FilePath string
	Name     sql.NullString
	Ready    sql.NullBool
}

func (q *Queries) DahuaGetFileForThumbnail(ctx context.Context, id int64) (DahuaGetFileForThumbnailRow, error) {
	row := q.db.QueryRowContext(ctx, dahuaGetFileForThumbnail, id)
	var i DahuaGetFileForThumbnailRow
	err := row.Scan(
		&i.ID,
		&i.DeviceID,
		&i.Type,
		&i.FilePath,
		&i.Name,
		&i.Ready,
	)
	return i, err
}

const dahuaGetOldestFileStartTime = `-- name: DahuaGetOldestFileStartTime :one
SELECT
  start_time
FROM
  dahua_files
WHERE
  device_id = ?
ORDER BY
  start_time ASC
LIMIT
  1
`

func (q *Queries) DahuaGetOldestFileStartTime(ctx context.Context, deviceID int64) (types.Time, error) {
	row := q.db.QueryRowContext(ctx, dahuaGetOldestFileStartTime, deviceID)
	var start_time types.Time
	err := row.Scan(&start_time)
	return start_time, err
}

const dahuaGetStorageDestination = `-- name: DahuaGetStorageDestination :one
SELECT
  id, name, storage, server_address, port, username, password, remote_directory
FROM
  dahua_storage_destinations
WHERE
  id = ?
`

func (q *Queries) DahuaGetStorageDestination(ctx context.Context, id int64) (DahuaStorageDestination, error) {
	row := q.db.QueryRowContext(ctx, dahuaGetStorageDestination, id)
	var i DahuaStorageDestination
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Storage,
		&i.ServerAddress,
		&i.Port,
		&i.Username,
		&i.Password,
		&i.RemoteDirectory,
	)
	return i, err
}

const dahuaGetStorageDestinationByServerAddressAndStorage = `-- name: DahuaGetStorageDestinationByServerAddressAndStorage :one
SELECT
  id, name, storage, server_address, port, username, password, remote_directory
FROM
  dahua_storage_destinations
WHERE
  server_address = ?
  AND storage = ?
`

type DahuaGetStorageDestinationByServerAddressAndStorageParams struct {
	ServerAddress string
	Storage       models.Storage
}

func (q *Queries) DahuaGetStorageDestinationByServerAddressAndStorage(ctx context.Context, arg DahuaGetStorageDestinationByServerAddressAndStorageParams) (DahuaStorageDestination, error) {
	row := q.db.QueryRowContext(ctx, dahuaGetStorageDestinationByServerAddressAndStorage, arg.ServerAddress, arg.Storage)
	var i DahuaStorageDestination
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Storage,
		&i.ServerAddress,
		&i.Port,
		&i.Username,
		&i.Password,
		&i.RemoteDirectory,
	)
	return i, err
}

const dahuaGetStream = `-- name: DahuaGetStream :one
SELECT
  id, internal, device_id, channel, subtype, name, mediamtx_path
FROM
  dahua_streams
WHERE
  id = ?
`

func (q *Queries) DahuaGetStream(ctx context.Context, id int64) (DahuaStream, error) {
	row := q.db.QueryRowContext(ctx, dahuaGetStream, id)
	var i DahuaStream
	err := row.Scan(
		&i.ID,
		&i.Internal,
		&i.DeviceID,
		&i.Channel,
		&i.Subtype,
		&i.Name,
		&i.MediamtxPath,
	)
	return i, err
}

const dahuaListDevicePermissions = `-- name: DahuaListDevicePermissions :many
SELECT
  dahua_devices.id,
  coalesce(p.level, 2)
FROM
  dahua_devices
  LEFT JOIN dahua_permissions AS p ON p.device_id = dahua_devices.id
WHERE
  -- Allow if user is admin
  EXISTS (SELECT user_id FROM admins WHERE admins.user_id = ?1)
  -- Allow if user is a part of the group the owns the permission
  OR p.group_id IN (
    SELECT
      group_id
    FROM
      group_users
    WHERE
      group_users.user_id = ?1
  )
  -- Allow if user owns the permission
  OR p.user_id = ?1
GROUP BY
  -- Remove duplicate devices with different permissions
  dahua_devices.id
ORDER BY
  -- Get the highest permission level
  p.level DESC
`

type DahuaListDevicePermissionsRow struct {
	ID    int64
	Level models.DahuaPermissionLevel
}

func (q *Queries) DahuaListDevicePermissions(ctx context.Context, userID int64) ([]DahuaListDevicePermissionsRow, error) {
	rows, err := q.db.QueryContext(ctx, dahuaListDevicePermissions, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DahuaListDevicePermissionsRow
	for rows.Next() {
		var i DahuaListDevicePermissionsRow
		if err := rows.Scan(&i.ID, &i.Level); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dahuaListEventActions = `-- name: DahuaListEventActions :many
SELECT DISTINCT
  action
FROM
  dahua_events
`

func (q *Queries) DahuaListEventActions(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, dahuaListEventActions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var action string
		if err := rows.Scan(&action); err != nil {
			return nil, err
		}
		items = append(items, action)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dahuaListEventCodes = `-- name: DahuaListEventCodes :many
SELECT DISTINCT
  code
FROM
  dahua_events
`

func (q *Queries) DahuaListEventCodes(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, dahuaListEventCodes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		items = append(items, code)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dahuaListEventRules = `-- name: DahuaListEventRules :many
SELECT
  id, code, ignore_db, ignore_live, ignore_mqtt
FROM
  dahua_event_rules
`

func (q *Queries) DahuaListEventRules(ctx context.Context) ([]DahuaEventRule, error) {
	rows, err := q.db.QueryContext(ctx, dahuaListEventRules)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DahuaEventRule
	for rows.Next() {
		var i DahuaEventRule
		if err := rows.Scan(
			&i.ID,
			&i.Code,
			&i.IgnoreDb,
			&i.IgnoreLive,
			&i.IgnoreMqtt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dahuaListEventWorkerStates = `-- name: DahuaListEventWorkerStates :many
SELECT
  id, device_id, state, error, created_at,
  max(created_at)
FROM
  dahua_event_worker_states
GROUP BY
  device_id
`

type DahuaListEventWorkerStatesRow struct {
	ID        int64
	DeviceID  int64
	State     models.DahuaEventWorkerState
	Error     sql.NullString
	CreatedAt types.Time
	Max       interface{}
}

func (q *Queries) DahuaListEventWorkerStates(ctx context.Context) ([]DahuaListEventWorkerStatesRow, error) {
	rows, err := q.db.QueryContext(ctx, dahuaListEventWorkerStates)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DahuaListEventWorkerStatesRow
	for rows.Next() {
		var i DahuaListEventWorkerStatesRow
		if err := rows.Scan(
			&i.ID,
			&i.DeviceID,
			&i.State,
			&i.Error,
			&i.CreatedAt,
			&i.Max,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dahuaListFileCursors = `-- name: DahuaListFileCursors :many
SELECT
  c.device_id, c.quick_cursor, c.full_cursor, c.full_epoch, c.full_complete, c.scan, c.scan_percent, c.scan_type,
  count(f.device_id) AS files
FROM
  dahua_file_cursors AS c
  LEFT JOIN dahua_files AS f ON f.device_id = c.device_id
GROUP BY
  c.device_id
`

type DahuaListFileCursorsRow struct {
	DeviceID     int64
	QuickCursor  types.Time
	FullCursor   types.Time
	FullEpoch    types.Time
	FullComplete bool
	Scan         bool
	ScanPercent  float64
	ScanType     models.DahuaScanType
	Files        int64
}

func (q *Queries) DahuaListFileCursors(ctx context.Context) ([]DahuaListFileCursorsRow, error) {
	rows, err := q.db.QueryContext(ctx, dahuaListFileCursors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DahuaListFileCursorsRow
	for rows.Next() {
		var i DahuaListFileCursorsRow
		if err := rows.Scan(
			&i.DeviceID,
			&i.QuickCursor,
			&i.FullCursor,
			&i.FullEpoch,
			&i.FullComplete,
			&i.Scan,
			&i.ScanPercent,
			&i.ScanType,
			&i.Files,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dahuaListFileTypes = `-- name: DahuaListFileTypes :many
SELECT DISTINCT
  type
FROM
  dahua_files
`

func (q *Queries) DahuaListFileTypes(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, dahuaListFileTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var type_ string
		if err := rows.Scan(&type_); err != nil {
			return nil, err
		}
		items = append(items, type_)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dahuaListStorageDestinations = `-- name: DahuaListStorageDestinations :many
SELECT
  id, name, storage, server_address, port, username, password, remote_directory
FROM
  dahua_storage_destinations
`

func (q *Queries) DahuaListStorageDestinations(ctx context.Context) ([]DahuaStorageDestination, error) {
	rows, err := q.db.QueryContext(ctx, dahuaListStorageDestinations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DahuaStorageDestination
	for rows.Next() {
		var i DahuaStorageDestination
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Storage,
			&i.ServerAddress,
			&i.Port,
			&i.Username,
			&i.Password,
			&i.RemoteDirectory,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dahuaListStreams = `-- name: DahuaListStreams :many
SELECT
  id, internal, device_id, channel, subtype, name, mediamtx_path
FROM
  dahua_streams
ORDER BY
  device_id
`

func (q *Queries) DahuaListStreams(ctx context.Context) ([]DahuaStream, error) {
	rows, err := q.db.QueryContext(ctx, dahuaListStreams)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DahuaStream
	for rows.Next() {
		var i DahuaStream
		if err := rows.Scan(
			&i.ID,
			&i.Internal,
			&i.DeviceID,
			&i.Channel,
			&i.Subtype,
			&i.Name,
			&i.MediamtxPath,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dahuaListStreamsByDevice = `-- name: DahuaListStreamsByDevice :many
SELECT
  id, internal, device_id, channel, subtype, name, mediamtx_path
FROM
  dahua_streams
WHERE
  device_id = ?
`

func (q *Queries) DahuaListStreamsByDevice(ctx context.Context, deviceID int64) ([]DahuaStream, error) {
	rows, err := q.db.QueryContext(ctx, dahuaListStreamsByDevice, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DahuaStream
	for rows.Next() {
		var i DahuaStream
		if err := rows.Scan(
			&i.ID,
			&i.Internal,
			&i.DeviceID,
			&i.Channel,
			&i.Subtype,
			&i.Name,
			&i.MediamtxPath,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dahuaNormalizeFileCursors = `-- name: DahuaNormalizeFileCursors :exec
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
FROM dahua_devices
`

type DahuaNormalizeFileCursorsParams struct {
	QuickCursor types.Time
	FullCursor  types.Time
	FullEpoch   types.Time
	Scan        bool
	ScanPercent float64
	ScanType    models.DahuaScanType
}

func (q *Queries) DahuaNormalizeFileCursors(ctx context.Context, arg DahuaNormalizeFileCursorsParams) error {
	_, err := q.db.ExecContext(ctx, dahuaNormalizeFileCursors,
		arg.QuickCursor,
		arg.FullCursor,
		arg.FullEpoch,
		arg.Scan,
		arg.ScanPercent,
		arg.ScanType,
	)
	return err
}

const dahuaOrphanDeleteThumbnail = `-- name: DahuaOrphanDeleteThumbnail :exec
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
  )
`

func (q *Queries) DahuaOrphanDeleteThumbnail(ctx context.Context, createdAt types.Time) error {
	_, err := q.db.ExecContext(ctx, dahuaOrphanDeleteThumbnail, createdAt)
	return err
}

const dahuaOrphanListAferoFiles = `-- name: DahuaOrphanListAferoFiles :many
SELECT
  id, file_id, thumbnail_id, email_attachment_id, name, ready, size, created_at
FROM
  dahua_afero_files
WHERE
  file_id IS NULL
  AND thumbnail_id IS NULL
  AND email_attachment_id IS NULL
  AND ready = true
LIMIT
  ?
`

func (q *Queries) DahuaOrphanListAferoFiles(ctx context.Context, limit int64) ([]DahuaAferoFile, error) {
	rows, err := q.db.QueryContext(ctx, dahuaOrphanListAferoFiles, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DahuaAferoFile
	for rows.Next() {
		var i DahuaAferoFile
		if err := rows.Scan(
			&i.ID,
			&i.FileID,
			&i.ThumbnailID,
			&i.EmailAttachmentID,
			&i.Name,
			&i.Ready,
			&i.Size,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const dahuaReadyAferoFile = `-- name: DahuaReadyAferoFile :one
UPDATE dahua_afero_files
SET
  ready = true,
  size = ?,
  created_at = ?
WHERE
  id = ? RETURNING id
`

type DahuaReadyAferoFileParams struct {
	Size      int64
	CreatedAt types.Time
	ID        int64
}

func (q *Queries) DahuaReadyAferoFile(ctx context.Context, arg DahuaReadyAferoFileParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, dahuaReadyAferoFile, arg.Size, arg.CreatedAt, arg.ID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const dahuaUpdateDevice = `-- name: DahuaUpdateDevice :one
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
  id = ? RETURNING id
`

type DahuaUpdateDeviceParams struct {
	Name      string
	Url       types.URL
	Ip        string
	Username  string
	Password  string
	Location  types.Location
	Feature   models.DahuaFeature
	UpdatedAt types.Time
	ID        int64
}

func (q *Queries) DahuaUpdateDevice(ctx context.Context, arg DahuaUpdateDeviceParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, dahuaUpdateDevice,
		arg.Name,
		arg.Url,
		arg.Ip,
		arg.Username,
		arg.Password,
		arg.Location,
		arg.Feature,
		arg.UpdatedAt,
		arg.ID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const dahuaUpdateEventRule = `-- name: DahuaUpdateEventRule :exec
UPDATE dahua_event_rules
SET
  code = ?,
  ignore_db = ?,
  ignore_live = ?,
  ignore_mqtt = ?
WHERE
  id = ?
`

type DahuaUpdateEventRuleParams struct {
	Code       string
	IgnoreDb   bool
	IgnoreLive bool
	IgnoreMqtt bool
	ID         int64
}

func (q *Queries) DahuaUpdateEventRule(ctx context.Context, arg DahuaUpdateEventRuleParams) error {
	_, err := q.db.ExecContext(ctx, dahuaUpdateEventRule,
		arg.Code,
		arg.IgnoreDb,
		arg.IgnoreLive,
		arg.IgnoreMqtt,
		arg.ID,
	)
	return err
}

const dahuaUpdateFile = `-- name: DahuaUpdateFile :one
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
  AND file_path = ? RETURNING id
`

type DahuaUpdateFileParams struct {
	Channel     int64
	StartTime   types.Time
	EndTime     types.Time
	Length      int64
	Type        string
	Duration    int64
	Disk        int64
	VideoStream string
	Flags       types.StringSlice
	Events      types.StringSlice
	Cluster     int64
	Partition   int64
	PicIndex    int64
	Repeat      int64
	WorkDir     string
	WorkDirSn   bool
	UpdatedAt   types.Time
	Storage     models.Storage
	DeviceID    int64
	FilePath    string
}

func (q *Queries) DahuaUpdateFile(ctx context.Context, arg DahuaUpdateFileParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, dahuaUpdateFile,
		arg.Channel,
		arg.StartTime,
		arg.EndTime,
		arg.Length,
		arg.Type,
		arg.Duration,
		arg.Disk,
		arg.VideoStream,
		arg.Flags,
		arg.Events,
		arg.Cluster,
		arg.Partition,
		arg.PicIndex,
		arg.Repeat,
		arg.WorkDir,
		arg.WorkDirSn,
		arg.UpdatedAt,
		arg.Storage,
		arg.DeviceID,
		arg.FilePath,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const dahuaUpdateFileCursor = `-- name: DahuaUpdateFileCursor :one
UPDATE dahua_file_cursors
SET
  quick_cursor = ?,
  full_cursor = ?,
  full_epoch = ?,
  scan = ?,
  scan_percent = ?,
  scan_type = ?
WHERE
  device_id = ? RETURNING device_id, quick_cursor, full_cursor, full_epoch, full_complete, scan, scan_percent, scan_type
`

type DahuaUpdateFileCursorParams struct {
	QuickCursor types.Time
	FullCursor  types.Time
	FullEpoch   types.Time
	Scan        bool
	ScanPercent float64
	ScanType    models.DahuaScanType
	DeviceID    int64
}

func (q *Queries) DahuaUpdateFileCursor(ctx context.Context, arg DahuaUpdateFileCursorParams) (DahuaFileCursor, error) {
	row := q.db.QueryRowContext(ctx, dahuaUpdateFileCursor,
		arg.QuickCursor,
		arg.FullCursor,
		arg.FullEpoch,
		arg.Scan,
		arg.ScanPercent,
		arg.ScanType,
		arg.DeviceID,
	)
	var i DahuaFileCursor
	err := row.Scan(
		&i.DeviceID,
		&i.QuickCursor,
		&i.FullCursor,
		&i.FullEpoch,
		&i.FullComplete,
		&i.Scan,
		&i.ScanPercent,
		&i.ScanType,
	)
	return i, err
}

const dahuaUpdateFileCursorScanPercent = `-- name: DahuaUpdateFileCursorScanPercent :one
UPDATE dahua_file_cursors
SET
  scan_percent = ?
WHERE
  device_id = ? RETURNING device_id, quick_cursor, full_cursor, full_epoch, full_complete, scan, scan_percent, scan_type
`

type DahuaUpdateFileCursorScanPercentParams struct {
	ScanPercent float64
	DeviceID    int64
}

func (q *Queries) DahuaUpdateFileCursorScanPercent(ctx context.Context, arg DahuaUpdateFileCursorScanPercentParams) (DahuaFileCursor, error) {
	row := q.db.QueryRowContext(ctx, dahuaUpdateFileCursorScanPercent, arg.ScanPercent, arg.DeviceID)
	var i DahuaFileCursor
	err := row.Scan(
		&i.DeviceID,
		&i.QuickCursor,
		&i.FullCursor,
		&i.FullEpoch,
		&i.FullComplete,
		&i.Scan,
		&i.ScanPercent,
		&i.ScanType,
	)
	return i, err
}

const dahuaUpdateStorageDestination = `-- name: DahuaUpdateStorageDestination :one
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
  id = ? RETURNING id
`

type DahuaUpdateStorageDestinationParams struct {
	Name            string
	Storage         models.Storage
	ServerAddress   string
	Port            int64
	Username        string
	Password        string
	RemoteDirectory string
	ID              int64
}

func (q *Queries) DahuaUpdateStorageDestination(ctx context.Context, arg DahuaUpdateStorageDestinationParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, dahuaUpdateStorageDestination,
		arg.Name,
		arg.Storage,
		arg.ServerAddress,
		arg.Port,
		arg.Username,
		arg.Password,
		arg.RemoteDirectory,
		arg.ID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const dahuaUpdateStream = `-- name: DahuaUpdateStream :one
UPDATE dahua_streams
SET
  name = ?,
  mediamtx_path = ?
WHERE
  id = ? RETURNING id, internal, device_id, channel, subtype, name, mediamtx_path
`

type DahuaUpdateStreamParams struct {
	Name         string
	MediamtxPath string
	ID           int64
}

func (q *Queries) DahuaUpdateStream(ctx context.Context, arg DahuaUpdateStreamParams) (DahuaStream, error) {
	row := q.db.QueryRowContext(ctx, dahuaUpdateStream, arg.Name, arg.MediamtxPath, arg.ID)
	var i DahuaStream
	err := row.Scan(
		&i.ID,
		&i.Internal,
		&i.DeviceID,
		&i.Channel,
		&i.Subtype,
		&i.Name,
		&i.MediamtxPath,
	)
	return i, err
}

const dahuaUpdateStreamForInternal = `-- name: DahuaUpdateStreamForInternal :exec
UPDATE dahua_streams
SET
  internal = false
WHERE
  device_id = ?
`

func (q *Queries) DahuaUpdateStreamForInternal(ctx context.Context, deviceID int64) error {
	_, err := q.db.ExecContext(ctx, dahuaUpdateStreamForInternal, deviceID)
	return err
}
