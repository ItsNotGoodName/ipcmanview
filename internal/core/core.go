package core

import (
	"database/sql"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

func SplitAddress(address string) [2]string {
	s := strings.SplitN(address, ":", 2)
	if len(s) != 2 {
		return [2]string{address}
	}
	return [2]string{s[0], s[1]}
}

func Address(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}

func NewTimeRange(start, end time.Time) (models.TimeRange, error) {
	if end.Before(start) {
		return models.TimeRange{}, errors.New("invalid time range: end is before start")
	}

	return models.TimeRange{
		Start: start,
		End:   end,
	}, nil
}

type MultiReadCloser struct {
	io.Reader
	Closers []func() error
}

func (c MultiReadCloser) Close() error {
	var multiErr error
	for _, closer := range c.Closers {
		err := closer()
		if err != nil {
			multiErr = errors.Join(multiErr, err)
		}
	}
	return multiErr
}

func Int64ToNullInt64(a int64) sql.NullInt64 {
	if a == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: a,
		Valid: true,
	}
}
