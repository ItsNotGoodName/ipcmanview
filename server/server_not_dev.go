//go:build !dev

package server

import (
	"github.com/go-chi/chi/v5"
)

func CORS(r chi.Router) {}
