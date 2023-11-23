package webdahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
)

type ScanType string

var (
	ScanTypeFull  ScanType = "full"
	ScanTypeQuick ScanType = "quick"
)

const scanVolatileDuration = 8 * time.Hour

func DefaultFileCursor() sqlc.CreateDahuaFileCursorParams {
	now := time.Now()
	return sqlc.CreateDahuaFileCursorParams{
		QuickCursor: now.Add(-scanVolatileDuration),
		FullCursor:  now,
		FullEpoch:   dahua.ScanEpoch,
	}
}

func updateFileCursor(fileCursor sqlc.DahuaFileCursor, scanPeriod dahua.ScanPeriod, scanType ScanType) sqlc.DahuaFileCursor {
	switch scanType {
	case ScanTypeFull:
		// Update FullCursor
		if scanPeriod.Start.Before(fileCursor.FullCursor) {
			fileCursor.FullCursor = scanPeriod.Start
		}
	case ScanTypeQuick:
		// Update QuickCursor
		quickCursor := time.Now().Add(-scanVolatileDuration)
		if scanPeriod.End.Before(quickCursor) {
			fileCursor.QuickCursor = scanPeriod.End
		} else {
			fileCursor.QuickCursor = quickCursor
		}
	default:
		panic("unknown type")
	}

	return fileCursor
}

func getScanRange(fileCursor sqlc.DahuaFileCursor, scanType ScanType) models.DahuaScanRange {
	switch scanType {
	case ScanTypeFull:
		return models.DahuaScanRange{
			Start: fileCursor.FullEpoch,
			End:   fileCursor.FullCursor,
		}
	case ScanTypeQuick:
		return models.DahuaScanRange{
			Start: fileCursor.QuickCursor,
			End:   time.Now(),
		}
	default:
		panic("unknown type")
	}
}

// Scan should only be called once per camera.
func Scan(ctx context.Context, db sqlc.DB, rpcClient dahuarpc.Client, camera models.DahuaCamera, scanType ScanType) error {
	fileCursor, err := db.GetDahuaFileCursor(ctx, camera.ID)
	if err != nil {
		return err
	}

	updated_at := time.Now()
	iterator := dahua.NewScanPeriodIterator(getScanRange(fileCursor, scanType))
	mediaFilesC := make(chan []mediafilefind.FindNextFileInfo)

	for scanPeriod, ok := iterator.Next(); ok; scanPeriod, ok = iterator.Next() {
		errC := dahua.Scan(ctx, rpcClient, scanPeriod, camera.Location.Location, mediaFilesC)

	inner:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-errC:
				if err != nil {
					return err
				}
				break inner
			case mediaFiles := <-mediaFilesC:
				files, err := dahua.NewDahuaFiles(camera.ID, mediaFiles, int(camera.Seed), camera.Location.Location)
				if err != nil {
					return err
				}

				for _, f := range files {
					_, err := db.UpsertDahuaFiles(ctx, sqlc.CreateDahuaFileParams{
						CameraID:    camera.ID,
						Channel:     int64(f.Channel),
						StartTime:   f.StartTime.UTC(),
						EndTime:     f.EndTime.UTC(),
						Length:      int64(f.Length),
						Type:        f.Type,
						FilePath:    f.FilePath,
						Duration:    int64(f.Duration),
						Disk:        int64(f.Disk),
						VideoStream: f.VideoStream,
						Flags:       models.StringSlice{Slice: f.Flags},
						Events:      models.StringSlice{Slice: f.Events},
						Cluster:     int64(f.Cluster),
						Partition:   int64(f.Partition),
						PicIndex:    int64(f.PicIndex),
						Repeat:      int64(f.Repeat),
						WorkDir:     f.WorkDir,
						WorkDirSn:   int64(f.WorkDirSN),
						UpdatedAt:   updated_at.UTC(),
					})
					if err != nil {
						return err
					}
				}
			}
		}

		err := db.DeleteDahuaFile(ctx, sqlc.DeleteDahuaFileParams{
			UpdatedAt: updated_at.UTC(),
			CameraID:  camera.ID,
			Start:     scanPeriod.Start.UTC(),
			End:       scanPeriod.End.UTC(),
		})
		if err != nil {
			return err
		}

		fileCursor = updateFileCursor(fileCursor, scanPeriod, scanType)
		fileCursor, err = db.UpdateDahuaFileCursor(ctx, sqlc.UpdateDahuaFileCursorParams{
			QuickCursor: fileCursor.QuickCursor,
			FullCursor:  fileCursor.FullCursor,
			FullEpoch:   fileCursor.FullEpoch,
			CameraID:    camera.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
