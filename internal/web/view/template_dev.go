//go:build dev

package view

import (
	"html/template"
)

func (t Renderer) Template(name string) (*template.Template, error) {
	nameHTML := name + ".html"
	tmpl, err := parseTemplate(nameHTML, t.config)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
