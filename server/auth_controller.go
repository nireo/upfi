package server

import (
	"fmt"
	"net/http"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Successfully Uploaded file")
}

func AuthRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	newUser := models.User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		UUID:     "secret_id_yep",
	}

	db := lib.GetDatabase()
	db.Create(newUser)
	fmt.Fprintf(w, "Successfully registered")
}
