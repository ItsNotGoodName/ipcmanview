package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
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
	case models.DahuaScanType_Full:
		// Update FullCursor
		if scanPeriod.Start.Before(fileCursor.FullCursor.Time) {
			fileCursor.FullCursor = types.NewTime(scanPeriod.Start)
		}
	case models.DahuaScanType_Quick:
		// Update QuickCursor
		quickCursor := time.Now().Add(-scanVolatileDuration)
		if scanPeriod.End.Before(quickCursor) {
			fileCursor.QuickCursor = types.NewTime(scanPeriod.End)
		} else {
			fileCursor.QuickCursor = types.NewTime(quickCursor)
		}
	case models.DahuaScanType_Reverse:
	default:
		panic("unknown type")
	}

	return fileCursor
}

func getScanRange(ctx context.Context, fileCursor repo.DahuaFileCursor, scanType models.DahuaScanType) (models.TimeRange, error) {
	switch scanType {
	case models.DahuaScanType_Full:
		return models.TimeRange{
			Start: fileCursor.FullEpoch.Time,
			End:   fileCursor.FullCursor.Time,
		}, nil
	case models.DahuaScanType_Quick:
		return models.TimeRange{
			Start: fileCursor.QuickCursor.Time,
			End:   time.Now(),
		}, nil
	case models.DahuaScanType_Reverse:
		startTime, err := app.DB.C().DahuaGetOldestFileStartTime(ctx, fileCursor.DeviceID)
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
func ScanReset(ctx context.Context, deviceID int64) error {
	fileCursor := newFileCursor()
	_, err := app.DB.C().DahuaUpdateFileCursor(ctx, repo.DahuaUpdateFileCursorParams{
		QuickCursor: fileCursor.QuickCursor,
		FullCursor:  fileCursor.FullCursor,
		FullEpoch:   fileCursor.FullEpoch,
		DeviceID:    deviceID,
		ScanPercent: 0,
		Scanning:    false,
		ScanType:    models.DahuaScanType_Unknown,
	})
	return err
}

// Scan cannot be called concurrently for the same device.
func Scan(ctx context.Context, rpcClient dahuarpc.Conn, device Conn, scanType models.DahuaScanType) error {
	fileCursor, err := app.DB.C().DahuaUpdateFileCursorScanPercent(ctx, repo.DahuaUpdateFileCursorScanPercentParams{
		DeviceID:    device.ID,
		ScanPercent: 0,
	})
	if err != nil {
		return err
	}

	scanRange, err := getScanRange(ctx, fileCursor, scanType)
	if err != nil {
		return err
	}
	iterator := NewScannerPeriodIterator(scanRange)

	updatedAt := types.NewTime(time.Now())
	mediaFilesC := make(chan []mediafilefind.FindNextFileInfo)

	fileCursor, err = app.DB.C().DahuaUpdateFileCursor(ctx, repo.DahuaUpdateFileCursorParams{
		QuickCursor: fileCursor.QuickCursor,
		FullCursor:  fileCursor.FullCursor,
		FullEpoch:   fileCursor.FullEpoch,
		DeviceID:    device.ID,
		ScanPercent: iterator.Percent(),
		Scanning:    true,
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
						created, err := upsertFile(ctx, repo.DahuaCreateFileParams{
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
							Storage:     StorageFromFilePath(f.FilePath),
							Source:      models.DahuaFileSource_Device,
							UpdatedAt:   updatedAt,
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

		err = app.DB.C().DahuaDeleteFile(ctx, repo.DahuaDeleteFileParams{
			DeviceID:  device.ID,
			Start:     types.NewTime(scannerPeriod.Start.UTC()),
			End:       types.NewTime(scannerPeriod.End.UTC()),
			UpdatedAt: updatedAt,
			Source:    models.DahuaFileSource_Device,
		})
		if err != nil {
			return err
		}

		fileCursor = updateScanFileCursor(fileCursor, scannerPeriod, scanType)

		fileCursor, err = app.DB.C().DahuaUpdateFileCursor(ctx, repo.DahuaUpdateFileCursorParams{
			QuickCursor: fileCursor.QuickCursor,
			FullCursor:  fileCursor.FullCursor,
			FullEpoch:   fileCursor.FullEpoch,
			DeviceID:    device.ID,
			ScanPercent: iterator.Percent(),
			Scanning:    true,
			ScanType:    scanType,
		})
		if err != nil {
			return err
		}
		app.Hub.DahuaFileCursorUpdated(bus.DahuaFileCursorUpdated{
			Cursor: fileCursor,
		})

		if createdCount > 0 {
			app.Hub.DahuaFileCreated(bus.DahuaFileCreated{
				DeviceID: device.ID,
				TimeRange: models.TimeRange{
					Start: scannerPeriod.Start,
					End:   scanRange.End,
				},
				Count: int64(createdCount),
			})
		}
	}
	fileCursor, err = app.DB.C().DahuaUpdateFileCursor(ctx, repo.DahuaUpdateFileCursorParams{
		QuickCursor: fileCursor.QuickCursor,
		FullCursor:  fileCursor.FullCursor,
		FullEpoch:   fileCursor.FullEpoch,
		DeviceID:    device.ID,
		ScanPercent: iterator.Percent(),
		Scanning:    false,
		ScanType:    scanType,
	})
	if err != nil {
		return err
	}

	return nil
}
