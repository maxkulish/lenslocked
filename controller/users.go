package controller

import (
	"fmt"
	"lenslocked/database"
	"lenslocked/models"
	"lenslocked/views"
	"net/http"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
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

func NewUser(us *models.UserService) *Users {

	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
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
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Pass,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	signIn(w, &user)
	http.Redirect(w, r, "/cookietest", http.StatusFound)

	_, _ = fmt.Fprintln(w, user)
}

// Login is used to verify the provided email adress and
// password and then log the user in if they are correct
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Pass)

	if err != nil {
		switch err {
		case database.ErrNotFound:
			_, _ = fmt.Fprintln(w, "Invalid email address.")
		case database.ErrInvalidPass:
			_, _ = fmt.Fprintln(w, "Invalid password provided")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	signIn(w, user)
	_, _ = fmt.Fprint(w, user)

}

// signIn is used to sign in the given user in via cookies
func signIn(w http.ResponseWriter, user *models.User) {
	cookie := http.Cookie{
		Name:  "email",
		Value: user.Email,
	}

	http.SetCookie(w, &cookie)
}

// CookieTest is used to display cookies set on the current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("email")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintln(w, "Email is:", cookie.Value)
}
