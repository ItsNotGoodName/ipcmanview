package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmango/internal/models"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
	"github.com/ItsNotGoodName/ipcmango/pkg/qes"
)

func newScanFull(cam models.DahuaScanCamera) models.DahuaScanTask {
	return models.DahuaScanTask{
		CameraID: cam.ID,
		ScanRange: models.DahuaScanRange{
			Start: cam.FullEpoch,
			End:   cam.FullCursor,
		},
		Type: models.DahuaScanTypeFull,
	}
}

func newScanQuick(cam models.DahuaScanCamera) models.DahuaScanTask {
	return models.DahuaScanTask{
		CameraID: cam.ID,
		ScanRange: models.DahuaScanRange{
			Start: cam.QuickCursor,
			End:   time.Now(),
		},
		Type: models.DahuaScanTypeQuick,
	}
}

func newScanManual(cam models.DahuaScanCamera, scanRange models.DahuaScanRange) models.DahuaScanTask {
	return models.DahuaScanTask{
		CameraID:  cam.ID,
		ScanRange: scanRange,
		Type:      models.DahuaScanTypeManual,
	}
}

func scanTaskStart(ctx context.Context, db qes.Querier, gen dahua.GenRPC, scanTask models.DahuaScanTask, scanLock ScanLock) error {
	// TODO: finish this
	scanCamera := models.DahuaScanCamera{}

	iter := NewScanPeriodIterator(scanTask.ScanRange)
	for {
		scanPeriod, ok := iter.Next()
		if !ok {
			break
		}

		_, err := Scan(ctx, db, gen, scanCamera, scanPeriod)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
