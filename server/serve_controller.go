package server

import (
	"html/template"
	"net/http"

	"github.com/nireo/upfi/lib"
)

type AuthenticatedPage struct {
	Authenticated bool
}

func ServeHomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	store := lib.GetStore()
	session, _ := store.Get(r, "auth")

	// check if the user authenticated, so we can display the right navbar
	authenticated := true
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		authenticated = false
	}

	tmpl := template.Must(template.ParseFiles("./static/home.html"))
	if err := tmpl.Execute(w, &AuthenticatedPage{Authenticated: authenticated}); err != nil {
		lib.HttpInternalErrorHandler(w, r)
		return
	}
}
