//go:build !dev

package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed views
var viewsFS embed.FS

func ViewsFS() fs.ReadDirFS {
	return viewsFS
}

//go:embed dist
var distFS embed.FS

func AssetFS() http.FileSystem {
	return http.FS(echo.MustSubFS(distFS, "dist"))
}

//go:embed dist/.vite/manifest.json
var manifestJSON []byte

func Head() template.HTML {
	type Manifest struct {
		CSS     []string `json:"css"`
		File    string   `json:"file"`
		IsEntry bool     `json:"isEntry"`
		Src     string   `json:"src"`
	}

	var manifestMap map[string]Manifest
	if err := json.Unmarshal(manifestJSON, &manifestMap); err != nil {
		panic(err)
	}

	var manifest Manifest
	for _, man := range manifestMap {
		if man.IsEntry {
			manifest = man
			break
		}
	}

	var headTags string
	for _, v := range manifest.CSS {
		headTags += fmt.Sprintf(`<link rel="stylesheet" href="/%s" />`, v)
	}
	headTags += fmt.Sprintf(`<script type="module" src="/%s"></script>`, manifest.File)

	return template.HTML(headTags)
}
