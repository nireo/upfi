package web

import (
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/templates"
)

// ServeSettingsPage serves the user a settings page, in which they can configure their account settings.
// Also does checking if the user is logged in. After all the checking serve a html template, which is used
// to display current user configuration.
func ServeSettingsPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	username := r.Header.Get("username")

	// Find the user from the database, so that we can display the user's current settings.
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	params := templates.SettingsParams{
		User:          user,
		Authenticated: true,
		Title:         "settings",
	}

	// Serve the settings page with the given parameters.
	templates.Settings(w, params)
}

// HandleSettingChange applies the new settings to the new user's database entry. This handler will get
// called by the ServeSettingsPage handler's html template.
func HandleSettingChange(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// The auth token middleware appends the user's username in to the request header, if the
	// execution is successful.
	username := r.Header.Get("username")
	db := lib.GetDatabase()

	// Load the user using the username, so that we can change the settings, and then later
	// save to changes to the database.
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	// Since all upfi forms are handeled using multipart forms, we need to parse the values
	err = r.ParseMultipartForm(1 << 20)
	if err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}

	// Check that the username field exists to prevent an index out of bounds error.
	if len(r.Form["username"]) == 0 {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	// Check if the user has decided to update their username
	newUsername := r.Form["username"][0]
	if !lib.IsUsernameValid(newUsername) {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	// Check that there are no conflicts with an existing user
	if err := db.Where(&models.User{Username: newUsername}).Error; err == nil {
		ErrorPageHandler(w, r, lib.ConflictErrorPage)
		return
	}

	// Update the new username and save the changes to the database
	user.Username = newUsername
	db.Save(&user)

	// Send user status codes which indicate that the request was successful
	http.Redirect(w, r, "/settings", http.StatusMovedPermanently)
}

// DeleteUser handles a total account deletion that includes deleting all information about the user and his/her files.
// Also deletes the user's token so that user can't make requests with a invalid username that doesn't exist.
func DeleteUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// The auth token middleware appends the user's username in to the request header, if the
	// execution is successful.
	username := r.Header.Get("username")

	// Find the user since we need the user struct to delete the user from the database, also we need the
	// user's uuid to delete all of his/her files.
	db := lib.GetDatabase()
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	// Remove the hole directory we created at registeration.
	if err := os.Remove(lib.AddRootToPath("files/") + user.UUID); err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}

	// After all the other things have been deleted, we delete the user entry from the database.
	db.Delete(&user)

	// Remove the user's authentication cookie
	c := &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	}
	http.SetCookie(w, c)
	// Redirect the user to the home page, so they don't get stuck in authorized pages.
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// UpdatePassword handles the update of a user's password given a new password and the old password.
// Also does checking if the current password provided is matching, if not the password won't be updated.
func UpdatePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// The auth token middleware appends the user's username in to the request header, if the
	// execution is successful.
	username := r.Header.Get("username")

	db := lib.GetDatabase()
	// Load the user model since we need it to check the validity of the current password and to
	// update the user model with the new hashed password
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	// Parse the multipart form that is used in the template file, in the settings, that sent the request
	// to this handler.
	err = r.ParseMultipartForm(1 << 20)
	if err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}

	if len(r.Form["password"]) == 0 || len(r.Form["newPassword"]) == 0 {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	// Take the current and new password from the request and do some checking on them.
	currentPassword := r.Form["password"][0]
	newPassword := r.Form["newPassword"][0]

	// We don't need to check the validity of the currentPassword since this password has already
	// been checked when the user registered.
	if !lib.IsPasswordValid(newPassword) {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	// Check that the current password in the form matches the one on the user model.
	if !lib.CheckPasswordHash(currentPassword, user.Password) {
		ErrorPageHandler(w, r, lib.ForbiddenErrorPage)
		return
	}

	// Since all the checking is valid, we hash the new password and the update the password fields on the
	// database entry.
	newHashedPassword, err := lib.HashPassword(newPassword)
	if err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}
	user.Password = newHashedPassword
	db.Save(&user)

	// Redirect the user back to the /settings page, where the request originally came from.
	http.Redirect(w, r, "/settings", http.StatusMovedPermanently)
}
