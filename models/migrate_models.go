package models

import (
	"log"

	"gorm.io/gorm"
)

// MigrateModels gets run in the main function and it migrates all of the database models
// to the database. This gets run everytime the service is restarted.
func MigrateModels(db *gorm.DB) {
	if err := db.AutoMigrate(&User{}, &File{}, &FileShare{}); err != nil {
		log.Fatal(err)
	}
}
