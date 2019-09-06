package controller

import (
	"fmt"
	"lenslocked/models"
	"lenslocked/rand"
	"lenslocked/views"
	"log"
	"net/http"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

type SignupForm struct {
	Name  string `schema:"name"`
	Email string `schema:"email"`
	Pass  string `schema:"password"`
}

type LoginForm struct {
	Email string `schema:"email"`
	Pass  string `schema:"password"`
}

func NewUser(us models.UserService) *Users {

	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

// New is used to render the form where a user can
// create a new user account.
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}

// This is used to process sign up form when a user tries
// to create a new account
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	var vd views.Data

	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Pass,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)

	_, _ = fmt.Fprintln(w, user)
}

// Login is used to verify the provided email adress and
// password and then log the user in if they are correct
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

	var vd views.Data
	var form LoginForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Pass)

	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("Invalid email address")
		case models.ErrPasswordIncorrect:
			vd.AlertError("Invalid password provided")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)

}

// signIn is used to sign in the given user in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	return nil
}
