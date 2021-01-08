package json_api

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
	Master string `json:"master"`
}

func Register(ctx *fasthttp.RequestCtx) {
	var body registerRequestBody
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		// Since the parsing didn't work the request body doesn't have all the needed fields
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		return
	}

	username := body.Username
	password := body.Password

	// The masterPass is the file encryption master password, which is used to encrypt all the files.
	// It's checked during uploading and downloading of files. Also it can be the same as the normal password,
	// but this isn't as secure as using different passwords.
	masterPass := body.Master

	if len(username) < 3 || len(password) < 8 || len(masterPass) < 8 {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		return
	}

	if len(username) > 20 || len(password) > 32 || len(masterPass) > 32 {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		return
	}

	// Check that the username is unique, and if there exists a user with that name return a conflicting status.
	if _, err := models.FindOneUser(&models.User{Username: username}); err == nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusConflict), fasthttp.StatusConflict)
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
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
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
}
