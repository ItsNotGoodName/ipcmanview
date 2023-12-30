package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
)

func init() {
	var err error
	ScannerEpoch, err = time.ParseInLocation(time.DateTime, "2009-12-31 00:00:00", time.UTC)
	if err != nil {
		panic(err)
	}
}

const (
	// MaxScannerPeriod is longest allowed time range for a mediafilefind query because some devices give invalid data past the MaxScannerPeriod.
	MaxScannerPeriod = 30 * 24 * time.Hour
)

// ScannerEpoch is the oldest time a file can exist.
var ScannerEpoch time.Time

// ScannerPeriod is INCLUSIVE Start and EXCLUSIVE End.
type ScannerPeriod struct {
	Start time.Time
	End   time.Time
}

// ScannerPeriodIterator generates scan periods that are equal to or less than MaxScanPeriod.
type ScannerPeriodIterator struct {
	start  time.Time
	end    time.Time
	cursor time.Time
}

func NewScannerPeriodIterator(scanRange models.TimeRange) *ScannerPeriodIterator {
	return &ScannerPeriodIterator{
		start:  scanRange.Start,
		end:    scanRange.End,
		cursor: scanRange.End,
	}
}

func (s *ScannerPeriodIterator) Next() (ScannerPeriod, bool) {
	if s.start.Equal(s.cursor) {
		return ScannerPeriod{}, false
	}

	var (
		start  time.Time
		end    time.Time
		cursor time.Time
	)
	{
		end = s.cursor
		maybe_start := s.cursor.Add(-MaxScannerPeriod)
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

	return ScannerPeriod{Start: start, End: end}, true
}

func (s *ScannerPeriodIterator) Percent() float64 {
	if s.cursor.Equal(s.start) {
		return 100.0
	}
	return (s.end.Sub(s.cursor).Hours() / s.end.Sub(s.start).Hours()) * 100
}

func (s *ScannerPeriodIterator) Cursor() time.Time {
	return s.cursor
}

func Scanner(
	ctx context.Context,
	rpcClient dahuarpc.Conn,
	scanPeriod ScannerPeriod,
	location *time.Location,
	resC chan<- []mediafilefind.FindNextFileInfo,
) (context.CancelFunc, <-chan error) {
	ctx, cancel := context.WithCancel(ctx)
	errC := make(chan error, 1)

	go func() {
		errC <- scanner(ctx, rpcClient, scanPeriod, location, resC)
		cancel()
	}()

	return cancel, errC
}

func scanner(
	ctx context.Context,
	rpcClient dahuarpc.Conn,
	scanPeriod ScannerPeriod,
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
