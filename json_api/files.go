package json_api

import (
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
