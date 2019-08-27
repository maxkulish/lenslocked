package main

import (
	"fmt"
	"lenslocked/controller"
	"lenslocked/middleware"
	"lenslocked/models"
	"net/http"

	"github.com/gorilla/mux"
)

const env = "dev"

func notFound404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprint(w, "I can't find this page")
}

func main() {
	services, err := models.NewServices(env)
	if err != nil {
		panic(err)
	}

	defer services.Close()
	_ = services.AutoMigrate()

	// Reset all tables
	//_ = services.FullReset()

	staticC := controller.NewStatic()
	usersC := controller.NewUser(services.User)
	galleriesC := controller.NewGalleries(services.Gallery)

	requireUserMW := middleware.RequireUser{
		UserService: services.User,
	}

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	// Gallery routes
	galleryNew := requireUserMW.Apply(galleriesC.New)
	r.Handle("/galleries/new", galleryNew).Methods("GET")
	r.HandleFunc("/galleries", galleriesC.Create).Methods("POST")

	fmt.Println("Starting the server on :3000")

	handler404 := http.HandlerFunc(notFound404)
	r.NotFoundHandler = handler404

	_ = http.ListenAndServe(":3000", r)
}
