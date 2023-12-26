// webnext contains the SPA.
package webnext

import (
	"embed"
	"net/http"
)

//go:generate pnpm install
//go:generate pnpm run build

//go:embed dist
var dist embed.FS

func DistFS() http.FileSystem {
	return http.FS(dist)
}
