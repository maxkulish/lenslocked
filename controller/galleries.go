package controller

import (
	"fmt"
	"lenslocked/models"
	"lenslocked/views"
	"log"
	"net/http"
)

type Galleries struct {
	New *views.View
	gs  models.GalleryService
}

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {

	var vd views.Data

	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	gallery := models.Gallery{
		Title: form.Title,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	_, _ = fmt.Fprintln(w, gallery)
}
