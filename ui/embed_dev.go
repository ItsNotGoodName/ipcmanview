//go:build dev

package ui

import (
	"io/fs"
)

var FS = empty{}

type empty struct{}

func (empty) Open(string) (fs.File, error) {
	return nil, fs.ErrNotExist
}
