package dahua

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmango/internal/models"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
	"github.com/ItsNotGoodName/ipcmango/pkg/qes"
)

func NewScanTaskFull(cam models.DahuaScanCursor) (models.DahuaScanQueueTask, error) {
	if cam.FullComplete {
		return models.DahuaScanQueueTask{}, fmt.Errorf("full scan complete")
	}

	return models.DahuaScanQueueTask{
		CameraID: cam.CameraID,
		Range: models.DahuaScanRange{
			Start: cam.FullEpoch,
			End:   cam.FullCursor,
		},
		Kind: models.DahuaScanKindFull,
	}, nil
}

func NewScanTaskQuick(cam models.DahuaScanCursor) models.DahuaScanQueueTask {
	return models.DahuaScanQueueTask{
		CameraID: cam.CameraID,
		Range: models.DahuaScanRange{
			Start: cam.QuickCursor,
			End:   time.Now(),
		},
		Kind: models.DahuaScanKindQuick,
	}
}

func NewScanTaskManual(cam models.DahuaScanCursor, scanRange models.DahuaScanRange) models.DahuaScanQueueTask {
	return models.DahuaScanQueueTask{
		CameraID: cam.CameraID,
		Range:    scanRange,
		Kind:     models.DahuaScanKindManual,
	}
}

// ScanTaskQueueExecute runs the queued scan task. This should not being called concurrently with the same camera.
func ScanTaskQueueExecute(ctx context.Context, db qes.Querier, gen dahua.GenRPC, queueTask models.DahuaScanQueueTask) error {
	scanCursor, err := DB.ScanCursorGet(ctx, db, queueTask.CameraID)
	if err != nil {
		return err
	}

	activeTask, err := DB.ScanActiveTaskCreate(ctx, db, queueTask)
	if err != nil {
		return err
	}

	var errString string
	if err = scanTaskActiveExecute(ctx, db, gen, scanCursor, activeTask); err != nil {
		errString = err.Error()
		// Update full cursor from active tasks' cursor
		if queueTask.Kind == models.DahuaScanKindFull {
			err := DB.ScanCursorUpdateFullCursorFromActiveScanTaskCursor(ctx, db, queueTask.CameraID)
			if err != nil {
				return err
			}
		}
	} else {
		// Update scan cursors on success
		switch queueTask.Kind {
		case models.DahuaScanKindFull:
			err := DB.ScanCursorUpdateFullCursor(ctx, db, queueTask.CameraID, queueTask.Range.Start)
			if err != nil {
				return err
			}
		case models.DahuaScanKindQuick:
			err := DB.ScanCursorUpdateQuickCursor(ctx, db, queueTask.CameraID, ScanQuickCursorFromScanRange(queueTask.Range))
			if err != nil {
				return err
			}
		}
	}

	return DB.ScanActiveTaskComplete(ctx, db, activeTask.CameraID, errString)
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
