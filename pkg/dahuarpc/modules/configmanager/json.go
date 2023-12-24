package configmanager

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type MergeOption struct {
	Path  string
	Value any
}

// Merge only sets if it already exists in the js.
func Merge(js string, options []MergeOption) (string, error) {
	var err error
	for _, opt := range options {
		if !gjson.Get(js, opt.Path).Exists() {
			continue
		}
		js, err = sjson.Set(js, opt.Path, opt.Value)
		if err != nil {
			return "", err
		}
	}
	return js, nil
}
