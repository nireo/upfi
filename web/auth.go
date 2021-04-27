package web

import (
	"net/http"
	"os"
	"time"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/templates"
)

// ServeRegisterPage returns the register html page to the user.
func ServeRegisterPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	templates.Register(w, templates.RegisterParams{
		Title:         "register",
		Authenticated: lib.IsAuth(r),
	})
}

// ServeLoginPage returns the login html page to the user.
func ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	templates.Login(w, templates.LoginParams{
		Authenticated: lib.IsAuth(r),
		Title:         "login",
	})
}

// Register handles the register request from the /register page html form. It creates checks for conflicting
// usernames and creates a folder to the store all of the user's files in. Finally it creates a database entry
// with all the information in given in the form.
func Register(w http.ResponseWriter, r *http.Request) {
	// check if the user is already logged in
	if lib.IsAuth(r) {
		return
	}

	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check that the username and the password fields are not empty. If they are empty, return the
	// user with a bad request status.
	if len(r.Form["username"]) == 0 || len(r.Form["password"]) == 0 || len(r.Form["master"]) == 0 {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	// Store the form values into variables, so that the code looks cleaner
	username := r.Form["username"][0]
	password := r.Form["password"][0]

	// The masterPass is the file encryption master password, which is used to encrypt all the files.
	// It's checked during uploading and downloading files. Also it can be same as the normal password,
	// but this isn't as secure as using different passwords.
	masterPass := r.Form["master"][0]

	if !lib.IsUsernameValid(username) || !lib.IsPasswordValid(password) || !lib.IsPasswordValid(masterPass) {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	// Check that the username is unique, and if there exists a user with that name return a conflicting status.
	if _, err := models.FindOneUser(&models.User{Username: username}); err == nil {
		ErrorPageHandler(w, r, lib.ConflictErrorPage)
		return
	}

	// Hash the password of the user using bcrypt.
	passwordHash, err := lib.HashPassword(password)
	if err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}

	// Hash the master password using the same hashing as the normal password, so that we can easily
	// check the validity of the password.
	masterHash, err := lib.HashPassword(masterPass)
	if err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}

	// Create the database entry for the user, which contains the username, password and a newly generated unique id
	newUser := models.User{
		Username:             username,
		Password:             passwordHash,
		FileEncryptionMaster: masterHash,
		UUID:                 lib.GenerateUUID(),
	}

	// Use that unique id to create a folder in the files directory that in the future will contain all of the
	// user's files.
	err = os.Mkdir(lib.AddRootToPath("files/")+newUser.UUID, os.ModePerm)
	if err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}

	// Finally save that entry after creating the folder, since the folder creation is more likely to fail.
	db := lib.GetDatabase()
	db.Create(&newUser)

	// Create a new authentication token for the user so that he/she can use authenticated routes.
	token, err := lib.CreateToken(newUser.Username)
	if err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}

	// Use the token we created before and store it in a cookie, which will be checked when accessing
	// authenticated routes.
	expirationTime := time.Now().Add(time.Hour * 24)
	cookie := http.Cookie{Name: "token", Value: token, Expires: expirationTime}
	http.SetCookie(w, &cookie)

	// Redirect the new user to the files page where the user can add new files.
	http.Redirect(w, r, "/files", http.StatusPermanentRedirect)
}
