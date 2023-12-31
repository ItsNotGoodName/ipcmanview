package dahua

import (
	"context"
	"errors"
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

// ScanLockCreateTry only tries once to create a lock.
func ScanLockCreateTry(ctx context.Context, db repo.DB, deviceID int64) error {
	err := db.DeleteDahuaFileScanLockByAge(ctx, ScanLockStaleTime())
	if err != nil {
		return err
	}

	_, err = db.CreateDahuaFileScanLock(ctx, repo.CreateDahuaFileScanLockParams{
		DeviceID:  deviceID,
		TouchedAt: types.NewTime(time.Now()),
	})
	if err != nil {
		return err
	}

	return nil
}

// ScanLockCreate keeps trying to create a lock.
func ScanLockCreate(ctx context.Context, db repo.DB, deviceID int64) error {
	err := ScanLockCreateTry(ctx, db, deviceID)
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	t := time.NewTicker(scanLockHeartbeat)
	defer t.Stop()

	for {
		err := ScanLockCreateTry(ctx, db, deviceID)
		if err == nil {
			return nil
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}
	}
}

// ScanLockHeartbeat keeps the lock active until context is canceled or the cancel function is called.
// Lock is automatically deleted.
func ScanLockHeartbeat(ctx context.Context, db repo.DB, deviceID int64) context.CancelFunc {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		t := time.NewTicker(scanLockHeartbeat)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				ScanLockDelete(db, deviceID)
				return
			case <-t.C:
				err := db.TouchDahuaFileScanLock(ctx, repo.TouchDahuaFileScanLockParams{
					TouchedAt: types.NewTime(time.Now()),
					DeviceID:  deviceID,
				})
				if err != nil {
					log.Err(err).Msg("Failed to touch scan lock")
				}
			}
		}
	}()
	return cancel
}

func ScanLockDelete(db repo.DB, deviceID int64) {
	err := db.DeleteDahuaFileScanLock(context.Background(), deviceID)
	if err != nil {
		log.Err(err).Msg("Failed to delete scan lock")
	}
}
