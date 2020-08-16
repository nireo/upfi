package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/middleware"
	"github.com/nireo/upfi/models"

	"github.com/nireo/upfi/server"
	"net/http"
)

func main() {
	// Load database
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=upfi sslmode=disable")
	if err != nil {
		panic(err)
	}
	models.MigrateModels(db)
	defer db.Close()
	lib.SetDatabase(db)

	// 			Setup HTTP Handler
	// Auth routes
	http.HandleFunc("/register", middleware.Chain(server.AuthRegister, middleware.LogRequest()))
	http.HandleFunc("/login", middleware.Chain(server.AuthLogin, middleware.LogRequest()))

	// File routes
	http.HandleFunc("/upload", server.UploadFile)
	http.HandleFunc("/file", server.SingleFileController)
	http.HandleFunc("/files", server.FilesController)
	http.HandleFunc("/delete", server.DeleteFile)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8080", nil)
}
