package controller

import (
	"fmt"
	"lenslocked/views"
	"net/http"
)

type Users struct {
	NewView *views.View
}

type SignupForm struct {
	Email string `schema:"email"`
	Pass  string `schema:"password"`
}

func NewUser() *Users {

	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
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

	_, _ = fmt.Fprintln(w, form)
}
