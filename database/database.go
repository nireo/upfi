package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/nireo/upfi/models"
)

func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.File{})
}

func InitializeDatabase() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "./database.db")
	if err != nil {
		panic("Cannot connect to database")
	}

	MigrateModels(db)
	return db, err
}
