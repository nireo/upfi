package server

import (
	"fmt"
	"net/http"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
)

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

	if user.Password != password {
		fmt.Fprintf(w, "Incorrect credentials")
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded file")
}

func AuthRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	db := lib.GetDatabase()

	username := r.FormValue("username")

	var exists models.User
	if err := db.Where("username = ?", username).First(&exists); err == nil {
		fmt.Fprintf(w, "User with username %s already exists", username)
		return
	}

	newUser := models.User{
		Username: username,
		Password: r.FormValue("password"),
		UUID:     lib.GenerateUUID(),
	}

	db.NewRecord(newUser)
	db.Create(&newUser)

	fmt.Fprintf(w, "Successfully registered")
}
