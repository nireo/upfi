package optimized_api

import (
	"html/template"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

func ServeCreatePage(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./static/create_paste.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

func CreatePaste(ctx *fasthttp.RequestCtx) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		return
	}

	// get values from form
	title := form.Value["title"][0]
	description := form.Value["description"][0]
	content := form.Value["content"][0]

	// load the user
	db := lib.GetDatabase()
	user, err := models.FindOneUser(&models.User{Username: string(ctx.Request.Header.Peek("username"))})
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	newPasteEntry := &models.Paste{
		UserID:      user.ID,
		Title:       title,
		Description: description,
		Content:     content,
		UUID:        lib.GenerateUUID(),
	}

	db.Save(newPasteEntry)

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Redirect("/pastes", fasthttp.StatusMovedPermanently)
}

func DeletePaste(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	pasteID := ctx.UserValue("paste").(string)
	var paste models.Paste
	if err := db.Where(&models.Paste{UUID: pasteID}).First(&paste).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	if paste.UserID != user.ID {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusForbidden), fasthttp.StatusForbidden)
		return
	}

	db.Delete(&paste)

	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
	ctx.Redirect("/pastes", fasthttp.StatusMovedPermanently)
}
