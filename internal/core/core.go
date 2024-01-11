package core

import (
	"errors"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

func IPFromAddress(address string) string {
	s := strings.Split(address, ":")
	if len(s) != 2 {
		return ""
	}
	return s[0]
}

func IPFromURL(urL *url.URL) (string, error) {
	ip := urL.Hostname()

	ips, err := net.LookupIP(ip)
	if err != nil {
		return "", err
	}

	for _, i2 := range ips {
		if i2.To4() != nil {
			ip = i2.String()
			break
		}
	}

	return ip, nil
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

func StorageFromFilePath(filePath string) models.Storage {
	if strings.HasPrefix(filePath, "sftp://") {
		return models.StorageSFTP
	}
	if strings.HasPrefix(filePath, "ftp://") {
		return models.StorageFTP
	}
	// if strings.HasPrefix(filePath, "nfs://") {
	// 	return models.StorageNFS
	// }
	// if strings.HasPrefix(filePath, "smb://") {
	// 	return models.StorageSMB
	// }
	return models.StorageLocal
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
