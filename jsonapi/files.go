package jsonapi

import (
	"encoding/json"
	"os"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

func GetSingleFile(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	fileID := ctx.UserValue("file").(string)
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	if user.ID != file.UserID {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	lib.WriteResponseJSON(ctx, fasthttp.StatusOK, file)
}

type updateFileBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func UpdateFile(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	var body updateFileBody
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		ServeErrorJSON(ctx, lib.BadRequestErrorPage)
		return
	}

	title := body.Title
	description := body.Description

	// The title cannot be empty
	if title == "" || len(title) > 20 {
		ServeErrorJSON(ctx, lib.BadRequestErrorPage)
		return
	}

	if len(description) > 200 {
		ServeErrorJSON(ctx, lib.BadRequestErrorPage)
		return
	}

	fileID := ctx.UserValue("file").(string)
	db := lib.GetDatabase()

	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	if user.ID != file.UserID {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	file.Description = description
	file.Filename = title

	db.Save(&file)

	lib.WriteResponseJSON(ctx, fasthttp.StatusOK, file)
}

func GetUserFiles(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	var files []models.File
	db.Where(&models.File{UserID: user.ID}).Find(&files)
	lib.WriteResponseJSON(ctx, fasthttp.StatusOK, files)
}

func DeleteFile(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	file, err := models.FindOneFile(&models.File{UUID: ctx.UserValue("file").(string)})
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	if user.ID != file.ID {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	if err := os.Remove("./files/" + user.UUID + "/" + file.UUID + file.Extension); err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	db := lib.GetDatabase()
	db.Delete(&file)

	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}
