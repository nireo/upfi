package optimized_api

import (
	"text/template"

	"github.com/nireo/booru/lib"
	"github.com/nireo/upfi/models"
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
