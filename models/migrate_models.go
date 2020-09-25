package models

import "gorm.io/gorm"

func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&User{}, &File{})
}
