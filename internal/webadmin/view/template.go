package view

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
	"github.com/ItsNotGoodName/ipcmanview/internal/webadmin"
	"github.com/Masterminds/sprig/v3"
	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
)

type Config struct{}

type Data map[string]any

type Block struct {
	Name string
	Data any
}

func newFuncMap(config Config) template.FuncMap {
	return template.FuncMap{
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
			query := api.EncodeQuery(params)
			for i := 0; i < length; i += 2 {
				query.Set(vals[i].(string), fmt.Sprint(vals[i+1]))
			}
			return template.URL("?" + query.Encode())
		},
		// "QueryDelete": func(params any, vals ...string) template.URL {
		// 	query := api.EncodeQuery(params)
		// 	for _, v := range vals {
		// 		query.Del(v)
		// 	}
		// 	return template.URL(query.Encode())
		// },
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
	}
}

type parser struct {
	sprigFuncMap template.FuncMap
	funcMap      template.FuncMap
}

func (p parser) Template(name string) (*template.Template, error) {
	return template.
		New(name).
		Funcs(p.sprigFuncMap).
		Funcs(p.funcMap).
		ParseFS(webadmin.ViewsFS(), "views/partials/*.html", "views/"+name)
}

type TemplateContext struct {
	// Template is the current template that is being rendered.
	Template string
	URL      *url.URL
	Head     template.HTML
	Data     any
}

func WithRenderer(e *echo.Echo, config Config) (*echo.Echo, error) {
	files, err := webadmin.ViewsFS().ReadDir("views")
	if err != nil {
		return nil, err
	}

	parser := parser{
		sprigFuncMap: sprig.FuncMap(),
		funcMap:      newFuncMap(config),
	}

	templates := make(map[string]*template.Template)
	for _, f := range files {
		if !f.IsDir() {
			name := f.Name()
			baseName, _ := strings.CutSuffix(name, filepath.Ext(name))
			templates[baseName] = template.Must(parser.Template(name))
		}
	}

	e.Renderer = Renderer{
		head:      webadmin.Head(),
		templates: templates,
		parser:    parser,
	}

	return e, nil
}

type Renderer struct {
	head      template.HTML
	templates map[string]*template.Template
	parser    parser
}

func (r Renderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	tmpl, err := r.Template(name)
	if err != nil {
		return err
	}

	tmplData := TemplateContext{
		Template: name,
		URL:      c.Request().URL,
		Head:     r.head,
	}

	switch data := data.(type) {
	case Block:
		tmplData.Data = data.Data
		return tmpl.ExecuteTemplate(w, data.Name, tmplData)
	default:
		tmplData.Data = data
		return tmpl.Execute(w, tmplData)
	}
}
