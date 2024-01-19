//go:build dev

package view

import (
	"html/template"
)

func (r Renderer) Template(name string) (*template.Template, error) {
	nameHTML := name + ".html"
	tmpl, err := r.parser.Template(nameHTML)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
