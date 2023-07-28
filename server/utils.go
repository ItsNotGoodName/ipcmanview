package server

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MountFS adds GET handlers for all files and directories for the given filesystem.
func MountFS(r chi.Router, f fs.FS) {
	httpFS := http.FS(f)
	fsHandler := http.StripPrefix("/", http.FileServer(httpFS))

	if files, err := fs.ReadDir(f, "."); err == nil {
		for _, f := range files {
			name := f.Name()
			if f.IsDir() {
				r.Get("/"+name+"/*", fsHandler.ServeHTTP)
			} else if name == "index.html" {
				indexHandler := indexGet(httpFS)
				r.Get("/", indexHandler)
				r.Get("/index.html", indexHandler)
			} else {
				r.Get("/"+name, fsHandler.ServeHTTP)
			}
		}
	} else if err != fs.ErrNotExist {
		panic(err)
	}
}

// indexGet returns index.html from the given filesystem.
func indexGet(httpFS http.FileSystem) http.HandlerFunc {
	index, err := httpFS.Open("/index.html")
	if err != nil {
		panic(err)
	}

	stat, err := index.Stat()
	if err != nil {
		panic(err)
	}

	modtime := stat.ModTime()

	return func(rw http.ResponseWriter, r *http.Request) {
		http.ServeContent(rw, r, "index.html", modtime, index)
	}
}
