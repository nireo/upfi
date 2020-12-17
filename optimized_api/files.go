package optimized_api

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/nireo/upfi/crypt"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

// ServeUploadPage serves the requester a upload form, in which the user can upload files to their account.
// Also the route is protected, so that the security token is checked before calling this handler.
func ServeUploadPage(ctx *fasthttp.RequestCtx) {
	// Set the right Content-Type so that the html renders correctly.
	ctx.Response.Header.Set("Content-Type", "text/html")

	// Return the template, which just has a request form.
	tmpl := template.Must(template.ParseFiles("./static/upload.html"))
	err := tmpl.Execute(ctx, nil)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

// UploadFile handles the logic of uploading a file from the upload file form.
// Also the route is protected, so that the security token is checked before calling this handler.
func UploadFile(ctx *fasthttp.RequestCtx) {
	// Get the file from the request form.
	header, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error("File could not be parsed", fasthttp.StatusInternalServerError)
		return
	}

	// Parse the multipart form so that we can check for other values, such as custom filenames or descriptions.
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error("Content type needs to be multipart", fasthttp.StatusBadRequest)
		return
	}

	// Find the user who requested this handler.
	db := lib.GetDatabase()
	username := string(ctx.Request.Header.Peek("username"))
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ctx.Error("User not found", fasthttp.StatusNotFound)
		return
	}

	// Check that the user's master passwords is correct.
	if !lib.CheckPasswordHash(form.Value["master"][0], user.FileEncryptionMaster) {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusForbidden), fasthttp.StatusForbidden)
		return
	}

	// Construct a database entry
	newFileEntry := &models.File{
		Filename:    header.Filename,
		UUID:        lib.GenerateUUID(),
		Description: form.Value["description"][0],
		Size:        header.Size,
		UserID:      user.ID,
		Extension:   filepath.Ext(header.Filename),
	}

	// Save the new file to the user's own file folder.
	path := fmt.Sprintf("./files/%s/%s%s", user.UUID, newFileEntry.UUID, newFileEntry.Extension)
	/*
		if err := fasthttp.SaveMultipartFile(header, path); err != nil {
			InternalServerErrorHandler(ctx)
			return
		}
	*/

	multipartFile, err := header.Open()
	if err != nil {
		InternalServerErrorHandler(ctx)
		return
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, multipartFile); err != nil {
		InternalServerErrorHandler(ctx)
		return
	}

	if err := crypt.EncryptToDst(path, buf.Bytes(), form.Value["master"][0]); err != nil {
		InternalServerErrorHandler(ctx)
		return
	}

	db.Create(newFileEntry)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Redirect("/files", fasthttp.StatusMovedPermanently)
}

// GetSingleFile returns the database entry, which contains data about a file to the user. The user
// needs to provide a file id as a query. Also the files are kept private, so you need to own the file.
// Also the route is protected, so that the security token is checked before calling this handler.
func GetSingleFile(ctx *fasthttp.RequestCtx) {
	// get the user's username which was appended to the request header
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	// Find the user's database entry who is requesting this handler.
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	// Find the file
	fileID := ctx.UserValue("file").(string)
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	// Check that the user owns the file.
	if user.ID != file.UserID {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusForbidden), fasthttp.StatusForbidden)
		return
	}

	// Display the user with the file's information, this template also includes the option to download a file.
	ctx.Response.Header.Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./templates/single_file_template.html"))
	if err := tmpl.Execute(ctx, file); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

type FilePage struct {
	PageTitle string
	Files     []models.File
}

// GetUserFiles returns all the files that are related to the username which is requesting this
// handler. Then handler finds all the related files and constructs a template, which the user
// then can view as html content.
// Also the route is protected, so that the security token is checked before calling this handler.
func GetUserFiles(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	// Find the user's database entry who is requesting this handler.
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	// Find all file models which are related to the user in the database.
	var files []models.File
	db.Find(&files).Where(&models.File{UserID: user.ID})

	tmpl := template.Must(template.ParseFiles("./templates/files_template.html"))
	// construct a struct which contains the data we will give to the html template.
	data := FilePage{
		PageTitle: "Your files",
		Files:     files,
	}

	// Display the user's files to the user
	ctx.Response.Header.Set("Content-Type", "text/html")
	if err := tmpl.Execute(ctx, data); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

// DeleteFile is a handler that deletes a file owned by the user. The handler takes a file id as a query parameter
// and then does checking on the ownership of the file.
// Also the route is protected, so that the security token is checked before calling this handler.
func DeleteFile(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	// Find the database entry of the user that requested this handler.
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		NotFoundHandler(ctx)
		return
	}

	// Find the file, if the file does not exist, return a not found error
	fileID := ctx.UserValue("file").(string)
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		NotFoundHandler(ctx)
		return
	}

	// Check that the user owns the file.
	if user.ID != file.UserID {
		ForbiddenHandler(ctx)
		return
	}

	// Remove the file, if the file cannot be removed the return a internal server error to the user.
	if err := os.Remove("./files/" + user.UUID + "/" + fmt.Sprintf("%s%s", file.UUID, file.Extension)); err != nil {
		InternalServerErrorHandler(ctx)
		return
	}

	// Delete the database entry
	db.Delete(&file)

	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
	ctx.Redirect("/files", fasthttp.StatusMovedPermanently)
}

// UpdateFile is http handler which takes a file id as a query parameter and checks for the file's ownership.
// This handler can be used to update file title and description.
// Also the route is protected, so that the security token is checked before calling this handler.
func UpdateFile(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	// Parse the multipart form so that we can take the 'title' and 'description' fields.
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error("Content type needs to be multipart", fasthttp.StatusBadRequest)
		return
	}

	// Find the user that is requesting this handler.
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		NotFoundHandler(ctx)
		return
	}

	// Find the file
	fileID := ctx.UserValue("file").(string)
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		NotFoundHandler(ctx)
		return
	}

	// Check that the user owns the file
	if user.ID != file.UserID {
		ForbiddenHandler(ctx)
		return
	}

	title := form.Value["title"][0]
	description := form.Value["description"][0]

	// Update only the fields, which are not empty.
	if description != "" {
		file.Description = description
	}

	if title != "" {
		file.Filename = title
	}

	// Save the changes to the database.
	db.Save(&file)

	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
	ctx.Redirect("/files", fasthttp.StatusMovedPermanently)
}
