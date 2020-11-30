package lib

import "gorm.io/gorm"

var db *gorm.DB

// SetDatabase sets the global variable in this file to database instance created in the main function.
func SetDatabase(database *gorm.DB) {
	db = database
}

// GetDatabase returns a pointer, which points to the database instance.
func GetDatabase() *gorm.DB {
	return db
}
