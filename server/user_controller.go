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

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

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

	currentPassword := r.FormValue("password")
	newPassword := r.FormValue("newPassword")
	// check if the password is correct
	if !CheckPasswordHash(currentPassword, user.Password) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	newHashedPassword, err := HashPassword(newPassword)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user.Password = newHashedPassword
	db.Save(&user)
	http.Redirect(w, r, "http://localhost:8080/settings", http.StatusMovedPermanently)
}
