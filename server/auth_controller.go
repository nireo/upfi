package server

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"os"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AuthLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		store := lib.GetStore()
		session, _ := store.Get(r, "auth")
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := models.FindOneUser(&models.User{Username: username})
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if !CheckPasswordHash(password, user.Password) {
			http.Error(w, "Incorrect credentials", http.StatusForbidden)
			return
		}

		session.Values["username"] = username
		session.Values["authenticated"] = true
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "Error saving session", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "http://localhost:8080/files", http.StatusMovedPermanently)
	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles("./static/login.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func AuthRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		db := lib.GetDatabase()
		username := r.FormValue("username")
		password := r.FormValue("password")

		_, err := models.FindOneUser(&models.User{Username: username})
		if err == nil {
			http.Error(w, fmt.Sprintf("User with username %s already exists", username), http.StatusConflict)
			return
		}

		hash, err := HashPassword(password)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		newUser := models.User{
			Username: username,
			Password: hash,
			UUID:     lib.GenerateUUID(),
		}

		// create a directory which stores all of the user's files
		err = os.Mkdir("./files/"+newUser.UUID, os.ModePerm)
		if err != nil {
			http.Error(w, "Failed user directory creation", http.StatusInternalServerError)
			return
		}

		db.Create(&newUser)
		http.Redirect(w, r, "http://localhost:8080/login.html", http.StatusMovedPermanently)
	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles("./static/register.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}
