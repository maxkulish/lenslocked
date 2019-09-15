package views

import (
	"bytes"
	"errors"
	"github.com/gorilla/csrf"
	"html/template"
	"io"
	"lenslocked/contextd"
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
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Body: data,
		}
	}

	vd.User = contextd.User(r.Context())
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	_, _ = io.Copy(w, &buf)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

func NewView(layout string, files ...string) *View {

	addTemplatePath(files)
	addTemplateExt(files)

	files = append(files, layoutFiles(LayoutDir+"*"+TemplateExt)...)

	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("CSRF is not implemented")
		},
	}).ParseFiles(files...)
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
