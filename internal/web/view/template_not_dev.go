//go:build !dev

package view

import (
	"fmt"
	"html/template"

	"github.com/labstack/echo/v4"
)

func (t Renderer) Template(name string) (*template.Template, error) {
	tmpl, found := t.templates[name]
	if !found {
		return nil, echo.ErrInternalServerError.WithInternal(fmt.Errorf("template not found: %s", name))
	}

	return tmpl, nil
}
