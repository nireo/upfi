package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func getFileExtension(fileName string) string {
	splitted := strings.Split(fileName, ".")
	return splitted[len(splitted)-1]
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
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
	extension := getFileExtension(handler.Filename)

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

func main() {
	http.HandleFunc("/upload", uploadFile)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8080", nil)
}
