package jsonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nireo/upfi/crypt"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

// GetSingleFile returns the database entry, which contains data about a file. The user
// needs to provide a file id as a query. Also the files are kept private, so you need to own the file.
// Also the route is protected, so that the security token is checked before calling this handler.
func GetSingleFile(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	db := lib.GetDatabase()

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	fileID := ctx.UserValue("file").(string)
	file, err := models.FindFileAndCheckOwnership(user.ID, fileID)
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	lib.WriteResponseJSON(ctx, fasthttp.StatusOK, file)
}

// UploadFile handles the logic of uploading a file from the upload file form.
// Also the route is protected, so that the security token is checked before calling this handler.
func UploadFile(ctx *fasthttp.RequestCtx) {
	header, err := ctx.FormFile("file")
	if err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	db := lib.GetDatabase()
	user, err := models.FindOneUser(&models.User{Username: string(ctx.Request.Header.Peek("username"))})
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	if len(form.Value["master"]) == 0 {
		ServeErrorJSON(ctx, lib.BadRequestErrorPage)
		return
	}

	// Check that the user's master passwords is correct.
	if !lib.CheckPasswordHash(form.Value["master"][0], user.FileEncryptionMaster) {
		ServeErrorJSON(ctx, lib.ForbiddenErrorPage)
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

	// Define a path, where the file should be stored. Even though we encrypt the file, we
	// still want to keep the extension, since windows for example does not work without proper file
	// types.
	path := fmt.Sprintf("%s/%s/%s%s", lib.AddRootToPath("files"),
		user.UUID, newFileEntry.UUID, newFileEntry.Extension)

	// Read the file from the header. This is done because we need *multipart.File, which implements
	// io.Reader. This is needed to read the bytes in the file.
	multipartFile, err := header.Open()
	if err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	// Read the bytes of the file into a buffer.
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, multipartFile); err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	// Encrypt the data of the file using AESCipher and store it into the before defined path.
	if err := crypt.EncryptToDst(path, buf.Bytes(), form.Value["master"][0]); err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	// Read the mimetype so that we can set the content type properly
	// Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)
	// Copy the headers into the FileHeader buffer
	if _, err := multipartFile.Read(fileHeader); err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	newFileEntry.MIME = http.DetectContentType(fileHeader)

	db.Create(newFileEntry)
	lib.WriteResponseJSON(ctx, fasthttp.StatusOK, newFileEntry)
}

type updateFileBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateFile is http handler which takes a file id as a query parameter and checks for the file's ownership.
// This handler can be used to update file title and description.
// Also the route is protected, so that the security token is checked before calling this handler.
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

	file, err := models.FindFileAndCheckOwnership(user.ID, fileID)
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	file.Description = description
	file.Filename = title

	db.Save(&file)

	lib.WriteResponseJSON(ctx, fasthttp.StatusOK, file)
}

// GetUserFiles returns all the files that are related to the user who is requesting this
// handler. Then handler finds all the related files and constructs a simple json response.
// Also the route is protected, so that the security token is checked before calling this handler.
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

// DeleteFile is a handler that deletes a file owned by the user. The handler takes a file id as a query parameter
// and then does checking on the ownership of the file.
// Also the route is protected, so that the security token is checked before calling this handler.
func DeleteFile(ctx *fasthttp.RequestCtx) {
	username := string(ctx.Request.Header.Peek("username"))
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	fileID := ctx.UserValue("file").(string)
	file, err := models.FindFileAndCheckOwnership(user.ID, fileID)
	if err != nil {
		ServeErrorJSON(ctx, lib.NotFoundErrorPage)
		return
	}

	if err := os.Remove(lib.AddRootToPath("files/") + user.UUID + "/" + file.UUID + file.Extension); err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	db := lib.GetDatabase()
	db.Delete(&file)

	ctx.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}

type downloadFileBody struct {
	MasterPassword string `json:"master"`
}

func DownloadFile(ctx *fasthttp.RequestCtx) {
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

	var body downloadFileBody
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	path := fmt.Sprintf("%s/%s/%s%s", lib.AddRootToPath("files"),
		user.UUID, file.UUID, file.Extension)
	ctx.Response.Header.Set("Content-Type", file.MIME)

	tempUUID := lib.GenerateUUID()
	tempPath := fmt.Sprintf("%s/%s%s", lib.AddRootToPath("temp"), tempUUID, file.Extension)
	if err := crypt.DecryptToDst(tempPath, path, body.MasterPassword); err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}

	ctx.Response.SendFile(tempPath)

	if err := os.Remove(tempPath); err != nil {
		ServeErrorJSON(ctx, lib.InternalServerErrorPage)
		return
	}
}
