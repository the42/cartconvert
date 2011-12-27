package main

import (
	// "github.com/the42/cartconvert"
	// "io"
	"net/http"
	// "old/template"
)

// Better read it from an INI-File
const staticWebDir = "../static/"
const templateDir = "../templates/"

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticWebDir+r.URL.Path)
}

func main() {

	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":1111", nil)
}
