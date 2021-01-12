package jsonapi

import (
	"encoding/json"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

// WhoAmI handlers checks for a token and if a token was found return information about the user that has the token.
// This is used to set the user state in the front-end and have it up-to-date.
func WhoAmI(ctx *fasthttp.RequestCtx) {
	user, err := models.FindOneUser(&models.User{Username: string(ctx.Request.Header.Peek("username"))})
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	lib.WriteResponseJSON(ctx, fasthttp.StatusOK, user)
}

type handleSettingsChangeBody struct {
	Username string `json:"username"`
}

// HandleSettingsChange handles the change of the username. Also checks for conflicts and other problems
// with the new error.
func HandleSettingsChange(ctx *fasthttp.RequestCtx) {
	user, err := models.FindOneUser(&models.User{Username: string(ctx.Request.Header.Peek("username"))})
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	var body handleSettingsChangeBody
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		ServeErrorJSON(ctx, lib.BadRequestErrorPage)
		return
	}

	username := body.Username
	if !lib.IsUsernameValid(username) {
		ServeErrorJSON(ctx, lib.BadRequestErrorPage)
		return
	}

	// Check for password conflicts
	if _, err := models.FindOneUser(&models.User{Username: username}); err == nil {
		ServeErrorJSON(ctx, lib.ConflictErrorPage)
		return
	}

	// Since the password passed all the checks change it.
	user.Username = username

	db := lib.GetDatabase()
	db.Save(&user)

	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}

type updatePasswordBody struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// UpdatePassword is a handler that takes in the user's password and a new password and then updates the user's
// password with the new password. Also does all the needed checking on the new password and checks the old password.
func UpdatePassword(ctx *fasthttp.RequestCtx) {
	user, err := models.FindOneUser(&models.User{Username: string(ctx.Request.Header.Peek("username"))})
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	var body updatePasswordBody
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		ServeErrorJSON(ctx, lib.BadRequestErrorPage)
		return
	}

	currentPassword := body.CurrentPassword
	newPassword := body.NewPassword

	// We don't need to check the validity of the currentPassword since this password has already
	// been checked when the user registered.
	if !lib.IsPasswordValid(newPassword) {
		ServeErrorJSON(ctx, lib.BadRequestErrorPage)
		return
	}

	// Check that the current password in the form matches the one on the user model.
	if !lib.CheckPasswordHash(currentPassword, user.Password) {
		ServeErrorJSON(ctx, lib.ForbiddenErrorPage)
		return
	}

	// Since all the checking is valid, we hash the new password and the update the password fields on the
	// database entry.
	newHashedPassword, err := lib.HashPassword(newPassword)
	if err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}
	user.Password = newHashedPassword

	db := lib.GetDatabase()
	db.Save(&user)

	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}
