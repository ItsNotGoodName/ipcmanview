package api

import (
	"io"
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store dahua.Store
}

func NewHandler(store dahua.Store) Handler {
	return Handler{
		store: store,
	}
}

func (h Handler) WithID(next func(w http.ResponseWriter, r *http.Request, h Handler, id int64)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next(w, r, h, id)
	}
}

func Snapshot(w http.ResponseWriter, r *http.Request, h Handler, id int64) {
	ctx := r.Context()

	cgi, err := h.store.ClientCGI(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	snapshot, err := dahuacgi.SnapshotGet(ctx, cgi, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer snapshot.Close()

	header := w.Header()
	header.Add("Content-Type", snapshot.ContentType)
	header.Add("Content-Length", snapshot.ContentLength)
	io.Copy(w, snapshot)
}
