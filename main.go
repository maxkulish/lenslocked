package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprint(w, "To get in touch, please send an email to <"+
		"a href=\"mailto:m@example.com\">m@example.com")
}

func faq(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprint(w, "<h1>This is my FAQ page</h1>")

}

func notFound404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprint(w, "I can't find this page")
}

func main() {
	r := httprouter.New()
	r.GET("/", home)
	r.GET("/contact", contact)
	r.GET("/faq", faq)

	handler404 := http.HandlerFunc(notFound404)
	r.NotFound = handler404

	_ = http.ListenAndServe(":3000", r)
}
