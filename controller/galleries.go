package controller

import (
	"lenslocked/models"
	"lenslocked/views"
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
