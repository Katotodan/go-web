package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func Secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	fmt.Println("Session for secret")
	fmt.Println(session)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	fmt.Fprintln(w, "The cake is a lie!")
}

func Login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Auth goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	fmt.Println("Session")
	fmt.Println(session)
	session.Save(r, w)
	fmt.Println("Session after saving")
	fmt.Println(session)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

func Logging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Middleware called on " + r.URL.Path)
		h(w, r)
	}
}

func Foo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "foo")
}

func Bar(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "bar")
}
