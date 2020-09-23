package optimized_api

import (
	"fmt"
	"path/filepath"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

func UploadFile(ctx *fasthttp.RequestCtx) {
	// get file
	header, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error("File could not be parsed", fasthttp.StatusInternalServerError)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error("Content type needs to be multipart", fasthttp.StatusBadRequest)
		return
	}

	db := lib.GetDatabase()
	// find user
	user, err := models.FindOneUser(&models.User{Username: form.Value["username"][0]})
	if err != nil {
		ctx.Error("User not found", fasthttp.StatusNotFound)
		return
	}

	newFileEntry := &models.File{
		Filename:    header.Filename,
		UUID:        lib.GenerateUUID(),
		Description: form.Value["description"][0],
		Size:        header.Size,
		UserID:      user.ID,
		Extension:   filepath.Ext(header.Filename),
	}

	fileDirectory := fmt.Sprintf("./files/%s/%s%s", user.UUID, newFileEntry.UUID, newFileEntry.Extension)
	if err := fasthttp.SaveMultipartFile(header, fileDirectory); err != nil {
		ctx.Error("File could not be stored", fasthttp.StatusInternalServerError)
		return
	}

	db.NewRecord(newFileEntry)
	db.Create(newFileEntry)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}
