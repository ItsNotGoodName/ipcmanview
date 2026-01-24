package dahua

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type ScanCursor struct {
	Device_ID     int64
	Quick_Cursor  types.Time
	Full_Cursor   types.Time
	Full_Epoch    types.Time
	Full_Complete bool
	Updated_At    types.Time
}

// ScanEpoch is the oldest a file can be.
var ScanEpoch time.Time = core.Must2(time.ParseInLocation(time.DateTime, "2009-12-31 00:00:00", time.UTC))

const (
	// scanVolatilePeriod is latest time period that the device could still be writing files to disk.
	scanVolatilePeriod = 8 * time.Hour
	// scanMaxPeriod is the maximum time period a device can handle when scanning files before they give weird results.
	scanMaxPeriod = 30 * 24 * time.Hour
)

func NewScanCursorQuick(end time.Time) time.Time {
	return core.Oldest(end, time.Now().Add(-scanVolatilePeriod))
}

func NewScanCursorFull(start, epoch time.Time) time.Time {
	return core.Newest(start, epoch)
}

func ResetScanCursor(ctx context.Context, db sqlx.ExecerContext, deviceID int64) error {
	now := time.Now()
	quickCursor := types.NewTime(now.Add(-scanVolatilePeriod))
	fullCursor := types.NewTime(now)
	fullEpoch := types.NewTime(ScanEpoch)
	updatedAt := types.NewTime(now)

	_, err := db.ExecContext(ctx, `
		INSERT OR REPLACE INTO dahua_scan_cursors (
			device_id,
			quick_cursor,
			full_cursor,
			full_epoch,
			updated_at
		) VALUES (?, ?, ?, ?, ?)
	`,
		deviceID,
		quickCursor,
		fullCursor,
		fullEpoch,
		updatedAt,
	)
	return err
}

func NewScanRange(start, end time.Time) (ScanRange, SubScan, bool) {
	scanRange := ScanRange{
		Start: start,
		End:   end,
	}
	scanCursor, ok := scanRange.NextSubScan(start)
	return scanRange, scanCursor, ok
}

type ScanRange struct {
	Start time.Time
	End   time.Time
}

type SubScan struct {
	ScanStart time.Time
	ScanEnd   time.Time
	Cursor    time.Time
}

func (r ScanRange) NextSubScan(cursor time.Time) (SubScan, bool) {
	if r.Start.Equal(r.End) || cursor.Equal(r.End) {
		return SubScan{}, false
	}

	if r.Start.Before(r.End) {
		nextCursor := core.Oldest(cursor.Add(scanMaxPeriod), r.End)
		return SubScan{
			ScanStart: cursor,
			ScanEnd:   nextCursor,
			Cursor:    nextCursor,
		}, true
	} else {
		nextCursor := core.Newest(cursor.Add(-scanMaxPeriod), r.End)
		return SubScan{
			ScanStart: cursor,
			ScanEnd:   nextCursor,
			Cursor:    nextCursor,
		}, true
	}
}

func (r ScanRange) Percent(cursor time.Time) float64 {
	top := cursor.Sub(r.Start).Abs().Hours()
	bottom := r.End.Sub(r.Start).Abs().Hours()
	percent := (top / bottom) * 100
	return percent
}

func ScanManual(ctx context.Context, db *sqlx.DB, conn dahuarpc.Conn, deviceID int64, start, end time.Time) (ScanResult, error) {
	var result ScanResult
	scanRange, scanCursor, ok := NewScanRange(start, end)
	for ok {
		scanResult, err := scan(ctx, db, conn, deviceID, scanCursor.ScanStart, scanCursor.ScanEnd)
		if err != nil {
			return ScanResult{}, err
		}
		result.add(scanResult)

		scanCursor, ok = scanRange.NextSubScan(scanCursor.Cursor)
	}
	return result, nil
}

func ScanFull(ctx context.Context, db *sqlx.DB, conn dahuarpc.Conn, deviceID int64) (ScanResult, error) {
	var fileScanCursor ScanCursor
	err := db.GetContext(ctx, &fileScanCursor, `
		SELECT * FROM dahua_scan_cursors WHERE device_id = ?
	`, deviceID)
	if err != nil {
		return ScanResult{}, err
	}

	var result ScanResult
	scanRange, scanCursor, ok := NewScanRange(fileScanCursor.Full_Cursor.Time, fileScanCursor.Full_Epoch.Time)
	for ok {
		scanResult, err := scan(ctx, db, conn, deviceID, scanCursor.ScanStart, scanCursor.ScanEnd)
		if err != nil {
			return ScanResult{}, err
		}
		result.add(scanResult)

		fullCursor := types.NewTime(NewScanCursorFull(scanCursor.ScanStart, fileScanCursor.Full_Epoch.Time))

		_, err = db.ExecContext(ctx, `
			UPDATE dahua_scan_cursors SET full_cursor = ? WHERE device_id = ?
		`, fullCursor, deviceID)
		if err != nil {
			return ScanResult{}, err
		}

		scanCursor, ok = scanRange.NextSubScan(scanCursor.Cursor)
	}

	_, err = db.ExecContext(ctx, `
		UPDATE dahua_scan_cursors SET full_cursor = full_epoch WHERE device_id = ?
	`, deviceID)
	if err != nil {
		return ScanResult{}, err
	}

	return result, nil
}

func ScanQuick(ctx context.Context, db *sqlx.DB, conn dahuarpc.Conn, deviceID int64) (ScanResult, error) {
	var fileScanCursor ScanCursor
	err := db.GetContext(ctx, &fileScanCursor, `
		SELECT * FROM dahua_scan_cursors WHERE device_id = ?
	`, deviceID)
	if err != nil {
		return ScanResult{}, err
	}

	var result ScanResult
	scanRange, scanCursor, ok := NewScanRange(fileScanCursor.Quick_Cursor.Time, time.Now())
	for ok {
		scanResult, err := scan(ctx, db, conn, deviceID, scanCursor.ScanStart, scanCursor.ScanEnd)
		if err != nil {
			return ScanResult{}, err
		}
		result.add(scanResult)

		quickCursor := types.NewTime(NewScanCursorQuick(scanCursor.ScanEnd))

		_, err = db.ExecContext(ctx, `
			UPDATE dahua_scan_cursors SET quick_cursor = ? WHERE device_id = ?
		`, quickCursor, deviceID)
		if err != nil {
			return ScanResult{}, err
		}

		scanCursor, ok = scanRange.NextSubScan(scanCursor.Cursor)
	}

	return result, nil
}

type ScanResult struct {
	CreatedCount int64 `json:"created_count"`
	UpdatedCount int64 `json:"updated_count"`
	DeletedCount int64 `json:"deleted_count"`
}

func (lhs *ScanResult) add(rhs ScanResult) {
	lhs.CreatedCount += rhs.CreatedCount
	lhs.UpdatedCount += rhs.UpdatedCount
	lhs.DeletedCount += rhs.DeletedCount
}

// scan files on a device by a time range.
func scan(ctx context.Context, db *sqlx.DB, conn dahuarpc.Conn, deviceID int64, scanStart, scanEnd time.Time) (ScanResult, error) {
	// Get device context
	var data struct {
		types.Key
		Location types.Location
		Seed     int64
		Name     string
	}
	err := db.GetContext(ctx, &data, `
		SELECT d.uuid, d.id, coalesce(d.location, s.location) AS location, d.seed, d.name
		FROM dahua_devices AS d, (SELECT value as location FROM settings WHERE key = ?) AS s
		WHERE d.id = ?
	`, system.KeyLocation, deviceID)
	if err != nil {
		return ScanResult{}, err
	}

	dahuaStart, dahuaEnd, dahuaOrder := scanStart, scanEnd, mediafilefind.ConditionOrderAscent
	if dahuaEnd.Before(dahuaStart) {
		dahuaStart, dahuaEnd, dahuaOrder = scanEnd, scanStart, mediafilefind.ConditionOrderDescent
	}

	condition := mediafilefind.NewCondtion(
		dahuarpc.NewTimestamp(dahuaStart, data.Location.Location),
		dahuarpc.NewTimestamp(dahuaEnd, data.Location.Location),
		dahuaOrder,
	)

	updatedAt := types.NewTime(time.Now())

	var createdCount int64
	var updatedCount int64
	var stream *mediafilefind.Stream
	defer func() {
		if stream != nil {
			stream.Close()
		}
	}()

	// For picture and video conditions
	for _, kind := range []string{"picture", "video"} {
		switch kind {
		case "picture":
			stream, err = mediafilefind.OpenStream(ctx, conn, condition.Picture())
			if err != nil {
				return ScanResult{}, err
			}
		case "video":
			stream, err = mediafilefind.OpenStream(ctx, conn, condition.Video())
			if err != nil {
				return ScanResult{}, err
			}
		default:
			panic(fmt.Sprintf("invalid kind %s", kind))
		}

		// Until stream is empty
		for {
			files, next, err := stream.Next(ctx)
			if err != nil {
				return ScanResult{}, err
			}
			if !next {
				break
			}

			// For each file
			for _, v := range files {
				startTime, endTime, err := v.UniqueTime(int(data.Seed), data.Location.Location)
				if err != nil {
					slog.Error("Failed to get unique time for file", "error", err, "device", data.Name)
					continue
				}
				storage := StorageFromFilePath(v.FilePath)
				events := v.CleanEvents()

				result, err := db.ExecContext(ctx, `
						UPDATE dahua_files SET 
							channel = ?,
							start_time = ?,
							end_time = ?,
							length = ?,
							type = ?,
							file_path = ?,
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
							storage = ?,
							updated_at = ?
						WHERE device_id = ? AND file_path = ?
					`,
					v.Channel,
					types.NewTime(startTime),
					types.NewTime(endTime),
					v.Length,
					v.Type,
					v.FilePath,
					v.Duration,
					v.Disk,
					v.VideoStream,
					types.NewSlice(v.Flags),
					types.NewSlice(events),
					v.Cluster,
					v.Partition,
					v.PicIndex,
					v.Repeat,
					v.WorkDir,
					v.WorkDirSN,
					storage,
					updatedAt,
					deviceID,
					v.FilePath,
				)
				if err != nil {
					return ScanResult{}, err
				}
				count, err := result.RowsAffected()
				if err != nil {
					return ScanResult{}, err
				}
				if count != 0 {
					updatedCount++
					continue
				}

				fileID := ulid.MustNew(ulid.Timestamp(startTime), ulid.DefaultEntropy()).String()

				_, err = db.ExecContext(ctx, `
						INSERT INTO dahua_files (
							id,
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
							storage,
							updated_at
						)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
						ON CONFLICT (start_time) DO UPDATE SET id = id RETURNING id -- Assume that files with the same time are really rare
					`,
					fileID,
					deviceID,
					v.Channel,
					types.NewTime(startTime),
					types.NewTime(endTime),
					v.Length,
					v.Type,
					v.FilePath,
					v.Duration,
					v.Disk,
					v.VideoStream,
					types.NewSlice(v.Flags),
					types.NewSlice(events),
					v.Cluster,
					v.Partition,
					v.PicIndex,
					v.Repeat,
					v.WorkDir,
					v.WorkDirSN,
					storage,
					updatedAt,
				)
				if err != nil {
					return ScanResult{}, err
				}

				createdCount++
			}
		}

		stream.Close()
	}

	// Delete stale files
	result, err := db.ExecContext(ctx, `
		DELETE FROM dahua_files
		WHERE
			device_id = ?
			AND ? < start_time
			AND start_time <= ?
			AND updated_at < ?
	`, deviceID, types.NewTime(scanStart), types.NewTime(scanEnd), updatedAt)
	if err != nil {
		return ScanResult{}, err
	}

	deletedCount, err := result.RowsAffected()
	if err != nil {
		return ScanResult{}, err
	}

	return ScanResult{
		CreatedCount: createdCount,
		UpdatedCount: updatedCount,
		DeletedCount: deletedCount,
	}, err
}
