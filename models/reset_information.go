package models

import (
	"log"
	"os"
	"time"

	"github.com/nireo/upfi/lib"
)

// ResetInformation removes all contents of the files folder and all of the
// database entries.
func ResetInformation() {
	startingTime := time.Now()
	if err := os.RemoveAll("./files/"); err != nil {
		log.Fatal(err)
	}
	db := lib.GetDatabase()

	var users []User
	db.Find(&users)
	for _, user := range users {
		db.Delete(user)
	}

	var files []File
	db.Find(&files)
	for _, file := range files {
		db.Delete(file)
	}

	// Create a new emtpy files folder
	if err := os.Mkdir("./files", 0755); err != nil {
		log.Fatal("Could not setup ./files folder")
	}

	log.Printf("Successfully removed all data. took: %v", time.Since(startingTime))
}
