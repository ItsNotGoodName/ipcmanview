package dahua

import (
	"net/url"
	"slices"
)

func toHTTPURL(u *url.URL) *url.URL {
	if slices.Contains([]string{"http", "https"}, u.Scheme) {
		return u
	}

	switch u.Port() {
	case "443":
		u.Scheme = "https"
	default:
		u.Scheme = "http"
	}

	return u
}
