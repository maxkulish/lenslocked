package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	_ = http.ListenAndServe(":3000", nil)
}
