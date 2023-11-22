package webdahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
)

const scanVolatileDuration = 8 * time.Hour

func NewFileCursor(cameraID int64) sqlc.DahuaFileCursor {
	now := time.Now()
	return sqlc.DahuaFileCursor{
		CameraID:     cameraID,
		QuickCursor:  now.Add(-scanVolatileDuration),
		FullCursor:   now,
		FullEpoch:    dahua.ScanEpoch,
		FullEpochEnd: now,
		FullComplete: false,
	}
}

func FileQuickScan(ctx context.Context, db sqlc.DB, rpcClient dahuarpc.Client, camera sqlc.GetDahuaCameraRow) error {
	// _, err := db.CreateDahuaFileScanLock(ctx, sqlc.CreateDahuaFileScanLockParams{
	// 	CameraID:  camera.ID,
	// 	CreatedAt: time.Now(),
	// })
	// if err != nil {
	// 	return err
	// }
	// defer db.DeleteDahuaCamera(context.TODO(), camera.ID)

	cursor, err := db.GetDahuaFileCursor(ctx, camera.ID)
	if err != nil {
		return err
	}

	scanRange, err := dahua.NewDahuaScanRange(cursor.QuickCursor, time.Now())
	if err != nil {
		return err
	}
	iter := dahua.NewScanPeriodIterator(scanRange)

	filesC := make(chan []mediafilefind.FindNextFileInfo)

	for period, ok := iter.Next(); ok; period, ok = iter.Next() {
		errC := dahua.Scan(ctx, rpcClient, period, camera.Location.Location, filesC)

	inner:
		for {
			select {
			case <-ctx.Done():
				return err
			case err := <-errC:
				if err != nil {
					return sendStreamError(c, stream, err)
				}
				break inner
			case files := <-filesC:
				res, err := dahua.NewDahuaFiles(conn.Camera.ID, files, dahua.GetSeed(conn.Camera), conn.Camera.Location.Location)
				if err != nil {
					return sendStreamError(c, stream, err)
				}

				if err := sendStream(c, stream, res); err != nil {
					return sendStreamError(c, stream, err)
				}
			}
		}
	}

	dahua.Scan(ctx, rpcClient)

	return nil
}
