package server

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AuthLogin(w http.ResponseWriter, r *http.Request) {
	store := lib.GetStore()
	session, _ := store.Get(r, "auth")
	if r.Method != http.MethodPost {
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	db := lib.GetDatabase()
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !checkPasswordHash(password, user.Password) {
		http.Error(w, "Incorrect credentials", http.StatusForbidden)
		return
	}

	session.Values["username"] = username
	session.Values["authenticated"] = true
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "http://localhost:8080/files", http.StatusMovedPermanently)
}

func AuthRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	db := lib.GetDatabase()
	username := r.FormValue("username")
	password := r.FormValue("password")

	var exists models.User
	if err := db.Where("username = ?", username).First(&exists); err == nil {
		http.Error(w, fmt.Sprintf("User with username %s already exists", username), http.StatusConflict)
		return
	}

	hash, err := hashPassword(password)
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

	db.NewRecord(newUser)
	db.Create(&newUser)
	http.Redirect(w, r, "http://localhost:8080/login.html", http.StatusMovedPermanently)
}
