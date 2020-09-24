package main

import (
	"flag"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/middleware"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/optimized_api"

	"net/http"

	"github.com/nireo/upfi/server"
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

	// options are 'optimized' and 'default'
	var apiVersion string
	flag.StringVar(&apiVersion, "api", "default", "Choose the api version")

	if apiVersion == "default" {
		// Setup HTTP Handler
		// Auth routes
		http.HandleFunc("/register", middleware.Chain(server.AuthRegister, middleware.LogRequest()))
		http.HandleFunc("/login", middleware.Chain(server.AuthLogin, middleware.LogRequest()))

		// File routes
		http.HandleFunc("/download", server.DownloadFile)
		http.HandleFunc("/upload", server.UploadFile)
		http.HandleFunc("/file", server.SingleFileController)
		http.HandleFunc("/files", server.FilesController)
		http.HandleFunc("/delete", server.DeleteFile)
		http.HandleFunc("/update", server.UpdateFile)

		// User routes
		http.HandleFunc("/settings", server.SettingsPage)
		http.HandleFunc("/password", server.UpdatePassword)
		http.HandleFunc("/remove", server.DeleteUser)

		// Serve routes
		http.HandleFunc("/", server.ServeHomePage)

		// http.Handle("/", http.FileServer(http.Dir("./static")))
		_ = http.ListenAndServe(":8080", nil)
	} else if apiVersion == "optimized" {
		optimized_api.SetupOptimizedApi()
	}
}
