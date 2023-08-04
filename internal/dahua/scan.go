package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmango/internal/models"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/modules/mediafilefind"
	"github.com/ItsNotGoodName/ipcmango/pkg/qes"
)

const (
	// MaxScanPeriod is longest allowed time range for a mediafilefind query.
	MaxScanPeriod = 30 * 24 * time.Hour
)

type ScanPeriod struct {
	Start time.Time
	End   time.Time
}

// ScanPeriodIterator generates scan periods that honor MaxScanPeriod.
type ScanPeriodIterator struct {
	start  time.Time
	end    time.Time
	cursor time.Time
}

func NewScanPeriodIterator(scanRange models.DahuaScanRange) *ScanPeriodIterator {
	return &ScanPeriodIterator{
		start:  scanRange.Start,
		end:    scanRange.End,
		cursor: scanRange.End,
	}
}

func (s *ScanPeriodIterator) Next() (ScanPeriod, bool) {
	panic("not implemented")
}

func (s *ScanPeriodIterator) Percent() float64 {
	panic("not implemented")
}

func ScanQuickCursor() time.Time {
	return time.Now().Add(-8 * time.Hour)
}

func ScanQuickCursorFromScanRange(scanRange models.DahuaScanRange) time.Time {
	quickCursor := ScanQuickCursor()
	if scanRange.End.Before(quickCursor) {
		return scanRange.End
	}

	return quickCursor
}

func Scan(ctx context.Context, db qes.Querier, gen dahua.GenRPC, scanCamera models.DahuaScanCamera, scanPeriod ScanPeriod) (ScanResult, error) {
	baseCondition := mediafilefind.
		NewCondtion(
			dahua.NewTimestamp(scanPeriod.Start, scanCamera.Location),
			dahua.NewTimestamp(scanPeriod.End, scanCamera.Location),
		)

	var upserted int64
	updatedAt := time.Now()

	// Pictures
	{
		pictureStream, err := mediafilefind.NewStream(ctx, gen, baseCondition.Picture())
		if err != nil {
			return ScanResult{}, err
		}
		defer pictureStream.Close(gen)

		for {
			files, err := pictureStream.Next(ctx, gen)
			if err != nil {
				return ScanResult{}, err
			}
			if files == nil {
				break
			}

			count, err := DB.ScanCameraFilesUpsert(ctx, db, scanCamera, files, updatedAt)
			if err != nil {
				return ScanResult{}, err
			}

			upserted += count
		}
	}

	// Videos
	{
		videoStream, err := mediafilefind.NewStream(ctx, gen, baseCondition.Video())
		if err != nil {
			return ScanResult{}, err
		}
		defer videoStream.Close(gen)

		for {
			files, err := videoStream.Next(ctx, gen)
			if err != nil {
				return ScanResult{}, err
			}
			if files == nil {
				break
			}

			count, err := DB.ScanCameraFilesUpsert(ctx, db, scanCamera, files, updatedAt)
			if err != nil {
				return ScanResult{}, err
			}

			upserted += count
		}
	}

	deleted, err := DB.ScanCameraFilesDelete(ctx, db, scanCamera.ID, scanPeriod, updatedAt)
	if err != nil {
		return ScanResult{}, err
	}

	return ScanResult{
		Upserted: upserted,
		Deleted:  deleted,
	}, nil
}

type ScanResult struct {
	Upserted int64
	Deleted  int64
}
