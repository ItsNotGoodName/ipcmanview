//go:build dev

package apiws

import "net/http"

func init() {
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
}
