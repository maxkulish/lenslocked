package views

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

const (
	LayoutDir   string = "views/layouts/"
	TemplateExt string = ".gohtml"
	TemplateDir string = "views/"
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

	addTemplatePath(files)
	addTemplateExt(files)

	files = append(files, layoutFiles(LayoutDir+"*"+TemplateExt)...)

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

// addTemplatePath takes in a slice of strings
// representing file paths for templates, and it prepends
// the TemplateDir directory to each string in the slice
//
// Eg the in {"home"} wold result in the output
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
