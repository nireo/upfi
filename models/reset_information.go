package models

import (
	"log"
	"os"
	"time"

	"github.com/nireo/upfi/lib"
	"golang.org/x/exp/errors/fmt"
)

// ResetInformation removes all contents of the files folder and all of the
// database entries.
func ResetInformation() {
	startingTime := time.Now()
	if err := os.RemoveAll(lib.AddRootToPath("files/")); err != nil {
		log.Fatal(err)
	}
	db := lib.GetDatabase()

	var users []User
	db.Find(&users)
	for _, user := range users {
		fmt.Printf("Deleted user %s\n", user.Username)
		db.Delete(User{Username: user.Username})
	}

	var files []File
	db.Find(&files)
	for _, file := range files {
		db.Delete(File{UUID: file.UUID})
	}

	// Create a new empty files folder
	if err := os.Mkdir(lib.AddRootToPath("files"), 0755); err != nil {
		log.Fatal("Could not setup files folder")
	}

	log.Printf("Successfully removed all data. took: %v", time.Since(startingTime))
}
