package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		_, _ = fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
	} else if r.URL.Path == "/contact" {
		_, _ = fmt.Fprint(w, "To get in touch, please send an email to <"+
			"a href=\"mailto:m@example.com\">m@example.com")
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "<h1>We could not find the page you were looking for :( ")
	}

}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name/spanish", Hello)
	_ = http.ListenAndServe(":3000", router)
}
