package dahuacore

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
)

func init() {
	var err error
	ScanEpoch, err = time.ParseInLocation(time.DateTime, "2009-12-31 00:00:00", time.UTC)
	if err != nil {
		panic(err)
	}
}

const (
	// MaxScanPeriod is longest allowed time range for a mediafilefind query because some cameras give invalid data past the MaxScanPeriod.
	MaxScanPeriod = 30 * 24 * time.Hour
)

// ScanEpoch is the oldest time a camera file can exist.
var ScanEpoch time.Time

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

func NewScanPeriodIterator(scanRange models.TimeRange) *ScanPeriodIterator {
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
			start = s.start
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

func Scan(
	ctx context.Context,
	rpcClient dahuarpc.Conn,
	scanPeriod ScanPeriod,
	location *time.Location,
	resC chan<- []mediafilefind.FindNextFileInfo,
) <-chan error {
	errC := make(chan error, 1)

	go func() {
		err := scan(ctx, rpcClient, scanPeriod, location, resC)
		errC <- err
	}()

	return errC
}

func scan(
	ctx context.Context,
	rpcClient dahuarpc.Conn,
	scanPeriod ScanPeriod,
	location *time.Location,
	resC chan<- []mediafilefind.FindNextFileInfo,
) error {
	baseCondition := mediafilefind.NewCondtion(
		dahuarpc.NewTimestamp(scanPeriod.Start, location),
		dahuarpc.NewTimestamp(scanPeriod.End, location),
	)

	// Pictures
	{
		pictureStream, err := mediafilefind.NewStream(ctx, rpcClient, baseCondition.Picture())
		if err != nil {
			return err
		}
		defer pictureStream.Close(rpcClient)

		for {
			files, err := pictureStream.Next(ctx, rpcClient)
			if err != nil {
				return err
			}
			if files == nil {
				break
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case resC <- files:
			}
		}
	}

	// Videos
	{
		videoStream, err := mediafilefind.NewStream(ctx, rpcClient, baseCondition.Video())
		if err != nil {
			return err
		}
		defer videoStream.Close(rpcClient)

		for {
			files, err := videoStream.Next(ctx, rpcClient)
			if err != nil {
				return err
			}
			if files == nil {
				break
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case resC <- files:
			}
		}
	}

	return nil
}
