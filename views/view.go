package views

import (
	"html/template"
	"log"
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
