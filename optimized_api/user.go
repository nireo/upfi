package optimized_api

import (
	"fmt"
	"os"
	"text/template"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

// ServeSettingsPage serves the user a settings page, in which they can configure their account settings.
// Also does checking if the user is logged in. After all the checking serve a html template, which is used
// to display current user configuration.
func ServeSettingsPage(ctx *fasthttp.RequestCtx) {
	// Set the right Content-Type so that they html renders corretly.
	ctx.Response.Header.Set("Content-Type", "text/html")

	// The auth token middleware appends the user's username in to the request header, if the
	// execution is successful.
	username := string(ctx.Request.Header.Peek("username"))

	db := lib.GetDatabase()

	// Find the user from the database, so that we can display the user's current settings.
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Serve the template file, with the user information we loaded before.
	tmpl := template.Must(template.ParseFiles("./static/settings_template.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

// HandleSettingsChange applies the new settings to the new user's database entry. This handler will get
// called by the ServeSettingsPage handler's html template.
func HandleSettingChange(ctx *fasthttp.RequestCtx) {
	// The auth token middleware appends the user's username in to the request header, if the
	// execution is successful.
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	// Load the user using the username, so that we can change the settings, and then later
	// save to changes to the database.
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Since all upfi forms are handeled using multipart forms, we need to parse the values
	form, err := ctx.MultipartForm()
	if err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}

	// Check that the username field exists to prevent an index out of bounds error.
	if len(form.Value["username"]) == 0 {
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	// Check if the user has decided to update their username
	newUsername := form.Value["username"][0]
	if !lib.IsUsernameValid(newUsername) {
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	// Check that there are no conflicts with an existing user
	if err := db.Where(&models.User{Username: newUsername}).Error; err == nil {
		ErrorPageHandler(ctx, lib.ConflictErrorPage)
		return
	}

	// Update the new username and save the changes to the database
	user.Username = newUsername
	db.Save(&user)

	// Send user status codes which indicate that the request was successful
	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
	ctx.Redirect("/settings", fasthttp.StatusMovedPermanently)
}

// DeleteUser handles a total account deletion that includes deleting all information about the user and his/her files.
// Also deletes the user's token so that user can't make requests with a invalid username that doesn't exist.
func DeleteUser(ctx *fasthttp.RequestCtx) {
	// The auth token middleware appends the user's username in to the request header, if the
	// execution is successful.
	username := string(ctx.Request.Header.Peek("username"))

	// Find the user since we need the user struct to delete the user from the database, also we need the
	// user's uuid to delete all of his/her files.
	db := lib.GetDatabase()
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Remove the hole directory we created at registeration.
	if err := os.Remove(fmt.Sprintf("./files/%s", user.UUID)); err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}

	// After all the other things have been deleted, we delete the user entry from the database.
	db.Delete(&user)

	// Remove the user's authentication cookie
	ctx.Request.Header.DelAllCookies()

	// Redirect the user to the home page, so they don't get stuck in authorized pages.
	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
	ctx.Redirect("/", fasthttp.StatusMovedPermanently)
}

// UpdatePassword handles the update of a user's password given a new password and the old password.
// Also does checking if the current password provided is matching, if not the password won't be updated.
func UpdatePassword(ctx *fasthttp.RequestCtx) {
	// The auth token middleware appends the user's username in to the request header, if the
	// execution is successful.
	username := string(ctx.Request.Header.Peek("username"))

	db := lib.GetDatabase()
	// Load the user model since we need it to check the validity of the current password and to
	// update the user model with the new hashed password
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Parse the multipart form that is used in the template file, in the settings, that sent the request
	// to this handler.
	form, err := ctx.MultipartForm()
	if err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}

	if len(form.Value["password"]) == 0 || len(form.Value["newPassword"]) == 0 {
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	// Take the current and new password from the request and do some checking on them.
	currentPassword := form.Value["password"][0]
	newPassword := form.Value["newPassword"][0]

	// We don't need to check the validity of the currentPassword since this password has already
	// been checked when the user registered.
	if !lib.IsPasswordValid(newPassword) {
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	// Check that the current password in the form matches the one on the user model.
	if !lib.CheckPasswordHash(currentPassword, user.Password) {
		ErrorPageHandler(ctx, lib.ForbiddenErrorPage)
		return
	}

	// Since all the checking is valid, we hash the new password and the update the password fields on the
	// database entry.
	newHashedPassword, err := lib.HashPassword(newPassword)
	if err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}
	user.Password = newHashedPassword
	db.Save(&user)

	// Redirect the user back to the /settings page, where the request originally came from.
	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
	ctx.Redirect("/settings", fasthttp.StatusMovedPermanently)
}
