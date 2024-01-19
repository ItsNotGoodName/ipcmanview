//go:build !dev

package view

import (
	"fmt"
	"html/template"

	"github.com/labstack/echo/v4"
)

func (r Renderer) Template(name string) (*template.Template, error) {
	tmpl, found := r.templates[name]
	if !found {
		return nil, echo.ErrInternalServerError.WithInternal(fmt.Errorf("template not found: %s", name))
	}

	return tmpl, nil
}
