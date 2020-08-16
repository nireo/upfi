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
	Files []models.File
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	store := lib.GetStore()
	db := lib.GetDatabase()
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var user models.User
	if err := db.Where(&models.User{Username: session.Values["username"].(string)}).First(&user).Error; err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error: Cannot find form file in request")
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

	tempFile, err := ioutil.TempFile("./files", "file-*."+extension)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully Uploaded file")
}

func FilesController(w http.ResponseWriter, r *http.Request) {
	store := lib.GetStore()
	db := lib.GetDatabase()
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var user models.User
	if err := db.Where(&models.User{Username: session.Values["username"].(string)}).First(&user).Error; err != nil {
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
		Files: files,
	}

	tmpl.Execute(w, data)
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

	var user models.User
	if err := db.Where(&models.User{Username: session.Values["username"].(string)}).First(&user).Error; err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var file models.File
	if err := db.Where(&models.File{UUID: keys[0]}).First(&file).Error; err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	tmpl := template.Must(template.ParseFiles("./templates/single_file_template.html"))
	tmpl.Execute(w, file)
}