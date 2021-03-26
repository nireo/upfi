package templateapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/nireo/upfi/crypt"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/templates"
	"github.com/valyala/fasthttp"
)

func formatFileSize(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// ServeUploadPage serves the requester a upload form, in which the user can upload files to their account.
// Also the route is protected, so that the security token is checked before calling this handler.
func ServeUploadPage(ctx *fasthttp.RequestCtx) {
	// Set the right Content-Type so that the html renders correctly.
	ctx.Response.Header.Set("Content-Type", "text/html")
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	templates.Upload(ctx, templates.UploadParams{
		Title: "upload",
		// we don't need to check if the user is authenticated, since a token is required to access route
		Authenticated: true,
	})
}

// UploadFile handles the logic of uploading a file from the upload file form.
// Also the route is protected, so that the security token is checked before calling this handler.
func UploadFile(ctx *fasthttp.RequestCtx) {
	// Get the file from the request form.
	header, err := ctx.FormFile("file")
	if err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}

	// Parse the multipart form so that we can check for other values, such as custom filenames or descriptions.
	form, err := ctx.MultipartForm()
	if err != nil {
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	// Find the user who requested this handler.
	db := lib.GetDatabase()
	username := string(ctx.Request.Header.Peek("username"))
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	if len(form.Value["master"]) == 0 {
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	// Check that the user's master passwords is correct.
	if !lib.CheckPasswordHash(form.Value["master"][0], user.FileEncryptionMaster) {
		ErrorPageHandler(ctx, lib.ForbiddenErrorPage)
		return
	}

	// Construct a database entry
	newFileEntry := &models.File{
		Filename:    header.Filename,
		UUID:        lib.GenerateUUID(),
		Description: form.Value["description"][0],
		Size:        header.Size,
		SizeHuman:   formatFileSize(header.Size),
		UserID:      user.ID,
		Extension:   filepath.Ext(header.Filename),
	}

	// Define a path, where the file should be stored. Even though we encrypt the file, we
	// still want to keep the extension, since windows for example does not work without proper file
	// types.
	path := fmt.Sprintf("%s/%s/%s%s", lib.AddRootToPath("files"),
		user.UUID, newFileEntry.UUID, newFileEntry.Extension)

	// Read the file from the header. This is done because we need *multipart.File, which implements
	// io.Reader. This is needed to read the bytes in the file.
	multipartFile, err := header.Open()
	if err != nil {
		fmt.Println("error 84")
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}

	// Read the bytes of the file into a buffer.
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, multipartFile); err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}

	// Encrypt the data of the file using AESCipher and store it into the before defined path.
	if err := crypt.EncryptToDst(path, buf.Bytes(), form.Value["master"][0]); err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}

	// Read the mimetype so that we can set the content type properly
	// Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)
	// Copy the headers into the FileHeader buffer
	if _, err := multipartFile.Read(fileHeader); err != nil && err != io.EOF {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}

	newFileEntry.MIME = http.DetectContentType(fileHeader)

	db.Create(newFileEntry)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Redirect("/files", fasthttp.StatusMovedPermanently)
}

// DownloadFile handler lets the user download a file. It also checks that the user owns the file he is trying download.
// Also the route is protected, so that the security token is checked before calling this handler.
func DownloadFile(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	fileID := ctx.UserValue("file").(string)
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Check that the user owns the file.
	if user.ID != file.UserID {
		ErrorPageHandler(ctx, lib.ForbiddenErrorPage)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		fmt.Println(err)
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	if len(form.Value["master"]) == 0 {
		fmt.Println(err)
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	path := fmt.Sprintf("%s/%s/%s%s", lib.AddRootToPath("files"),
		user.UUID, file.UUID, file.Extension)
	ctx.Response.Header.Set("Content-Type", file.MIME)

	tempUUID := lib.GenerateUUID()
	tempPath := fmt.Sprintf("%s/%s%s", lib.AddRootToPath("temp"),
		tempUUID, file.Extension)
	if err := crypt.DecryptToDst(tempPath, path, form.Value["master"][0]); err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}

	ctx.Response.SendFile(tempPath)

	if err := os.Remove(tempPath); err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
		return
	}
}

// GetSingleFile returns the database entry, which contains data about a file to the user. The user
// needs to provide a file id as a query. Also the files are kept private, so you need to own the file.
// Also the route is protected, so that the security token is checked before calling this handler.
func GetSingleFile(ctx *fasthttp.RequestCtx) {
	// Get the user's username which was appended to the request header
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	// Find the user's database entry who is requesting this handler.
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		fmt.Print("error in user")
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Find the file
	fileID := ctx.UserValue("file").(string)
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		fmt.Print("error in db")
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Check that the user owns the file.
	if user.ID != file.UserID {
		fmt.Println("error 196")
		// We return a not found error, since we don't want the unauthorized user to know about the file's existence.
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Display the user with the file's information, this template also includes the option to download a file.
	ctx.Response.Header.Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles(
		lib.AddRootToPath("templates/single_file_template.html")))
	if err := tmpl.Execute(ctx, file); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

type filePage struct {
	PageTitle string
	Files     []models.File
}

// GetUserFiles returns all the files that are related to the username which is requesting this
// handler. Then handler finds all the related files and constructs a template, which the user
// then can view as html content.
// Also the route is protected, so that the security token is checked before calling this handler.
func GetUserFiles(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	// Find the user's database entry who is requesting this handler.
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Find all file models which are related to the user in the database.
	var files []models.File
	db.Where(&models.File{UserID: user.ID}).Find(&files)

	for _, f := range files {
		f.SizeHuman = formatFileSize(f.Size)
	}

	pageParams := templates.FilesParams{
		Title: "your files",
		Files: files,
		// No need to check if the user is authenticated
		Authenticated: true,
	}

	if err := templates.Files(ctx, pageParams); err != nil {
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
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Find the file, if the file does not exist, return a not found error
	fileID := ctx.UserValue("file").(string)
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Check that the user owns the file.
	if user.ID != file.UserID {
		// Return a not found error, since we don't want the unauthorized user to know about the
		// file's existence
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Remove the file, if the file cannot be removed the return a internal server error to the user.
	if err := os.Remove(lib.AddRootToPath("files/") + user.UUID + "/" + fmt.Sprintf("%s%s", file.UUID, file.Extension)); err != nil {
		ErrorPageHandler(ctx, lib.InternalServerErrorPage)
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
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	// Check that the input fields are included, because without this check there will be a
	// index out of bounds error, if any of the fields are missing.
	if len(form.Value["title"]) == 0 || len(form.Value["description"]) == 0 {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
		return
	}

	title := form.Value["title"][0]
	description := form.Value["description"][0]

	// The description can be empty, but the title cannot
	if title == "" {
		ErrorPageHandler(ctx, lib.BadRequestErrorPage)
		return
	}

	// Find the user that is requesting this handler.
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Find the file
	fileID := ctx.UserValue("file").(string)
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

	// Check that the user owns the file
	if user.ID != file.UserID {
		ErrorPageHandler(ctx, lib.NotFoundErrorPage)
		return
	}

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
