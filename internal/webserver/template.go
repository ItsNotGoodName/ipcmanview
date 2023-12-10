package webserver

import (
	"fmt"
	"html/template"
	"io"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/build"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/internal/web"
	"github.com/Masterminds/sprig/v3"
	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
)

func RegisterRenderer(e *echo.Echo) error {
	r, err := NewRenderer()
	if err != nil {
		return err
	}
	e.Renderer = r
	return nil
}

func parseTemplate(name string) (*template.Template, error) {
	return template.
		New(name).
		Funcs(sprig.FuncMap()).
		Funcs(template.FuncMap{
			"Build": func() build.Build {
				return build.Current
			},
			"Title": func(sub string) string {
				if sub != "" {
					return "IPCManView - " + sub
				}
				return "IPCManView"
			},
			"TimeHumanize": func(date any) string {
				var t time.Time
				switch date := date.(type) {
				case types.Time:
					t = date.Time
				case time.Time:
					t = date
				default:
					panic("invalid date type")
				}
				return humanize.Time(t)
			},
			"BytesHumanize": func(bytes int64) string {
				return humanize.Bytes(uint64(bytes))
			},
			"SLFormatDate": func(date any) template.HTML {
				var t time.Time
				switch date := date.(type) {
				case types.Time:
					t = date.Time
				case time.Time:
					t = date
				default:
					panic("invalid date type")
				}
				return template.HTML(fmt.Sprintf(`<sl-format-date month="numeric" day="numeric" year="numeric" hour="numeric" minute="numeric" hour-format="12" second="numeric" date="%s"></sl-format-date>`, t.Format(time.RFC3339)))
			},
			"Query": func(params any, vals ...any) template.URL {
				length := len(vals)
				query := api.NewQuery(params)
				for i := 0; i < length; i += 2 {
					query.Set(vals[i].(string), fmt.Sprint(vals[i+1]))
				}
				return template.URL(query.Encode())
			},
			"QueryDelete": func(params any, vals ...string) template.URL {
				query := api.NewQuery(params)
				for _, v := range vals {
					query.Del(v)
				}
				return template.URL(query.Encode())
			},
			"FormFormatDate": func(date any) string {
				var t time.Time
				switch date := date.(type) {
				case types.Time:
					t = date.Time
				case time.Time:
					t = date
				default:
					panic("invalid date type")
				}
				return t.Format("2018-06-12T19:30")
			},
		}).
		ParseFS(web.ViewsFS(), "views/partials/*.html", "views/"+name)
}

type TemplateContext struct {
	// Template is the current template that is being rendered.
	Template string
	URL      *url.URL
	Head     template.HTML
	Data     any
}

// TemplateBlock only renders the template block.
type TemplateBlock struct {
	Name string
	any
}

func NewRenderer() (Renderer, error) {
	files, err := web.ViewsFS().ReadDir("views")
	if err != nil {
		return Renderer{}, err
	}

	templates := make(map[string]*template.Template)
	for _, f := range files {
		if !f.IsDir() {
			name := f.Name()
			baseName, _ := strings.CutSuffix(name, filepath.Ext(name))
			templates[baseName] = template.Must(parseTemplate(name))
		}
	}

	return Renderer{
		templates: templates,
		head:      web.Head(),
	}, nil
}

type Renderer struct {
	templates map[string]*template.Template
	head      template.HTML
}

func (t Renderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	tmpl, err := t.Template(name)
	if err != nil {
		return err
	}

	tmplData := TemplateContext{
		Template: name,
		URL:      c.Request().URL,
		Head:     t.head,
	}

	switch data := data.(type) {
	case TemplateBlock:
		tmplData.Data = data.any
		return tmpl.ExecuteTemplate(w, data.Name, tmplData)
	default:
		tmplData.Data = data
		return tmpl.Execute(w, tmplData)
	}
}
