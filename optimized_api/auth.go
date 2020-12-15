package optimized_api

import (
	"html/template"
	"os"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

// ServeRegisterPage returns the requester with the html of the register page, but even though the
// page servead, is a static page it's given out as a template, so that we can more easily add data in
// the future, and the implementation is quite minimal.
func ServeRegisterPage(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./static/register.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

// ServeLoginPage returns the requester with the html of the login page, but even though the
// page served, is a static page it's given out as a template, so that we can more easily add data in
// the future, and the implementation is quite minimal.
func ServeLoginPage(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./static/login.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

// Register handles the register request from the /register page html form. It creates checks for conflicting
// usernames and creates a folder to the store all of the user's files in. Finally it creates a database entry
// with all the information in given in the form.
func Register(ctx *fasthttp.RequestCtx) {
	// Parse the multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error("Content type needs to be multipart", fasthttp.StatusBadRequest)
		return
	}

	// Check that the username and the password fields are not empty. If they are empty, return the
	// user with a bad request status.
	if len(form.Value["username"]) == 0 || len(form.Value["password"][0]) == 0 {
		ctx.Error("Both username and password fields must be added", fasthttp.StatusBadRequest)
		return
	}

	// Store the form values into variables, so that the code looks cleaner
	username := form.Value["username"][0]
	password := form.Value["password"][0]

	// The masterPass is the file encryption master password, which is used to encrypt all the files.
	// It's checked during uploading and downloading files. Also it can be same as the normal password,
	// but this isn't as secure as using different passwords.
	masterPass := form.Value["master"][0]

	// Check that the username is unique, and if there exists a user with that name return a conflicting status.
	if _, err := models.FindOneUser(&models.User{Username: username}); err != nil {
		ctx.Error("User already exists with that username", fasthttp.StatusConflict)
		return
	}

	// Hash the password of the user using bcrypt.
	passwordHash, err := lib.HashPassword(password)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	// Hash the master password using the same hashing as the normal password, so that we can easily
	// check the validity of the password.
	masterHash, err := lib.HashPassword(masterPass)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	// Create the database entry for the user, which contains the username, password and a newly generated unique id
	newUser := models.User{
		Username:             username,
		Password:             passwordHash,
		FileEncryptionMaster: masterHash,
		UUID:                 lib.GenerateUUID(),
	}

	// Use that unique id to create a folder in the ./files directory that in the future will contain all of the
	// user's files.
	err = os.Mkdir("./files/"+newUser.UUID, os.ModePerm)
	if err != nil {
		ctx.Error("Failed user directory creation", fasthttp.StatusInternalServerError)
		return
	}

	// Finally save that entry after creating the folder, since the folder creation is more likely to fail.
	db := lib.GetDatabase()
	db.Create(&newUser)

	// Create a new authentication token for the user so that he/she can use authenticated routes.
	token, err := lib.CreateToken(newUser.Username)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	// Use the token we created before and store it in a cookie, which will be checked when accessing
	// authenticated routes.
	var cookie fasthttp.Cookie
	cookie.SetKey("token")
	cookie.SetValue(token)
	ctx.Response.Header.SetCookie(&cookie)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	// Redirect the new user to the files page where the user can add new files.
	ctx.Redirect("/files", fasthttp.StatusMovedPermanently)
}

// Login handles the login request from the /login page. It firstly checks that the a user
// with the given username does exist and then checks that user's hash using bcrypt to the
// password given in the form.
func Login(ctx *fasthttp.RequestCtx) {
	// Parse the multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error("Content type needs to be multipart", fasthttp.StatusBadRequest)
		return
	}

	// Store the form fields in variables, so the code is cleaner.
	username := form.Value["username"][0]
	password := form.Value["password"][0]

	// Check that a user with the given username actually exists.
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	// Compare the hash on the database model to the hash of the given password.
	if !lib.CheckPasswordHash(password, user.Password) {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusForbidden), fasthttp.StatusForbidden)
		return
	}

	// Create a token for the user, with which the user can use different authenticated routes
	token, err := lib.CreateToken(user.Username)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	// Store the token in a cookie, which the authentication middleware checks when
	// accessing authenticated routes.
	var cookie fasthttp.Cookie
	cookie.SetKey("token")
	cookie.SetValue(token)
	ctx.Response.Header.SetCookie(&cookie)
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)

	// Redirect the user to the files page.
	ctx.Redirect("/files", fasthttp.StatusMovedPermanently)
}
