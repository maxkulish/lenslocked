package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"lenslocked/contextd"
	"lenslocked/models"
	"lenslocked/views"
	"log"
	"net/http"
	"strconv"
)

const (
	ShowGallery = "show_gallery"
	EditGallery = "edit_gallery"
)

type Galleries struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	route     *mux.Router
}

func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		New:       views.NewView("bootstrap", "galleries/new"),
		ShowView:  views.NewView("bootstrap", "galleries/show"),
		EditView:  views.NewView("bootstrap", "galleries/edit"),
		IndexView: views.NewView("bootstrap", "galleries/index"),
		gs:        gs,
		route:     r,
	}
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// GET /galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {

	// Check if user is owner of the gallery
	user := contextd.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var vd views.Data
	vd.Body = galleries
	g.IndexView.Render(w, r, vd)
}

// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var vd views.Data
	vd.Body = gallery
	g.ShowView.Render(w, r, vd)
}

// GET /galleries/:id/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	// Check if user is owner of the gallery
	user := contextd.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Body = gallery
	g.EditView.Render(w, r, vd)
}

// GET /galleries/:id/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	// Check if user is owner of the gallery
	user := contextd.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	var form GalleryForm
	var vd views.Data
	vd.Body = gallery
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Gallery successfully updated!",
	}

	g.EditView.Render(w, r, vd)
}

// GET /galleries/:id/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	// Check if user is owner of the gallery
	user := contextd.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Body = gallery
		g.EditView.Render(w, r, vd)
		return
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}

	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went wrong", http.StatusInternalServerError)
		}
		return nil, err

	}

	return gallery, nil
}

// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {

	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}

	user := contextd.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}

	url, err := g.route.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		// Todo: make this go to the index page
		http.Redirect(w, r, "/", http.StatusFound)
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}
