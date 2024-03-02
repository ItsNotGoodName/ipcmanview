package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
)

func init() {
	var err error
	scanEpoch, err = time.ParseInLocation(time.DateTime, "2009-12-31 00:00:00", time.UTC)
	if err != nil {
		panic(err)
	}
}

// scanEpoch is the oldest time a file can exist.
var scanEpoch time.Time

const scanVolatileDuration = 8 * time.Hour

func newFileCursor() repo.DahuaCreateFileCursorParams {
	now := time.Now()
	return repo.DahuaCreateFileCursorParams{
		QuickCursor: types.NewTime(now.Add(-scanVolatileDuration)),
		FullCursor:  types.NewTime(now),
		FullEpoch:   types.NewTime(scanEpoch),
	}
}

func updateScanFileCursor(fileCursor repo.DahuaFileCursor, scanPeriod ScannerPeriod, scanType models.DahuaScanType) repo.DahuaFileCursor {
	switch scanType {
	case models.DahuaScanTypeFull:
		// Update FullCursor
		if scanPeriod.Start.Before(fileCursor.FullCursor.Time) {
			fileCursor.FullCursor = types.NewTime(scanPeriod.Start)
		}
	case models.DahuaScanTypeQuick:
		// Update QuickCursor
		quickCursor := time.Now().Add(-scanVolatileDuration)
		if scanPeriod.End.Before(quickCursor) {
			fileCursor.QuickCursor = types.NewTime(scanPeriod.End)
		} else {
			fileCursor.QuickCursor = types.NewTime(quickCursor)
		}
	case models.DahuaScanTypeReverse:
	default:
		panic("unknown type")
	}

	return fileCursor
}

func getScanRange(ctx context.Context, db sqlite.DB, fileCursor repo.DahuaFileCursor, scanType models.DahuaScanType) (models.TimeRange, error) {
	switch scanType {
	case models.DahuaScanTypeFull:
		return models.TimeRange{
			Start: fileCursor.FullEpoch.Time,
			End:   fileCursor.FullCursor.Time,
		}, nil
	case models.DahuaScanTypeQuick:
		return models.TimeRange{
			Start: fileCursor.QuickCursor.Time,
			End:   time.Now(),
		}, nil
	case models.DahuaScanTypeReverse:
		startTime, err := db.C().DahuaGetOldestFileStartTime(ctx, fileCursor.DeviceID)
		if err != nil {
			if core.IsNotFound(err) {
				return models.TimeRange{}, nil
			}
			return models.TimeRange{}, err
		}

		start := startTime.Time.Add(-MaxScannerPeriod / 2)
		end := startTime.Time.Add(MaxScannerPeriod / 2)

		return models.TimeRange{
			Start: start,
			End:   end,
		}, nil
	default:
		panic("unknown type")
	}
}

// ScanReset cannot be called concurrently for the same device.
func ScanReset(ctx context.Context, db sqlite.DB, deviceID int64) error {
	fileCursor := newFileCursor()
	_, err := db.C().DahuaUpdateFileCursor(ctx, repo.DahuaUpdateFileCursorParams{
		QuickCursor: fileCursor.QuickCursor,
		FullCursor:  fileCursor.FullCursor,
		FullEpoch:   fileCursor.FullEpoch,
		DeviceID:    deviceID,
		ScanPercent: 0,
		Scan:        false,
		ScanType:    models.DahuaScanTypeUnknown,
	})
	return err
}

// Scan cannot be called concurrently for the same device.
func Scan(ctx context.Context, db sqlite.DB, bus *event.Bus, rpcClient dahuarpc.Conn, device Conn, scanType models.DahuaScanType) error {
	fileCursor, err := db.C().DahuaUpdateFileCursorScanPercent(ctx, repo.DahuaUpdateFileCursorScanPercentParams{
		DeviceID:    device.ID,
		ScanPercent: 0,
	})
	if err != nil {
		return err
	}

	scanRange, err := getScanRange(ctx, db, fileCursor, scanType)
	if err != nil {
		return err
	}
	iterator := NewScannerPeriodIterator(scanRange)

	updated_at := types.NewTime(time.Now())
	mediaFilesC := make(chan []mediafilefind.FindNextFileInfo)

	fileCursor, err = db.C().DahuaUpdateFileCursor(ctx, repo.DahuaUpdateFileCursorParams{
		QuickCursor: fileCursor.QuickCursor,
		FullCursor:  fileCursor.FullCursor,
		FullEpoch:   fileCursor.FullEpoch,
		DeviceID:    device.ID,
		ScanPercent: iterator.Percent(),
		Scan:        true,
		ScanType:    scanType,
	})
	if err != nil {
		return err
	}

	for scannerPeriod, ok := iterator.Next(); ok; scannerPeriod, ok = iterator.Next() {
		createdCount, err := func() (int, error) {
			cancel, errC := ScannerScan(ctx, rpcClient, scannerPeriod, device.Location, mediaFilesC)
			defer cancel()

			var createdCount int

			for {
				select {
				case <-ctx.Done():
					return 0, ctx.Err()
				case err := <-errC:
					return createdCount, err
				case mediaFiles := <-mediaFilesC:
					files, err := NewDahuaFiles(mediaFiles, int(device.Seed), device.Location)
					if err != nil {
						return createdCount, err
					}

					for _, f := range files {
						created, err := upsertDahuaFiles(ctx, db, repo.DahuaCreateFileParams{
							DeviceID:    device.ID,
							Channel:     int64(f.Channel),
							StartTime:   types.NewTime(f.StartTime),
							EndTime:     types.NewTime(f.EndTime),
							Length:      int64(f.Length),
							Type:        f.Type,
							FilePath:    f.FilePath,
							Duration:    int64(f.Duration),
							Disk:        int64(f.Disk),
							VideoStream: f.VideoStream,
							Flags:       types.StringSlice{Slice: f.Flags},
							Events:      types.StringSlice{Slice: f.Events},
							Cluster:     int64(f.Cluster),
							Partition:   int64(f.Partition),
							PicIndex:    int64(f.PicIndex),
							Repeat:      int64(f.Repeat),
							WorkDir:     f.WorkDir,
							WorkDirSn:   f.WorkDirSN,
							UpdatedAt:   updated_at,
							Storage:     StorageFromFilePath(f.FilePath),
						})
						if err != nil {
							return createdCount, err
						}
						if created {
							createdCount++
						}
					}
				}
			}
		}()
		if err != nil {
			return err
		}

		err = db.C().DahuaDeleteFile(ctx, repo.DahuaDeleteFileParams{
			UpdatedAt: updated_at,
			DeviceID:  device.ID,
			Start:     types.NewTime(scannerPeriod.Start.UTC()),
			End:       types.NewTime(scannerPeriod.End.UTC()),
		})
		if err != nil {
			return err
		}

		fileCursor = updateScanFileCursor(fileCursor, scannerPeriod, scanType)

		fileCursor, err = db.C().DahuaUpdateFileCursor(ctx, repo.DahuaUpdateFileCursorParams{
			QuickCursor: fileCursor.QuickCursor,
			FullCursor:  fileCursor.FullCursor,
			FullEpoch:   fileCursor.FullEpoch,
			DeviceID:    device.ID,
			ScanPercent: iterator.Percent(),
			Scan:        true,
			ScanType:    scanType,
		})
		if err != nil {
			return err
		}
		bus.DahuaFileCursorUpdated(event.DahuaFileCursorUpdated{
			Cursor: fileCursor,
		})

		if createdCount > 0 {
			bus.DahuaFileCreated(event.DahuaFileCreated{
				DeviceID: device.ID,
				TimeRange: models.TimeRange{
					Start: scannerPeriod.Start,
					End:   scanRange.End,
				},
				Count: int64(createdCount),
			})
		}
	}
	fileCursor, err = db.C().DahuaUpdateFileCursor(ctx, repo.DahuaUpdateFileCursorParams{
		QuickCursor: fileCursor.QuickCursor,
		FullCursor:  fileCursor.FullCursor,
		FullEpoch:   fileCursor.FullEpoch,
		DeviceID:    device.ID,
		ScanPercent: iterator.Percent(),
		Scan:        false,
		ScanType:    scanType,
	})
	if err != nil {
		return err
	}

	return nil
}

func upsertDahuaFiles(ctx context.Context, db sqlite.DB, arg repo.DahuaCreateFileParams) (bool, error) {
	_, err := db.C().DahuaUpdateFile(ctx, repo.DahuaUpdateFileParams{
		DeviceID:    arg.DeviceID,
		Channel:     arg.Channel,
		StartTime:   arg.StartTime,
		EndTime:     arg.EndTime,
		Length:      arg.Length,
		Type:        arg.Type,
		FilePath:    arg.FilePath,
		Duration:    arg.Duration,
		Disk:        arg.Disk,
		VideoStream: arg.VideoStream,
		Flags:       arg.Flags,
		Events:      arg.Events,
		Cluster:     arg.Cluster,
		Partition:   arg.Partition,
		PicIndex:    arg.PicIndex,
		Repeat:      arg.Repeat,
		WorkDir:     arg.WorkDir,
		WorkDirSn:   arg.WorkDirSn,
		UpdatedAt:   arg.UpdatedAt,
		Storage:     arg.Storage,
	})
	if err == nil {
		return false, nil
	}
	if !core.IsNotFound(err) {
		return false, err
	}

	if _, err := db.C().DahuaCreateFile(ctx, arg); err != nil {
		return false, err
	}

	return true, nil
}
