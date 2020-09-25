package lib

import "gorm.io/gorm"

var db *gorm.DB

func SetDatabase(database *gorm.DB) {
	db = database
}

func GetDatabase() *gorm.DB {
	return db
}
