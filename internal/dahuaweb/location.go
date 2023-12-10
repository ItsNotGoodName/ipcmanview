package dahuaweb

import (
	_ "embed"
	"strings"
)

var Locations []string

//go:embed locations.txt
var locationsStr string

func init() {
	for _, location := range strings.Split(locationsStr, "\n") {
		if location != "" {
			Locations = append(Locations, location)
		}
	}
}
