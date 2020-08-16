package server

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/nireo/upfi/lib"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error: Cannot find form file in request")
		return
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	defer file.Close()
	extension := lib.GetFileExtension(handler.Filename)

	tempFile, err := ioutil.TempFile("files", "file-*."+extension)
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

type File struct {
	Title string
}

type FilePage struct {
	PageTitle string
	Files []File
}

func FilesController(w http.ResponseWriter, r *http.Request) {
	store := lib.GetStore()
	session, _ := store.Get(r, "auth")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	tmpl := template.Must(template.ParseFiles("./templates/files_template.html"))
	data := FilePage{
		PageTitle: "Your files",
		Files: []File {
			{Title: "Videos"},
			{Title: "Messages"},
			{Title: "Documents"},
		},
	}

	tmpl.Execute(w, data)
}
