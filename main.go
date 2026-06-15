package main

import (
	"fmt"
	"net/http"
)

func main() {

	fs := http.FileServer(http.Dir("static"))

	// Getting file from static folder
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/about.html")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, youve requested: %s\n", r.URL.Path)
	})

	http.ListenAndServe(":3000", nil)
}
