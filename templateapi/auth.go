package templateapi

import (
	"os"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

// ServeRegisterPage returns the register html page to the user.
func ServeRegisterPage(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	ctx.SendFile("./static/register.html")
}

// ServeLoginPage returns the login html page to the user.
func ServeLoginPage(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	ctx.SendFile("./static/login.html")
}

// Register handles the register request from the /register page html form. It creates checks for conflicting
// usernames and creates a folder to the store all of the user's files in. Finally it creates a database entry
// with all the information in given in the form.
func Register(ctx *fasthttp.RequestCtx) {
	// Parse the multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	// Check that the username and the password fields are not empty. If they are empty, return the
	// user with a bad request status.
	if len(form.Value["username"]) == 0 || len(form.Value["password"][0]) == 0 || len(form.Value["master"]) == 0 {
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	// Store the form values into variables, so that the code looks cleaner
	username := form.Value["username"][0]
	password := form.Value["password"][0]

	// The masterPass is the file encryption master password, which is used to encrypt all the files.
	// It's checked during uploading and downloading files. Also it can be same as the normal password,
	// but this isn't as secure as using different passwords.
	masterPass := form.Value["master"][0]

	if !lib.IsUsernameValid(username) || !lib.IsPasswordValid(password) || !lib.IsPasswordValid(masterPass) {
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
	}

	// Check that the username is unique, and if there exists a user with that name return a conflicting status.
	if _, err := models.FindOneUser(&models.User{Username: username}); err == nil {
		ErrorPageHandler(ctx, lib.ConflictErrorPage)
		return
	}

	// Hash the password of the user using bcrypt.
	passwordHash, err := lib.HashPassword(password)
	if err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}

	// Hash the master password using the same hashing as the normal password, so that we can easily
	// check the validity of the password.
	masterHash, err := lib.HashPassword(masterPass)
	if err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
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
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}

	// Finally save that entry after creating the folder, since the folder creation is more likely to fail.
	db := lib.GetDatabase()
	db.Create(&newUser)

	// Create a new authentication token for the user so that he/she can use authenticated routes.
	token, err := lib.CreateToken(newUser.Username)
	if err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
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
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	// Check that the fields exist
	if len(form.Value["username"]) == 0 || len(form.Value["password"]) == 0 {
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	// Store the form fields in variables, so the code is cleaner.
	username := form.Value["username"][0]
	password := form.Value["password"][0]

	// Check that a user with the given username actually exists.
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		// Return a forbiden, since we don't want to tell another user, if some user has an account.
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Compare the hash on the database model to the hash of the given password.
	if !lib.CheckPasswordHash(password, user.Password) {
		ErrorPageHandler(ctx, lib.ForbiddenErrorPage)
		return
	}

	// Create a token for the user, with which the user can use different authenticated routes
	token, err := lib.CreateToken(user.Username)
	if err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
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
