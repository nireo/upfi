package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/middleware"

	"github.com/nireo/upfi/server"
	"net/http"
)

func main() {
	// Load database
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=upfi sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	lib.SetDatabase(db)

	// Setup HTTP Handler
	http.HandleFunc("/upload", server.UploadFile)
	http.HandleFunc("/register", middleware.Chain(server.AuthRegister, middleware.CheckMethod("POST"), middleware.LogRequest()))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8080", nil)
}
