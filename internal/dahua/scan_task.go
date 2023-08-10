package dahua

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmango/internal/models"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
	"github.com/ItsNotGoodName/ipcmango/pkg/qes"
)

func NewScanTaskFull(cursor models.DahuaScanCursor) (models.DahuaScanQueueTask, error) {
	if cursor.FullComplete {
		return models.DahuaScanQueueTask{}, fmt.Errorf("full scan complete")
	}

	return models.DahuaScanQueueTask{
		CameraID: cursor.CameraID,
		Range: models.DahuaScanRange{
			Start: cursor.FullEpoch,
			End:   cursor.FullCursor,
		},
		Kind: models.DahuaScanKindFull,
	}, nil
}

func NewScanTaskQuick(cursor models.DahuaScanCursor) models.DahuaScanQueueTask {
	// Prevent start from being older than when the last full scan ends
	start := cursor.QuickCursor
	if start.Before(cursor.FullEpochEnd) { // && cursor.FullComplete {
		start = ScanQuickCursorFromCursor(cursor.FullEpochEnd)
	}

	return models.DahuaScanQueueTask{
		CameraID: cursor.CameraID,
		Range: models.DahuaScanRange{
			Start: start,
			End:   time.Now(),
		},
		Kind: models.DahuaScanKindQuick,
	}
}

func NewScanTaskManual(cursor models.DahuaScanCursor, scanRange models.DahuaScanRange) models.DahuaScanQueueTask {
	return models.DahuaScanQueueTask{
		CameraID: cursor.CameraID,
		Range:    scanRange,
		Kind:     models.DahuaScanKindManual,
	}
}

// ScanTaskQueueExecute runs the queued scan task.
func ScanTaskQueueExecute(ctx context.Context, db qes.Querier, gen dahua.GenRPC, queueTask models.DahuaScanQueueTask) error {
	activeTask, err := DB.ScanActiveTaskCreate(ctx, db, queueTask)
	if err != nil {
		return err
	}

	scanErrString, err := func() (string, error) {
		// WARNING: this assumes the scan cursor will not be modified by other functions
		scanCursor, err := DB.ScanCursorGet(ctx, db, activeTask.CameraID)
		if err != nil {
			return err.Error(), err
		}

		// Run the scan
		scanErr := scanTaskActiveExecute(ctx, db, gen, scanCursor, activeTask)
		if scanErr != nil {
			// Sad path, scan encounterd some sort of error
			if activeTask.Kind == models.DahuaScanKindFull {
				err := DB.ScanCursorUpdateFullCursorFromActiveScanTaskCursor(ctx, db, activeTask.CameraID)
				if err != nil {
					return err.Error() + scanErr.Error(), err
				}
			}

			return scanErr.Error(), nil
		}

		// Happy path, scan was successful
		switch activeTask.Kind {
		case models.DahuaScanKindFull:
			err := DB.ScanCursorUpdateFullCursor(ctx, db, activeTask.CameraID, activeTask.Range.Start)
			if err != nil {
				return err.Error(), err
			}
		case models.DahuaScanKindQuick:
			err := DB.ScanCursorUpdateQuickCursor(ctx, db, activeTask.CameraID, ScanQuickCursorFromScanRange(activeTask.Range))
			if err != nil {
				return err.Error(), err
			}
		}

		return "", nil
	}()

	return errors.Join(err, DB.ScanActiveTaskComplete(ctx, db, activeTask.CameraID, scanErrString))
}

func scanTaskActiveExecute(ctx context.Context, db qes.Querier, gen dahua.GenRPC, scanCursor models.DahuaScanCursor, activeTask models.DahuaScanActiveTask) error {
	scanPeriodIterator := NewScanPeriodIterator(activeTask.Range)
	progress := activeTask.NewProgress()

	for {
		scanPeriod, ok := scanPeriodIterator.Next()
		if !ok {
			break
		}

		res, err := Scan(ctx, db, gen, scanCursor, scanPeriod)
		if err != nil {
			return err
		}

		progress.Upserted += int(res.Upserted)
		progress.Deleted += int(res.Deleted)
		progress.Percent = scanPeriodIterator.Percent()
		progress.Cursor = scanPeriodIterator.Cursor()

		progress, err = DB.ScanProgressUpdate(ctx, db, progress)
		if err != nil {
			return err
		}
	}

	return nil
}
