package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/nireo/upfi/crypt"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/templates"
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
func ServeUploadPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	templates.Upload(w, templates.UploadParams{
		Title:         "upload",
		Authenticated: true,
	})
}

func UploadFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseMultipartForm(50 << 20) // 50 mb
	file, header, err := r.FormFile("file")
	if err != nil {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	db := lib.GetDatabase()
	username := r.Header.Get("username")
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	if len(r.Form["master"]) == 0 {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	if !lib.CheckPasswordHash(r.Form["master"][0], user.FileEncryptionMaster) {
		ErrorPageHandler(w, r, lib.ForbiddenErrorPage)
		return
	}

	// Construct a database entry
	newFileEntry := &models.File{
		Filename:    header.Filename,
		UUID:        lib.GenerateUUID(),
		Description: r.Form["description"][0],
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

	// Read the bytes of the file into a buffer.
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}

	// Encrypt the data of the file using AESCipher and store it into the before defined path.
	if err := crypt.EncryptToDst(path, buf.Bytes(), r.Form["master"][0]); err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}

	// Read the mimetype so that we can set the content type properly
	// Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)
	// Copy the headers into the FileHeader buffer
	if _, err := file.Read(fileHeader); err != nil && err != io.EOF {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}

	newFileEntry.MIME = http.DetectContentType(fileHeader)

	db.Create(newFileEntry)
	http.Redirect(w, r, "/files", http.StatusMovedPermanently)
}

// GetSingleFile returns the database entry, which contains data about a file to the user. The user
// needs to provide a file id as a query. Also the files are kept private, so you need to own the file.
// Also the route is protected, so that the security token is checked before calling this handler.
func GetSingleFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the user's username which was appended to the request header
	username := r.Header.Get("username")
	db := lib.GetDatabase()

	// Find the user's database entry who is requesting this handler.
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	// Find the file
	fileID := ps.ByName("file")
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	// Check that the user owns the file.
	if user.ID != file.UserID {
		// We return a not found error, since we don't want the unauthorized user to know about the file's existence.
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	// Display the user with the file's information, this template also includes the option to download a file.
	w.Header().Set("Content-Type", "text/html")
	params := templates.SingleFileParams{
		Authenticated: true,
		Title:         file.Filename,
		File:          file,
	}

	templates.SingleFile(w, params)
}

// GetUserFiles returns all the files that are related to the username which is requesting this
// handler. Then handler finds all the related files and constructs a template, which the user
// then can view as html content.
// Also the route is protected, so that the security token is checked before calling this handler.
func GetUserFiles(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	username := r.Header.Get("username")
	db := lib.GetDatabase()

	// Find the user's database entry who is requesting this handler.
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
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

	if err := templates.Files(w, pageParams); err != nil {
		return
	}
}

// UpdateFile is http handler which takes a file id as a query parameter and checks for the file's ownership.
// This handler can be used to update file title and description.
// Also the route is protected, so that the security token is checked before calling this handler.
func UpdateFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("username")
	db := lib.GetDatabase()

	// Parse the multipart form so that we can take the 'title' and 'description' fields.
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	// Check that the input fields are included, because without this check there will be a
	// index out of bounds error, if any of the fields are missing.
	if len(r.Form["title"]) == 0 || len(r.Form["description"]) == 0 {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	title := r.Form["title"][0]
	description := r.Form["description"][0]

	// The description can be empty, but the title cannot
	if title == "" {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	// Find the user that is requesting this handler.
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	// Find the file
	fileID := ps.ByName("file")
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	// Check that the user owns the file
	if user.ID != file.UserID {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
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

	http.Redirect(w, r, "/files", http.StatusMovedPermanently)
}

// DeleteFile is a handler that deletes a file owned by the user. The handler takes a file id as a query parameter
// and then does checking on the ownership of the file.
// Also the route is protected, so that the security token is checked before calling this handler.
func DeleteFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("username")
	db := lib.GetDatabase()

	// Find the database entry of the user that requested this handler.
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	// Find the file, if the file does not exist, return a not found error
	fileID := ps.ByName("file")
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	// Check that the user owns the file.
	if user.ID != file.UserID {
		// Return a not found error, since we don't want the unauthorized user to know about the
		// file's existence
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	// Remove the file, if the file cannot be removed the return a internal server error to the user.
	if err := os.Remove(lib.AddRootToPath("files/") + user.UUID + "/" + fmt.Sprintf("%s%s", file.UUID, file.Extension)); err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}

	// Delete the database entry
	db.Delete(&file)
	http.Redirect(w, r, "/files", http.StatusMovedPermanently)
}

// DownloadFile handler lets the user download a file. It also checks that the user owns the file he is trying download.
// Also the route is protected, so that the security token is checked before calling this handler.
func DownloadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("username")
	db := lib.GetDatabase()

	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	fileID := ps.ByName("file")
	var file models.File
	if err := db.Where(&models.File{UUID: fileID}).First(&file).Error; err != nil {
		ErrorPageHandler(w, r, lib.NotFoundErrorPage)
		return
	}

	// Check that the user owns the file.
	if user.ID != file.UserID {
		ErrorPageHandler(w, r, lib.ForbiddenErrorPage)
		return
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	if len(r.Form["master"]) == 0 {
		ErrorPageHandler(w, r, lib.BadRequestErrorPage)
		return
	}

	path := fmt.Sprintf("%s/%s/%s%s", lib.AddRootToPath("files"),
		user.UUID, file.UUID, file.Extension)
	w.Header().Set("Content-Type", file.MIME)

	tempUUID := lib.GenerateUUID()
	tempPath := fmt.Sprintf("%s/%s%s", lib.AddRootToPath("temp"),
		tempUUID, file.Extension)
	if err := crypt.DecryptToDst(tempPath, path, r.Form["master"][0]); err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}

	http.ServeFile(w, r, tempPath)
	if err := os.Remove(tempPath); err != nil {
		ErrorPageHandler(w, r, lib.InternalServerErrorPage)
		return
	}
}
