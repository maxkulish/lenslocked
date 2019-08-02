package controller

import (
	"fmt"
	"lenslocked/models"
	"lenslocked/views"
	"net/http"
)

type Users struct {
	NewView *views.View
	us      *models.UserService
}

type SignupForm struct {
	Name  string `schema:"name"`
	Email string `schema:"email"`
	Pass  string `schema:"password"`
}

func NewUser(us *models.UserService) *Users {

	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		us:      us,
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	_ = u.NewView.Render(w, nil)
}

// This is used to process sign up form when a user tries
// to create a new account
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user := models.User{
		Name:  form.Name,
		Email: form.Email,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintln(w, form)
}
