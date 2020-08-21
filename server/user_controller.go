package server

import (
	"text/template"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"

	"net/http"
)

func SettingsPage(w http.ResponseWriter, r *http.Request) {
	store := lib.GetStore()
	db := lib.GetDatabase()
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		lib.HttpForbiddenHandler(w, r)
		return
	}

	user, err := models.FindOneUser(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		lib.HttpForbiddenHandler(w, r)
		return
	}
	switch r.Method {
	case http.MethodPost:
		user.Username = r.FormValue("username")
		db.Save(&db)
	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles("./templates/settings_template.html"))
		if err := tmpl.Execute(w, r); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
