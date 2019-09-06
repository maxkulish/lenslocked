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

	r := mux.NewRouter()
	staticC := controller.NewStatic()
	usersC := controller.NewUser(services.User)
	galleriesC := controller.NewGalleries(services.Gallery, r)

	userMW := middleware.User{UserService: services.User}
	requireUserMW := middleware.RequireUser{
		User: userMW,
	}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	// Gallery routes
	r.Handle("/galleries/new", requireUserMW.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMW.ApplyFn(galleriesC.Index)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMW.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMW.ApplyFn(galleriesC.Edit)).
		Methods("GET").
		Name(controller.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMW.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMW.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).
		Methods("GET").
		Name(controller.ShowGallery)

	fmt.Println("Starting the server on :3000")

	handler404 := http.HandlerFunc(notFound404)
	r.NotFoundHandler = handler404

	_ = http.ListenAndServe(":3000", userMW.Apply(r))
}
