package server

import (
	"fmt"
	"github.com/nireo/upfi/models"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/nireo/upfi/lib"
)

type FilePage struct {
	PageTitle string
	Files     []models.File
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	store := lib.GetStore()
	db := lib.GetDatabase()
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := models.FindOneFile(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// create file entry to the database
	newFileEntry := &models.File{
		Filename:    handler.Filename,
		UUID:        lib.GenerateUUID(),
		Description: r.FormValue("description"),
		Size:        handler.Size,
		UserID:      user.ID,
	}
	db.NewRecord(newFileEntry)
	db.Create(newFileEntry)

	defer file.Close()
	extension := lib.GetFileExtension(handler.Filename)

	userDirectory := fmt.Sprintf("./files/%s", user.UUID)
	tempFile, err := ioutil.TempFile(userDirectory, newFileEntry.UUID+extension)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "http://localhost:8080/files", http.StatusMovedPermanently)
}

func FilesController(w http.ResponseWriter, r *http.Request) {
	store := lib.GetStore()
	db := lib.GetDatabase()
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := models.FindOneUser(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var files []models.File
	if err := db.Model(&user).Related(&files).Error; err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("./templates/files_template.html"))
	data := FilePage{
		PageTitle: "Your files",
		Files:     files,
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func SingleFileController(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["file"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide file ID", http.StatusBadRequest)
		return
	}
	store := lib.GetStore()
	db := lib.GetDatabase()
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := models.FindOneUser(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	file, err := models.FindOneFile(&models.File{UUID: keys[0]}))
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("./templates/single_file_template.html"))
	err := tmpl.Execute(w, file)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["file"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide file ID", http.StatusBadRequest)
		return
	}
	store := lib.GetStore()
	db := lib.GetDatabase()
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := models.FindOneUser(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	file, err := models.FindOneFile(&models.File{UUID: keys[0]})
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	if file.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	db.Delete(&file)
	http.Redirect(w, r, "http://localhost:8080/files", http.StatusMovedPermanently)
}

func UpdateFile(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["file"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide file ID", http.StatusBadRequest)
		return
	}
	store := lib.GetStore()
	db := lib.GetDatabase()
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := models.FindOneUser(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	file, err := models.FindOneFile(&models.File{UUID: keys[0]})
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	title := r.FormValue("title")
	if title != "" {
		file.Filename = title
	}

	description := r.FormValue("description")
	if description != "" {
		file.Description = description
	}

	db.Save(&file)
	http.Redirect(w, r, "http://localhost:8080/files", http.StatusMovedPermanently)
}

func ServeUpdateForm(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["file"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide file ID", http.StatusBadRequest)
		return
	}
	store := lib.GetStore()
	db := lib.GetDatabase()
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := models.FindOneUser(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	file, err := models.FindOneFile(&models.File{UUID: keys[0]})
	if err != nil {
		http.Error(w, "Forbidden", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("./templates/update_file_info_template.html"))

	err := tmpl.Execute(w, file)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
