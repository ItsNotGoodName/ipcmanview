package configmanager

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type MergeValues struct {
	Path  string
	Value any
}

// Merge merges values into the json.
// Values are only merged if the path already exists in the original json.
func Merge(json string, values []MergeValues) (string, error) {
	var err error
	for _, opt := range values {
		if !gjson.Get(json, opt.Path).Exists() {
			continue
		}
		json, err = sjson.Set(json, opt.Path, opt.Value)
		if err != nil {
			return "", err
		}
	}
	return json, nil
}
