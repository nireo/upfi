package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nireo/upfi/models"

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
		lib.HttpForbiddenHandler(w, r)
		return
	}

	user, err := models.FindOneUser(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		lib.HttpForbiddenHandler(w, r)
		return
	}

	switch r.Method {
	case http.MethodPost:
		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			lib.HttpInternalErrorHandler(w, r)
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
			Extension:   filepath.Ext(handler.Filename),
		}
		defer file.Close()

		fileDirectory := fmt.Sprintf("./files/%s/%s%s", user.UUID, newFileEntry.UUID, newFileEntry.Extension)
		dst, err := os.Create(fileDirectory)
		defer dst.Close()
		if err != nil {
			lib.HttpInternalErrorHandler(w, r)
			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			lib.HttpInternalErrorHandler(w, r)
			return
		}

		db.NewRecord(newFileEntry)
		db.Create(newFileEntry)
		http.Redirect(w, r, "http://localhost:8080/files", http.StatusMovedPermanently)
	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles("./static/upload.html"))
		err = tmpl.Execute(w, nil)
		if err != nil {
			lib.HttpInternalErrorHandler(w, r)
			return
		}
	}
}

func FilesController(w http.ResponseWriter, r *http.Request) {
	store := lib.GetStore()
	db := lib.GetDatabase()
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		lib.HttpForbiddenHandler(w, r)
		return
	}

	user, err := models.FindOneUser(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		lib.HttpForbiddenHandler(w, r)
		return
	}

	var files []models.File
	if err := db.Model(&user).Related(&files).Error; err != nil {
		lib.HttpInternalErrorHandler(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("./templates/files_template.html"))
	data := FilePage{
		PageTitle: "Your files",
		Files:     files,
	}

	if err = tmpl.Execute(w, data); err != nil {
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
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		lib.HttpForbiddenHandler(w, r)
		return
	}

	user, err := models.FindOneUser(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		lib.HttpForbiddenHandler(w, r)
		return
	}

	file, err := models.FindOneFile(&models.File{UUID: keys[0]})
	if err != nil {
		lib.HttpNotFoundHandler(w, r)
		return
	}

	if user.ID != file.UserID {
		lib.HttpForbiddenHandler(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("./templates/single_file_template.html"))
	if err = tmpl.Execute(w, file); err != nil {
		lib.HttpInternalErrorHandler(w, r)
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
		lib.HttpForbiddenHandler(w, r)
		return
	}

	user, err := models.FindOneUser(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		lib.HttpForbiddenHandler(w, r)
		return
	}

	file, err := models.FindOneFile(&models.File{UUID: keys[0]})
	if err != nil {
		lib.HttpNotFoundHandler(w, r)
		return
	}

	if file.UserID != user.ID {
		lib.HttpForbiddenHandler(w, r)
		return
	}

	os.Remove("./files/" + user.UUID + "/" + fmt.Sprintf("%s%s", file.UUID, file.Extension))
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
		lib.HttpForbiddenHandler(w, r)
		return
	}

	user, err := models.FindOneUser(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		lib.HttpForbiddenHandler(w, r)
		return
	}

	file, err := models.FindOneFile(&models.File{UUID: keys[0]})
	if err != nil {
		lib.HttpNotFoundHandler(w, r)
		return
	}

	if user.ID != file.UserID {
		// if the user doesn't own the file, he doesn't need to know of it's existence
		lib.HttpNotFoundHandler(w, r)
		return
	}

	switch r.Method {
	case http.MethodPost:
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
	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles("./templates/update_file_info_template.html"))
		if err = tmpl.Execute(w, file); err != nil {
			lib.HttpInternalErrorHandler(w, r)
			return
		}
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["file"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide file ID", http.StatusBadRequest)
		return
	}
	store := lib.GetStore()
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		lib.HttpForbiddenHandler(w, r)
		return
	}

	user, err := models.FindOneUser(&models.User{Username: session.Values["username"].(string)})
	if err != nil {
		lib.HttpForbiddenHandler(w, r)
		return
	}

	file, err := models.FindOneFile(&models.File{UUID: keys[0]})
	if err != nil {
		lib.HttpNotFoundHandler(w, r)
		return
	}

	if user.ID != file.UserID {
		lib.HttpNotFoundHandler(w, r)
		return
	}
	http.ServeFile(w, r, "./files/"+
		fmt.Sprintf("%s/%s%s", user.UUID, file.UUID, file.Extension))
}
