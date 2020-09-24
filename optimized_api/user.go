package optimized_api

import (
	"fmt"
	"os"
	"text/template"

	"github.com/nireo/booru/lib"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/server"
	"github.com/valyala/fasthttp"
)

func ServeSettingsPage(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")

	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("./static/settings_template.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

func HandleSettingChange(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	newUsername := form.Value["username"][0]
	if newUsername != "" {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		return
	}

	user.Username = newUsername

	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}

func DeleteUser(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	// remove the user files
	if err := os.Remove(fmt.Sprintf("./files/%s", user.UUID)); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	db.Delete(&user)

	// remove the tokens
	ctx.Request.Header.DelAllCookies()
	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
	ctx.Write([]byte("Successfully removed all your files!"))
}

func UpdatePassword(ctx *fasthttp.RequestCtx) {
	db := lib.GetDatabase()
	username := string(ctx.Request.Header.Peek("username"))

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	// still require the old password before we can change it
	currentPassword := form.Value["password"][0]
	newPassword := form.Value["newPassword"][0]

	if !server.CheckPasswordHash(currentPassword, user.Password) {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusForbidden), fasthttp.StatusForbidden)
		return
	}

	newHashedPassword, err := server.HashPassword(newPassword)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	user.Password = newHashedPassword
	db.Save(&user)

	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
	ctx.Write([]byte("Successfully changed your password!"))
}
