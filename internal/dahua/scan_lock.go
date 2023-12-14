package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/rs/zerolog/log"
)

const (
	scanLockHeartbeat = 10 * time.Second
	scanLockStale     = 30 * time.Second
)

func ScanLockStaleTime() types.Time {
	return types.NewTime(time.Now().Add(-scanLockStale))
}

func ScanLockCreate(ctx context.Context, db repo.DB, cameraID int64) error {
	err := db.DeleteDahuaFileScanLockByAge(ctx, ScanLockStaleTime())
	if err != nil {
		return err
	}

	_, err = db.CreateDahuaFileScanLock(ctx, repo.CreateDahuaFileScanLockParams{
		CameraID:  cameraID,
		TouchedAt: types.NewTime(time.Now()),
	})
	if err != nil {
		return err
	}

	return nil
}

// ScanLockHeartbeat keeps the lock active until context is canceled or the cancel function is called.
// Lock is automatically deleted.
func ScanLockHeartbeat(ctx context.Context, db repo.DB, cameraID int64) context.CancelFunc {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		t := time.NewTicker(scanLockHeartbeat)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				ScanLockDelete(db, cameraID)
				return
			case <-t.C:
				err := db.TouchDahuaFileScanLock(ctx, repo.TouchDahuaFileScanLockParams{
					TouchedAt: types.NewTime(time.Now()),
					CameraID:  cameraID,
				})
				if err != nil {
					log.Err(err).Msg("Failed to touch scan lock")
				}
			}
		}
	}()
	return cancel

}

func ScanLockDelete(db repo.DB, cameraID int64) {
	err := db.DeleteDahuaFileScanLock(context.Background(), cameraID)
	if err != nil {
		log.Err(err).Msg("Failed to delete scan lock")
	}
}
