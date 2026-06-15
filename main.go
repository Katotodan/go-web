package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Katotodan/go-web/controller"
	"github.com/Katotodan/go-web/db"
	"github.com/gorilla/mux"
)

func main() {
	err := db.ConnectDb()
	if err != nil {
		log.Fatalf("Critical Error: Could not connect to database: %v", err)
	}
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("static"))

	// Getting file from static folder
	// r.Handle("/static/", http.StripPrefix("/static/", fs))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/about.html")
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, youve requested: %s\n", r.URL.Path)
	})
	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)

	})

	// Create table handler
	r.HandleFunc("/new-table/", controller.CreateTable).Methods("POST")
	r.HandleFunc("/delete-table/{tableName}", controller.DropTable).Methods("DELETE")
	r.HandleFunc("/all/table", controller.GetAllTable).Methods("GET")
	r.HandleFunc("/insert/user", controller.InsertUser).Methods("POST")

	http.ListenAndServe(":3000", r)
}
