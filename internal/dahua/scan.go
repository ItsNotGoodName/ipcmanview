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
	// MaxScanPeriod is longest allowed time range for a mediafilefind query because some cameras give false data past the MaxScanPeriod.
	MaxScanPeriod = 30 * 24 * time.Hour
)

// // ScanEpoch is the oldest time a camera file can exist.
// var ScanEpoch time.Time
//
// func init() {
// 	var err error
// 	ScanEpoch, err = time.ParseInLocation(time.DateTime, "2009-12-31 00:00:00", time.UTC)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// ScanPeriod is INCLUSIVE Start and EXCLUSIVE End.
type ScanPeriod struct {
	Start time.Time
	End   time.Time
}

// ScanPeriodIterator generates scan periods that are equal to or less than MaxScanPeriod.
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
	if s.start.Equal(s.cursor) {
		return ScanPeriod{}, false
	}

	var (
		start  time.Time
		end    time.Time
		cursor time.Time
	)
	{
		end = s.cursor
		maybe_start := s.cursor.Add(-MaxScanPeriod)
		if maybe_start.Before(s.start) {
			cursor = s.start
			start = s.cursor
		} else {
			cursor = maybe_start
			start = maybe_start
		}
	}

	// Only mutation in this struct
	s.cursor = cursor

	return ScanPeriod{Start: start, End: end}, true
}

func (s *ScanPeriodIterator) Percent() float64 {
	return (s.end.Sub(s.cursor).Hours() / s.end.Sub(s.start).Hours()) * 100
}

func (s *ScanPeriodIterator) Cursor() time.Time {
	return s.cursor
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

func Scan(ctx context.Context, db qes.Querier, gen dahua.GenRPC, scanCamera models.DahuaScanCursor, scanPeriod ScanPeriod) (ScanResult, error) {
	baseCondition := mediafilefind.NewCondtion(
		dahua.NewTimestamp(scanPeriod.Start, scanCamera.Location.Location),
		dahua.NewTimestamp(scanPeriod.End, scanCamera.Location.Location),
	)

	var upserted int64
	scannedAt := time.Now()

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

			count, err := DB.ScanCameraFilesUpsert(ctx, db, scannedAt, scanCamera, files)
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

			count, err := DB.ScanCameraFilesUpsert(ctx, db, scannedAt, scanCamera, files)
			if err != nil {
				return ScanResult{}, err
			}

			upserted += count
		}
	}

	deleted, err := DB.ScanCameraFilesDelete(ctx, db, scannedAt, scanCamera.CameraID, scanPeriod)
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
