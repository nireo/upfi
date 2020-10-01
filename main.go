package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/middleware"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/optimized_api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/nireo/upfi/server"

	"github.com/joho/godotenv"
)

func main() {
	// load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Could not load the env file")
	}

	user := os.Getenv("db_username")
	dbPort := os.Getenv("db_port")
	host := os.Getenv("db_host")
	dbName := os.Getenv("db_name")

	serverPort := os.Getenv("port")

	// Load database
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, dbPort, user, dbName),
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	models.MigrateModels(db)
	lib.SetDatabase(db)

	// options are 'optimized' and 'default'
	var apiVersion string
	flag.StringVar(&apiVersion, "api", "default", "Choose the api version")

	if apiVersion == "optimized" {
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
		if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil); err != nil {
			log.Fatal("Error in http.ListenAndServe")
		}
	} else {
		optimized_api.SetupOptimizedApi(serverPort)
	}
}
