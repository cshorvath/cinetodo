package view

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"strings"

	"github.com/labstack/echo/v4"
)

// Renderer renders HTML templates for Echo using the standard library templates.
type Renderer struct {
	templates *template.Template
}

// NewRenderer loads every *.html template from the provided file system.
func NewRenderer(fsys fs.FS, funcs template.FuncMap) (*Renderer, error) {
	tmpl := template.New("").Funcs(mergeFuncMaps(defaultFuncMap(), funcs))
	files := make([]string, 0, 32)
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".html") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("view: no templates found")
	}
	if _, err := tmpl.ParseFS(fsys, files...); err != nil {
		return nil, err
	}
	return &Renderer{templates: tmpl}, nil
}

// Render executes the named template with the provided data.
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var ctxData map[string]interface{}
	switch v := data.(type) {
	case nil:
		ctxData = map[string]interface{}{}
	case map[string]interface{}:
		ctxData = v
	default:
		// Allow structs/slices to be passed directly without wrapping.
		return r.templates.ExecuteTemplate(w, name, data)
	}
	ctxData["ctx"] = c
	return r.templates.ExecuteTemplate(w, name, ctxData)
}

func mergeFuncMaps(base, override template.FuncMap) template.FuncMap {
	result := template.FuncMap{}
	for k, v := range base {
		result[k] = v
	}
	for k, v := range override {
		result[k] = v
	}
	return result
}

func defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"dict": func(values ...interface{}) map[string]interface{} {
			if len(values)%2 != 0 {
				panic("dict requires an even number of arguments")
			}
			m := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					panic("dict keys must be strings")
				}
				m[key] = values[i+1]
			}
			return m
		},
	}
}
