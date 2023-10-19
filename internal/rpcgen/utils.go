package rpcgen

import (
	"net/http"
	"strings"
)

func IsRPC(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, "/rpc/")
}
