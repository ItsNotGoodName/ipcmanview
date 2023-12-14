package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
)

type ScanType string

var (
	ScanTypeFull  ScanType = "full"
	ScanTypeQuick ScanType = "quick"
)

const scanVolatileDuration = 8 * time.Hour

func NewFileCursor() repo.CreateDahuaFileCursorParams {
	now := time.Now()
	return repo.CreateDahuaFileCursorParams{
		CameraID:    0,
		QuickCursor: types.NewTime(now.Add(-scanVolatileDuration)),
		FullCursor:  types.NewTime(now),
		FullEpoch:   types.NewTime(dahuacore.ScanEpoch),
		Percent:     0,
	}
}

func updateFileCursor(fileCursor repo.DahuaFileCursor, scanPeriod dahuacore.ScanPeriod, scanType ScanType) repo.DahuaFileCursor {
	switch scanType {
	case ScanTypeFull:
		// Update FullCursor
		if scanPeriod.Start.Before(fileCursor.FullCursor.Time) {
			fileCursor.FullCursor = types.NewTime(scanPeriod.Start)
		}
	case ScanTypeQuick:
		// Update QuickCursor
		quickCursor := time.Now().Add(-scanVolatileDuration)
		if scanPeriod.End.Before(quickCursor) {
			fileCursor.QuickCursor = types.NewTime(scanPeriod.End)
		} else {
			fileCursor.QuickCursor = types.NewTime(quickCursor)
		}
	default:
		panic("unknown type")
	}

	return fileCursor
}

func getScanRange(fileCursor repo.DahuaFileCursor, scanType ScanType) models.TimeRange {
	switch scanType {
	case ScanTypeFull:
		return models.TimeRange{
			Start: fileCursor.FullEpoch.Time,
			End:   fileCursor.FullCursor.Time,
		}
	case ScanTypeQuick:
		return models.TimeRange{
			Start: fileCursor.QuickCursor.Time,
			End:   time.Now(),
		}
	default:
		panic("unknown type")
	}
}

// ScanReset should only be called once per camera.
func ScanReset(ctx context.Context, db repo.DB, id int64) error {
	fileCursor := NewFileCursor()
	_, err := db.UpdateDahuaFileCursor(ctx, repo.UpdateDahuaFileCursorParams{
		QuickCursor: fileCursor.QuickCursor,
		FullCursor:  fileCursor.FullCursor,
		FullEpoch:   fileCursor.FullEpoch,
		CameraID:    id,
		Percent:     0,
	})
	return err
}

// Scan should only be called once per camera.
func Scan(ctx context.Context, db repo.DB, rpcClient dahuarpc.Conn, camera models.DahuaConn, scanType ScanType) error {
	fileCursor, err := db.UpdateDahuaFileCursorPercent(ctx, repo.UpdateDahuaFileCursorPercentParams{
		CameraID: camera.ID,
		Percent:  0,
	})
	if err != nil {
		return err
	}

	updated_at := types.NewTime(time.Now())
	iterator := dahuacore.NewScanPeriodIterator(getScanRange(fileCursor, scanType))
	mediaFilesC := make(chan []mediafilefind.FindNextFileInfo)

	for scanPeriod, ok := iterator.Next(); ok; scanPeriod, ok = iterator.Next() {
		errC := dahuacore.Scan(ctx, rpcClient, scanPeriod, camera.Location, mediaFilesC)

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
				files, err := dahuacore.NewDahuaFiles(camera.ID, mediaFiles, int(camera.Seed), camera.Location)
				if err != nil {
					return err
				}

				for _, f := range files {
					_, err := db.UpsertDahuaFiles(ctx, repo.CreateDahuaFileParams{
						CameraID:    camera.ID,
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
						WorkDirSn:   int64(f.WorkDirSN),
						UpdatedAt:   updated_at,
					})
					if err != nil {
						return err
					}
				}
			}
		}

		err := db.DeleteDahuaFile(ctx, repo.DeleteDahuaFileParams{
			UpdatedAt: updated_at,
			CameraID:  camera.ID,
			Start:     types.NewTime(scanPeriod.Start.UTC()),
			End:       types.NewTime(scanPeriod.End.UTC()),
		})
		if err != nil {
			return err
		}

		fileCursor = updateFileCursor(fileCursor, scanPeriod, scanType)
		fileCursor, err = db.UpdateDahuaFileCursor(ctx, repo.UpdateDahuaFileCursorParams{
			QuickCursor: fileCursor.QuickCursor,
			FullCursor:  fileCursor.FullCursor,
			FullEpoch:   fileCursor.FullEpoch,
			CameraID:    camera.ID,
			Percent:     iterator.Percent(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
