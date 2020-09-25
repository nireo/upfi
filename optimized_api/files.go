package optimized_api

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/server"
	"github.com/valyala/fasthttp"
)

func ServeUploadPage(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./static/upload.html"))
	err := tmpl.Execute(ctx, nil)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

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

	username := string(ctx.Request.Header.Peek("username"))

	db := lib.GetDatabase()
	// find user
	user, err := models.FindOneUser(&models.User{Username: username})
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

	db.Create(newFileEntry)

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Redirect("/files", fasthttp.StatusMovedPermanently)
}

func GetSingleFile(ctx *fasthttp.RequestCtx) {
	// get the user's username which was appended to the request header
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	// the parameter is given by fasthttprouter instead of fasthttp!
	fileID := ctx.UserValue("file").(string)
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	// check for file ownership
	if user.ID != file.UserID {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusForbidden), fasthttp.StatusForbidden)
		return
	}

	ctx.Response.Header.Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./templates/single_file_template.html"))
	if err := tmpl.Execute(ctx, file); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

func GetUserFiles(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	var files []models.File
	db.Find(&files).Where(&models.File{UserID: user.ID})

	tmpl := template.Must(template.ParseFiles("./templates/files_template.html"))
	data := server.FilePage{
		PageTitle: "Your files",
		Files:     files,
	}

	ctx.Response.Header.Set("Content-Type", "text/html")
	if err := tmpl.Execute(ctx, data); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

func DeleteFile(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	fileID := ctx.UserValue("file").(string)
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	if user.ID != file.UserID {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusForbidden), fasthttp.StatusForbidden)
		return
	}

	os.Remove("./files/" + user.UUID + "/" + fmt.Sprintf("%s%s", file.UUID, file.Extension))
	db.Delete(&file)

	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
	ctx.Redirect("/files", fasthttp.StatusMovedPermanently)
}

func UpdateFile(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error("Content type needs to be multipart", fasthttp.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	fileID := ctx.UserValue("file").(string)
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	if user.ID != file.UserID {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusForbidden), fasthttp.StatusForbidden)
		return
	}

	title := form.Value["title"][0]
	description := form.Value["description"][0]

	if description != "" {
		file.Description = description
	}

	if title != "" {
		file.Filename = title
	}

	db.Save(&file)

	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
	ctx.Redirect("/files", fasthttp.StatusMovedPermanently)
}
