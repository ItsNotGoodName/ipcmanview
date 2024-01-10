//go:build dev

package web

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get package directory")
	}

	cwd = filepath.Dir(filename)
}

var cwd string

func ViewsFS() fs.ReadDirFS {
	return os.DirFS(cwd).(fs.ReadDirFS)
}

var assetFSPath string

func AssetFS() http.FileSystem {
	return http.FS(os.DirFS(filepath.Join(cwd, "public")))
}

func Head() template.HTML {
	host := os.Getenv("VITE_HOST")
	return template.HTML(fmt.Sprintf(`<script type="module" src="http://%s:5173/@vite/client"></script><script type="module" src="http://%s:5173/src/main.ts"></script>`, host, host))
}
