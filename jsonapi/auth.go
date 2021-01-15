package jsonapi

import (
	"encoding/json"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
	"os"
)

type registerRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Master   string `json:"master"`
}

// Register handles the registeration. It creates checks for conflicting usernames and creates a folder
// to the store all of the user's files in. Finally it creates a database entry with all the information
// in given json.
func Register(ctx *fasthttp.RequestCtx) {
	var body registerRequestBody
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		// Since the parsing didn't work the request body doesn't have all the needed fields
		ServeErrorJSON(ctx, lib.BadRequestErrorPage)
		return
	}

	username := body.Username
	password := body.Password

	// The masterPass is the file encryption master password, which is used to encrypt all the files.
	// It's checked during uploading and downloading of files. Also it can be the same as the normal password,
	// but this isn't as secure as using different passwords.
	masterPass := body.Master

	if !lib.IsUsernameValid(username) || !lib.IsPasswordValid(password) || !lib.IsPasswordValid(masterPass) {
		ServeErrorJSON(ctx, lib.BadRequestErrorPage)
		return
	}

	// Check that the username is unique, and if there exists a user with that name return a conflicting status.
	if _, err := models.FindOneUser(&models.User{Username: username}); err == nil {
		ServeErrorJSON(ctx, lib.ConflictErrorPage)
		return
	}

	// Hash the password of the user using bcrypt.
	passwordHash, err := lib.HashPassword(password)
	if err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	// Hash the master password using the same hashing as the normal password, so that we can easily
	// check the validity of the password.
	masterHash, err := lib.HashPassword(masterPass)
	if err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
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
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	// Finally save that entry after creating the folder, since the folder creation is more likely to fail.
	db := lib.GetDatabase()
	db.Create(&newUser)

	// Create a new authentication token for the user so that he/she can use authenticated routes.
	token, err := lib.CreateToken(newUser.Username)
	if err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	// Use the token we created before and store it in a cookie, which will be checked when accessing
	// authenticated routes.
	var cookie fasthttp.Cookie
	cookie.SetKey("token")
	cookie.SetValue(token)
	ctx.Response.Header.SetCookie(&cookie)

	lib.WriteResponseJSON(ctx, fasthttp.StatusOK, newUser)
}

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login handles the login request. It firstly checks that the a user with the given username does exist
// and then checks that user's hash using bcrypt to the password given in the form.
func Login(ctx *fasthttp.RequestCtx) {
	var body loginRequestBody
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		// Since parsing the json wasn't successful, the body doesn't have all the required fields.
		ServeErrorJSON(ctx, lib.BadRequestErrorPage)
		return
	}

	// Store the form fields in variables, so the code is cleaner.
	username := body.Username
	password := body.Password

	// Check that a user with the given username actually exists.
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	// Compare the hash on the database model to the hash of the given password.
	if !lib.CheckPasswordHash(password, user.Password) {
		ServeErrorJSON(ctx, lib.ForbiddenErrorPage)
		return
	}

	// Create a token for the user, with which the user can use different authenticated routes
	token, err := lib.CreateToken(user.Username)
	if err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	// Store the token in a cookie, which the authentication middleware checks when
	// accessing authenticated routes.
	var cookie fasthttp.Cookie
	cookie.SetKey("token")
	cookie.SetValue(token)
	ctx.Response.Header.SetCookie(&cookie)
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
}
