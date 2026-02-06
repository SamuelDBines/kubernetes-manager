package web

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	t *template.Template
}

func NewRenderer(templatesDir string) (*Renderer, error) {
	// Parse all layouts + pages + components
	patterns := []string{
		filepath.Join(templatesDir, "layouts", "*.html"),
		filepath.Join(templatesDir, "pages", "*.html"),
		filepath.Join(templatesDir, "components", "*.html"),
	}

	t, err := template.ParseFiles(globMust(patterns)...)
	if err != nil {
		return nil, err
	}
	return &Renderer{t: t}, nil
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = r.t.ExecuteTemplate(w, name, data)
}

func globMust(patterns ...[]string) []string {
	var out []string
	for _, ps := range patterns {
		for _, p := range ps {
			matches, _ := filepath.Glob(p)
			out = append(out, matches...)
		}
	}
	return out
}
