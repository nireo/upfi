package models

import "gorm.io/gorm"

func MigrateModels(db *gorm.DB) {
	if err := db.AutoMigrate(&User{}, &File{}); err != nil {
		panic(err)
	}
}
