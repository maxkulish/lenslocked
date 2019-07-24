package controller

import (
	"fmt"
	"lenslocked/views"
	"net/http"
)

type Users struct {
	NewView *views.View
}

func NewUser() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
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
	_, _ = fmt.Fprintln(w, "This is a fake message. Pretend that we created the user account")
}
