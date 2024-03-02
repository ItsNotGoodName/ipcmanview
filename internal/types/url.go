package types

import (
	"database/sql/driver"
	"fmt"
	"net/url"
)

func NewURL(u *url.URL) URL {
	return URL{
		URL: u,
	}
}

// URL cannot be nil.
type URL struct {
	*url.URL
}

func (dst *URL) Scan(src any) error {
	switch src := src.(type) {
	case string:
		u, err := url.Parse(src)
		if err != nil {
			return err
		}
		dst.URL = u
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (u URL) Value() (driver.Value, error) {
	return u.String(), nil
}
