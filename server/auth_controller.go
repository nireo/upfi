package server

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"

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
	if r.Method != http.MethodPost {
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	db := lib.GetDatabase()
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		fmt.Fprintf(w, "User not found")
		return
	}

	if !checkPasswordHash(password, user.Password) {
		fmt.Fprintf(w, "Incorrect credentials")
		return
	}

	fmt.Fprintf(w, "Successfully logged in!")
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
		fmt.Fprintf(w, "User with username %s already exists", username)
		return
	}

	hash, err := hashPassword(password)
	if err != nil {
		fmt.Fprintf(w, "Internal server error")
		return
	}


	newUser := models.User{
		Username: username,
		Password: hash,
		UUID:     lib.GenerateUUID(),
	}

	db.NewRecord(newUser)
	db.Create(&newUser)

	fmt.Fprintf(w, "Successfully registered!")
}
