package models

import "github.com/jinzhu/gorm"

func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&User{}, &File{})
}
