package views

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

const (
	layoutDir   string = "views/layouts/"
	TemplateExt string = "*.gohtml"
)

type View struct {
	Template *template.Template
	Layout   string
}

// Render is userd to render the view with the predefined layout
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		log.Fatal(err)
	}
}

func NewView(layout string, files ...string) *View {

	files = append(files, layoutFiles(layoutDir+TemplateExt)...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// layoutFiles returns a slice of strings representing
// the layout files used in our application
func layoutFiles(path string) []string {
	layouts, err := filepath.Glob(path)
	if err != nil {
		log.Fatal(err)
	}

	return layouts
}
